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

	bodiesOrig * []quadtree.Body // original bodies position in the quatree
	bodies * []quadtree.Body // bodies position in the quatree
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

	// get country coordinates
	country.Unserialize()
	country.LoadConfig( true )
	country.LoadConfig( false )
	country.ComputeBaryCenters()
	
}

// load configuration from filename into counry 
// check that it matches the 
func (country * Country) LoadConfig( isOriginal bool) bool {

	// computing the file name from the step
	var step int
	
	if isOriginal { step = country.Step }

	filename := fmt.Sprintf( barnes_hut.CountryBodiesNamePattern, country.Name, step)
	Info.Printf( "LoadConfig file %s for country %s at step %d", filename, country.Name, step)

	file, err := os.Open(filename)
	if( err != nil) {
		log.Fatal(err)
		return false
	}

	// get the number of steps in the file name
	// var countryName string
	for index, runeValue := range filename {
    	Trace.Printf("%#U starts at byte position %d\n", runeValue, index)
	}
	ctry := filename[5:8]
	stepString := filename[9:14]
	
	nbItems, errScan := fmt.Sscanf(stepString, "%05d", & country.Step)
	if( errScan != nil) {
		log.Fatal(errScan)
		return false			
	}
	Trace.Printf( "nb item parsed in file name %d (should be one)\n", nbItems)
	
	jsonParser := json.NewDecoder(file)

	bodies := (make([]quadtree.Body, 0))
	if isOriginal {
		country.bodiesOrig = & bodies
		if err = jsonParser.Decode( country.bodiesOrig); err != nil {
			log.Fatal( fmt.Sprintf( "parsing config file %s", err.Error()))
		}
		Info.Printf( "nb item parsed in file %d\n", len( *country.bodiesOrig))
	} else {
		country.bodies = & bodies
		if err = jsonParser.Decode( country.bodies); err != nil {
			log.Fatal( fmt.Sprintf( "parsing config file %s", err.Error()))
		}
		Info.Printf( "nb item parsed in file %d\n", len( *country.bodies))
	}
	Info.Printf( "Country is %s, step is %d", ctry, country.Step)

	file.Close()
	
	return true
}

// compute villages barycenters
func (country * Country) ComputeBaryCenters() {
	
	country.villages = make( [][]Village, nbVillagePerAxe )
	
	for x,_  := range country.villages {
		country.villages[x] = make([]Village, nbVillagePerAxe)
	}

	// parse bodies
	for _,b := range *country.bodies {
		// compute village coordinate (from 0 to nbVillagePerAxe-1)
		x := int( math.Floor(float64( nbVillagePerAxe) * b.X))
		y := int( math.Floor(float64( nbVillagePerAxe) * b.Y))

		// add body to the barycenter of the village
		country.villages[x][y].addBody( b)
	}


}
