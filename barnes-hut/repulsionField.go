package barnes_hut

import (
	"github.com/thomaspeugeot/tkv/quadtree"
)

// a RepulsionField stores the computation
// of a scalar field with the values of the repulsion field (1/r)
// on a fixed area, at interpolation points ( GridFieldTicks interpolation points per dimension )
// this structure is transcient 
type RepulsionField struct {
	XMin, YMin, XMax, YMax float64 // coordinate of the rendering area
	GridFieldTicks int // nb of intervals where the field is computed
	// q * quadtree.Quadtree // quadtree used to compute the field
	values [][]float64 // values of the field
	q * quadtree.Quadtree // the quadtree against which the field is computed
}

func NewRepulsionField( XMin, YMin, XMax, YMax float64, GridFieldTicks int, q * quadtree.Quadtree) * RepulsionField {
	Info.Println("NewRepulsionField")

	var f RepulsionField
	// to be replaced with a proper init of struct to 
	f.XMin = XMin
	f.YMin = YMin
	f.XMax = XMax
	f.YMax = YMax

	f.GridFieldTicks = GridFieldTicks
	f.values = make ( [][]float64, GridFieldTicks)
	for i,_ := range f.values {
		f.values[i] = make ( []float64, GridFieldTicks)
	}
	return &f
}

// compute x,y coordinates in the square from the i,j coordinate in the field
func (f * RepulsionField) XY(i, j int) (x,y float64) {

	x = f.XMin + ((float64(i)+0.5) / float64( f.GridFieldTicks)) * (f.XMax - f.XMin)
	y = f.YMin + ((float64(j)+0.5) / float64( f.GridFieldTicks)) * (f.YMax - f.YMin)

	return x, y
}

// compute field for all interpolation points
func (f * RepulsionField) ComputeField() {
	Info.Println("ComputeField nbTicks ", f.GridFieldTicks, len( f.values))

	for i,vs := range f.values {
		for j,v := range vs {
			
			x, y := f.XY( i, j)
			

			var rootCoord quadtree.Coord
			f.ComputeFieldAtInterpolationPointRecursive( x, y, f.q, rootCoord)
			Info.Printf("computeField at %d %d %e %e, v = %e\n", i, j, x, y, v)
		}
	} 
}

func (f * RepulsionField) ComputeFieldAtInterpolationPointRecursive( x, y float64, q * quadtree.Quadtree, coord quadtree.Coord) float64 {

	return 0.0
}




