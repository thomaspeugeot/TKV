package barneshut

import (
	"math"
	"testing"

	"github.com/thomaspeugeot/tkv/quadtree"
)

// a RepulsionField stores the computation
// of a scalar field with the values of the repulsion field (1/r)
// on a fixed area, at interpolation points ( GridFieldTicks interpolation points per dimension )
// this structure is transcient
type RepulsionField struct {
	XMin, YMin, XMax, YMax float64 // coordinate of the rendering area
	GridFieldTicks         int     // nb of intervals where the field is computed
	// q * quadtree.Quadtree // quadtree used to compute the field
	values   [][]float64 // values of the field
	maxValue float64
	q        *quadtree.Quadtree // the quadtree against which the field is computed
	cutoff   float64            // the distance to the nearest body with void the repulsion field
}

func NewRepulsionField(XMin, YMin, XMax, YMax float64, GridFieldTicks int, q *quadtree.Quadtree, cutoff float64) *RepulsionField {
	Trace.Println("NewRepulsionField")

	var f RepulsionField
	// to be replaced with a proper init of struct to
	f.XMin = XMin
	f.YMin = YMin
	f.XMax = XMax
	f.YMax = YMax

	f.GridFieldTicks = GridFieldTicks
	f.q = q

	f.cutoff = cutoff

	f.values = make([][]float64, GridFieldTicks)
	for i := range f.values {
		f.values[i] = make([]float64, GridFieldTicks)
	}
	return &f
}

// compute x,y coordinates in the square from the i,j coordinate in the field
func (f *RepulsionField) XY(i, j int) (x, y float64) {

	x = f.XMin + ((float64(i)+0.5)/float64(f.GridFieldTicks))*(f.XMax-f.XMin)
	y = f.YMin + ((float64(j)+0.5)/float64(f.GridFieldTicks))*(f.YMax-f.YMin)

	return x, y
}

// compute repulsion at body A coordinates from body B
// repulsion field is in 1/r * M
func getRepulsionField(A, B *quadtree.Body) (v float64) {

	x := getModuloDistance(B.X, A.X)
	y := getModuloDistance(B.Y, A.Y)

	distQuared := (x*x + y*y)
	absDistance := math.Sqrt(distQuared + ETA)
	v = B.M / absDistance

	return v
}

// compute field for all interpolation points
// concurrently
func (f *RepulsionField) ComputeField() {
	Trace.Println("ComputeField nbTicks ", f.GridFieldTicks, len(f.values))
	f.maxValue = 0.0

	// done := make( chan float64)
	for i, vs := range f.values {
		for j := range vs {

			x, y := f.XY(i, j)
			var rootCoord quadtree.Coord
			// go func() {
			var fv float64 // field value
			// 	// am i sure that have not been changed by the next call to func ?
			f.ComputeFieldRecursive(x, y, f.q, rootCoord, &fv)
			// 	done <- fv
			// }()
			if fv > f.maxValue {
				f.maxValue = fv
			}
			vs[j] = fv

		}
	}
	// for i,vs := range f.values {
	// 	for j,_ := range vs {
	// 		var fv float64  // field value
	// 		fv = <- done
	// 		vs[j] = fv
	// 		if fv > f.maxValue { f.maxValue = fv}
	// 		x, y := f.XY( i, j)
	// 		Trace.Printf("computeField at %d %d %e %e, v = %e\n", i, j, x, y, vs[j])
	// 	}
	// }
	Trace.Printf("computeField maxValue %e\n", f.maxValue)
}

// compute repulsion field at interpolation point x, y and update v
func (f *RepulsionField) ComputeFieldRecursive(x, y float64, q *quadtree.Quadtree, coord quadtree.Coord, v *float64) {

	Trace.Printf("ComputeFieldRecursive at %e %e, quadtree %p, coord %s, input v = %e\n", x, y, q, coord.String(), *v)

	var body quadtree.Body
	body.X = x
	body.Y = y

	// compute the node box size
	level := coord.Level()
	boxSize := 1.0 / math.Pow(2.0, float64(level)) // if level = 0, this is 1.0

	node := &(q.Nodes[coord])
	distToNode := getModuloDistanceBetweenBodies(&body, &(node.Body))

	// avoid node with zero mass
	if node.M == 0 {
		return
	}

	// check if the COM of the node can be used
	if (boxSize / distToNode) < BN_THETA {

		*v += getRepulsionField(&body, &(node.Body))

	} else {
		if level < 8 {
			// parse sub nodes
			Trace.Printf("ComputeFieldRecursive go down at node %#v\n", node.Coord())

			coordNW, coordNE, coordSW, coordSE := quadtree.NodesBelow(coord)
			f.ComputeFieldRecursive(x, y, q, coordNW, v)
			f.ComputeFieldRecursive(x, y, q, coordNE, v)
			f.ComputeFieldRecursive(x, y, q, coordSW, v)
			f.ComputeFieldRecursive(x, y, q, coordSE, v)
		} else {

			// parse bodies of the node
			rank := 0
			rankOfBody := -1
			for b := node.First(); b != nil; b = b.Next() {
				if *b != body {

					dist := getModuloDistanceBetweenBodies(&body, b)

					if dist == 0.0 {
						var t testing.T
						q.CheckIntegrity(&t)

						// c1 := body.Coord()
						// c2 := b.Coord()
						Error.Printf("Problem at rank %d for body of rank %d on node %#v ",
							rank, rankOfBody, *node)
					}

					// if distance is inferior to cutoff, return
					if dist < f.cutoff {
						*v = 0.0
						return

					} else {
						*v += getRepulsionField(&body, b)

					}

					rank++
					Trace.Printf("ComputeFieldRecursive at leaf %#v rank %d x %9.3f y %9.3f\n", b.Coord(), rank, x, y)
				} else {
					rankOfBody = rank
				}
			}
		}
	}
	return
}
