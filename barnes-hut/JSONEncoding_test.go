package barneshut

import (
	"encoding/json"
	"testing"
)

// test
func TestEncodingSliceOfSlice(t *testing.T) {

	// var sliceOfSlice
	// sliceOfSlice := make( [][]float64, 2)
	// sliceOfSlice[0] = make( []float64, 2)
	// sliceOfSlice[1] = make( []float64, 2)

	sliceOfSlice := [][]float64{
		{0.0, 1.0},
		{2.0, 3.0},
	}
	_, err := json.MarshalIndent(sliceOfSlice, "", "\t")
	if err != nil {
		t.Errorf("error")
	}
	// fmt.Println( string( stats))
}
