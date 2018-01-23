package translation


import (
	"os"
	"log"
	"fmt"
	"math"
	"encoding/json"
	
	"github.com/thomaspeugeot/tkv/barnes-hut"
	"github.com/thomaspeugeot/tkv/quadtree"
	"github.com/thomaspeugeot/tkv/grump"

)

type Country struct {
	grump.Country

	NbBodies int // nb of bodies according to the filename

	bodiesOrig * []quadtree.Body // original bodies position in the quatree
	bodiesSpread * []quadtree.Body // bodies position in the quatree after the spread simulation
	VilCoordinates [][]int
	Step int // step when the simulation stopped
	
	villages [][]Village
}

type Point struct {
	X, Y float64
}

func MakePoint(x float64, y float64) Point {
	return Point{X: x, Y: y}
}

type PointList []Point

type BodySetChoice string
const (
	ORIGINAL_CONFIGURATION = "ORIGINAL_CONFIGURATION"
	SPREAD_CONFIGURATION = "SPREAD_CONFIGURATION"
)

// number of village per X or Y axis. For 10 000 villages, this number is 100
// this value can be set interactively during the run
var nbVillagePerAxe int = 100 
var numberOfVillagePerAxe float64 = 100.0

// init variables
func (country * Country) Init() {

	// unserialize from conf-<country trigram>.coord
	// store step because the unseralize set it to a wrong value
	step := country.Step
	country.Unserialize()
	country.Step = step

	Info.Printf("Init after Unserialize name %s", country.Name)
	Info.Printf("Init after Unserialize step %d", country.Step)

	country.LoadConfig( true ) // load config at the end  of the simulation
	country.LoadConfig( false ) // load config at the start of the simulation

	// init village array
	country.villages = make( [][]Village, nbVillagePerAxe )
	for x,_  := range country.villages {
		country.villages[x] = make([]Village, nbVillagePerAxe)
	}

	country.VilCoordinates = make( [][]int, country.NbBodies)
	for idx, _ := range country.VilCoordinates {
		country.VilCoordinates[idx] = make( []int, 2)
	}

	country.ComputeBaryCenters()
	
}

// load configuration from filename into country 
// check that it matches the 
func (country * Country) LoadConfig( isOriginal bool) bool {

	Info.Printf( "Load Config begin : Country is %s, step %d isOriginal %t", country.Name, country.Step, isOriginal)

	// computing the file name from the step
	step := 0
	
	// if isOrignal load the file with the step number 0, else use spread
	if ! isOriginal { step = country.Step }

	filename := fmt.Sprintf( barnes_hut.CountryBodiesNamePattern, country.Name, country.NbBodies, step)
	Info.Printf( "LoadConfig original %t file %s for country %s at step %d", isOriginal, filename, country.Name, step)

	file, err := os.Open(filename)
	if( err != nil) {
		log.Fatal(err)
		return false
	}

	jsonParser := json.NewDecoder(file)

	bodies := (make([]quadtree.Body, 0))
	if isOriginal {
		country.bodiesOrig = & bodies
		if err = jsonParser.Decode( country.bodiesOrig); err != nil {
			log.Fatal( fmt.Sprintf( "parsing config file %s", err.Error()))
		}
		Info.Printf( "nb item parsed in file for orig %d\n", len( *country.bodiesOrig))
	} else {
		country.bodiesSpread = & bodies
		if err = jsonParser.Decode( country.bodiesSpread); err != nil {
			log.Fatal( fmt.Sprintf( "parsing config file %s", err.Error()))
		}
		Info.Printf( "nb item parsed in file for spread %d\n", len( *country.bodiesSpread))
	}

	file.Close()

	Info.Printf( "Load Config end : Country is %s, step %d", country.Name, country.Step)
	
	return true
}

// compute villages barycenters
func (country * Country) ComputeBaryCenters() {
	Info.Printf("ComputeBaryCenters begins for country %s", country.Name)

	// parse bodiesSpread to compute bary centers 
	// use bodiesOrig to compute bary centers
	for index,b := range *country.bodiesSpread {

		// compute village coordinate (from 0 to nbVillagePerAxe-1)
		villageX := int( math.Floor(float64( nbVillagePerAxe) * b.X))
		villageY := int( math.Floor(float64( nbVillagePerAxe) * b.Y))

		Trace.Printf("Adding body index %d to village %d %d", index, villageX, villageY)

		// add body (original) to the barycenter of the village
		bOrig := (*country.bodiesOrig)[index]
		country.villages[villageX][villageY].addBody( bOrig)

		country.VilCoordinates[index][0] = villageX
		country.VilCoordinates[index][1] = villageY
	}
}

// given lat, lng
// return xClosest, yClosest, distance, latClosest, lngClosest of the closest body within the country (original position)
func (country * Country) ClosestBodyInOriginalPosition( lat, lng float64) (xRelClosest, yRelClosest, distance, latClosest, lngClosest, xSpread, ySpread float64, closestIndex int) {

	// compute relative coordinates within the square
	xRel, yRel := country.LatLng2XY( lat, lng)

	// parse all bodies and get closest body
	closestIndex = -1
	minDistance := 1000000000.0 // we start from away
	for index,b := range *country.bodiesOrig {
		distanceX := b.X - xRel
		distanceY := b.Y - yRel
		distance := math.Sqrt( (distanceX*distanceX) + (distanceY*distanceY))

		if( distance < minDistance ) { 
			closestIndex = index 
			minDistance = distance
		}
	}	

	xRelClosest = (*country.bodiesOrig)[closestIndex].X
	yRelClosest = (*country.bodiesOrig)[closestIndex].Y

	latOptimClosest, lngOptimClosest := country.XY2LatLng( xRelClosest, yRelClosest)
	

	Info.Printf( "ClosestBodyInOriginalPosition %f %f relative to country %f %f", lat, lng, xRel, yRel)
	Info.Printf( "ClosestBodyInOriginalPosition rel closest %f %f lat lng closest %f %f", xRelClosest, yRelClosest, latOptimClosest, lngOptimClosest)

	// compute x, y in spread bodies
	xSpread = (*country.bodiesSpread)[closestIndex].X
	ySpread = (*country.bodiesSpread)[closestIndex].Y

	Info.Printf( "ClosestBodyInOriginalPosition village %f %f index %d", xSpread, ySpread, closestIndex)

	return xRelClosest, yRelClosest, minDistance, latOptimClosest, lngOptimClosest, xSpread, ySpread, closestIndex
}

