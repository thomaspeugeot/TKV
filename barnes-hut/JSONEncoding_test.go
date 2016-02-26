package barnes_hut

import (
	"os"
	"fmt"
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
	stats, _ := json.MarshalIndent( sliceOfSlice, "","\t")
	fmt.Println( string( stats))
}

// test gif output
func TestEncodingGini(t *testing.T) {

	bodies := make([]quadtree.Body, 200000)
	SpreadOnCircle( & bodies)
	
	var r Run
	r.Init( & bodies)
	
	var output *os.File
	output, _ = os.Create("essai.gif")
	
	r.SetState( RUNNING)
	r.OutputGif( output, 3)

	// for step :=0; step < len(r.giniOverTime); step++ {

	// 	for i:= 0; i<=9; i++ {
	// 		s := fmt.Sprintf( "%f ", r.giniOverTime[ step][ i])
	// 		fmt.Printf( s)
	// 	}
	// 	fmt.Println()
	// }

	stats, _ := json.MarshalIndent( r.giniOverTime, "","\t")
	fmt.Println( string( stats))

	// visual verification
}

