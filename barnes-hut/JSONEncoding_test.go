package barnes_hut

import (
	"os"
	"github.com/thomaspeugeot/tkv/quadtree"
	"testing"
	"encoding/json"
)

// test 
func TestEncodingSliceOfSlice(t *testing.T) {

	// var sliceOfSlice 
	// sliceOfSlice := make( [][]float64, 2)
	// sliceOfSlice[0] = make( []float64, 2)
	// sliceOfSlice[1] = make( []float64, 2)

	sliceOfSlice := [][]float64 {
		{ 0.0, 1.0},
		{ 2.0, 3.0},
	}
	_, err := json.MarshalIndent( sliceOfSlice, "","\t")
	if( err != nil) { t.Errorf("error") }
	// fmt.Println( string( stats))
}

// test gif output
func TestEncodingGini(t *testing.T) {

	bodies := make([]quadtree.Body, 20000)
	SpreadOnCircle( & bodies)
	
	var r Run
	r.Init( & bodies)
	
	var output *os.File
	output, _ = os.Create("essai.gif")
	
	r.SetState( RUNNING)
	r.OutputGif( output, 3)

	_, err := json.MarshalIndent( r.giniOverTime, "","\t")
	if( err != nil) { t.Errorf("error") }

	// visual verification
}

