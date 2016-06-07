package barnes_hut

import (
	// "github.com/thomaspeugeot/tkv/quadtree"
)

// a RepulsionField stores the computation
// of a scalar field with the values of the repulsion field (1/r)
// on a fixed area
// this structure is transcient 
type RepulsionField struct {
	XMin, YMin, XMax, YMax float64 // coordinate of the rendering area
	GridFieldTicks int // nb of intervals where the field is computed
	// q * quadtree.Quadtree // quadtree used to compute the field
	values [][]float64 // values of the field
}

func NewRepulsionField(GridFieldTicks int) * RepulsionField {
	Info.Println("NewRepulsionField")

	var f RepulsionField
	f.GridFieldTicks = GridFieldTicks
	f.values = make ( [][]float64, GridFieldTicks)
	for i,_ := range f.values {
		f.values[i] = make ( []float64, GridFieldTicks)
	}
	return &f
}

func (f * RepulsionField) ComputeField() {
	Info.Println("ComputeField nbTicks ", f.GridFieldTicks, len( f.values))

	for i,vs := range f.values {
		for j,v := range vs {
			Trace.Printf("computeField at %d %d, v = %e\n", i, j, v)
		}
	} 
}


