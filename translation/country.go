package translation


import (
	"os"
	"log"
	"fmt"
	"math"
	// "bufio"
	"github.com/thomaspeugeot/tkv/barnes-hut"
	"github.com/thomaspeugeot/tkv/quadtree"
	"github.com/thomaspeugeot/tkv/grump"
	// "path/filepath"
	"encoding/json"
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

type BodySetChoice string
const (
	ORIGINAL_CONFIGURATION = "ORIGINAL_CONFIGURATION"
	SPREAD_CONFIGURATION = "SPREAD_CONFIGURATION"
)

// number of village per X or Y axis. For 10 000 villages, this number is 100
// this value can be set interactively during the run
var nbVillagePerAxe int = 100 


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

// load configuration from filename into counry 
// check that it matches the 
func (country * Country) LoadConfig( isOriginal bool) bool {

	Info.Printf( "Load Config begin : Country is %s, step %d", country.Name, country.Step)

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

func (country * Country) VillageCoordinates( lat, lng float64) (x, y int, distance, latClosest, lngClosest float64) {

	// compute relative coordinates within the square
	xRel, yRel := country.LatLng2XY( lat, lng)
	Info.Printf("VillageCoordinates lat %f,  lng %f", lat, lng)
	Info.Printf("Rel x %f, Rel y %f", xRel, yRel)

	// parse all bodies and get closest body
	closestIndex := -1
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
	Info.Printf("VillageCoordinates closestIndex %d, minDistance %f", closestIndex, minDistance)

	villageX := country.VilCoordinates[closestIndex][0]
	villageY := country.VilCoordinates[closestIndex][1]
	xRelClosest := (*country.bodiesOrig)[closestIndex].X
	yRelClosest := (*country.bodiesOrig)[closestIndex].Y

	latOptimClosest, lngOptimClosest := country.XY2LatLng( xRelClosest, yRelClosest)
	

	Info.Printf( "VillageCoordinates %f %f relative to country %f %f", lat, lng, xRel, yRel)
	Info.Printf( "VillageCoordinates rel closest %f %f lat lng Closet %f %f", xRelClosest, yRelClosest, latOptimClosest, lngOptimClosest)
	Info.Printf( "VillageCoordinates village %d %d", villageX, villageY)
	return villageX, villageY, minDistance, latOptimClosest, lngOptimClosest
}