func (country * Country) XYSpreadToLatLngOrig( x, y float64) (lat, lng float64) {

	Info.Printf( "XYSpreadToLatLngOrig input x %f y %f", x, y)

	// parse all bodies and get closest body
	closestIndex := -1
	minDistance := 1000000000.0 // we start from away
	for index,b := range *country.bodiesSpread {
		distanceX := b.X - x
		distanceY := b.Y - y
		distance := math.Sqrt( (distanceX*distanceX) + (distanceY*distanceY))

		if( distance < minDistance ) { 
			closestIndex = index 
			minDistance = distance
		}
	}	

	xRelClosest := (*country.bodiesOrig)[closestIndex].X
	yRelClosest := (*country.bodiesOrig)[closestIndex].Y
	latOptimClosest, lngOptimClosest := country.XY2LatLng( xRelClosest, yRelClosest)
	Info.Printf("XYSpreadToLatLngOrig target x %f y %f index %d distance %f", xRelClosest, yRelClosest, closestIndex, minDistance)

	Info.Printf("XYSpreadToLatLngOrig target lat %f lng %f", latOptimClosest, lngOptimClosest)

	return latOptimClosest, lngOptimClosest
}

// get the bodies of a village from x, y spread coordinates
func (country * Country) XYSpreadToLatLngOrigVillage( x, y float64) PointList {

	Info.Printf( "XYSpreadToLatLngOrigVillage Country %s x %f y %f", country.Name, x, y)
	
	points := make(PointList, 0)

	// compute village min & max coord
	xMinVillage := float64( int( x*numberOfVillagePerAxe))/numberOfVillagePerAxe
	xMaxVillage := float64( int( x*numberOfVillagePerAxe + 1.0))/numberOfVillagePerAxe
	yMinVillage := float64( int( y*numberOfVillagePerAxe))/numberOfVillagePerAxe
	yMaxVillage := float64( int( y*numberOfVillagePerAxe + 1.0))/numberOfVillagePerAxe
	
	Info.Printf( "XYSpreadToLatLngOrig input village Min x %f Max x %f", xMinVillage, xMaxVillage)
	
	// parse all bodies and get closest body
	for index,b := range *country.bodiesSpread {
		if (xMinVillage <= b.X) && (b.X < xMaxVillage) && (yMinVillage <= b.Y) && (b.Y < yMaxVillage) {

			xRelClosest := (*country.bodiesOrig)[index].X
			yRelClosest := (*country.bodiesOrig)[index].Y
			latOptimClosest, lngOptimClosest := country.XY2LatLng( xRelClosest, yRelClosest)
			
			points = append(points, MakePoint(latOptimClosest, lngOptimClosest))
		}
	}	

	Info.Printf( "XYSpreadToLatLngOrigVillage Country %s nb bodies in village %d", country.Name, len(points))
	return points
}

// given x, y of a point, return the border in the country
func (country * Country) VillageBorder( lat, lng float64) PointList {

	Info.Printf( "")
	Info.Printf( "VillageBorder country %s input lat %f lng %f", country.Name, lat, lng)
	
	// from input lat, lng, get the xSpread, ySpread
	_, _, _, _, _, xSpread, ySpread, _ := country.ClosestBodyInOriginalPosition(lat, lng)
	Info.Printf( "VillageBorder country %s input xSpread %f ySpread %f", country.Name, xSpread, ySpread)

	// compute village min & max coord
	xMinVillage := float64( int( xSpread*numberOfVillagePerAxe))/numberOfVillagePerAxe
	xMaxVillage := float64( int( xSpread*numberOfVillagePerAxe + 1.0))/numberOfVillagePerAxe
	yMinVillage := float64( int( ySpread*numberOfVillagePerAxe))/numberOfVillagePerAxe
	yMaxVillage := float64( int( ySpread*numberOfVillagePerAxe + 1.0))/numberOfVillagePerAxe
	
	Info.Printf( "VillageBorder input village Min x %f Max x %f", xMinVillage, xMaxVillage)
	
	// parse all bodies and if bodies has x & y spead close to input spread, include them in point list
	points := make(PointList, 0)
	for index,b := range *country.bodiesSpread {
		if (xMinVillage <= b.X) && (b.X < xMaxVillage) && (yMinVillage <= b.Y) && (b.Y < yMaxVillage) {

			xRelClosest := (*country.bodiesOrig)[index].X
			yRelClosest := (*country.bodiesOrig)[index].Y
			latOptimClosest, lngOptimClosest := country.XY2LatLng( xRelClosest, yRelClosest)
			
			points = append(points, MakePoint(latOptimClosest, lngOptimClosest))
		}
	}	

	Info.Printf( "VillageBorder nb of border points %d", len(points))
	Info.Printf( "")

	return points
}
