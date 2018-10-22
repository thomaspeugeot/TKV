package translation

import (
	"github.com/thomaspeugeot/tkv/quadtree"
)

type Village struct {
	quadtree.Body     // barycenter of the villages
	NbBodies      int // nb of bodies in the village (maybe not useful)
}

func (v *Village) reset() {
	v.M = 0
	v.NbBodies = 0
}

// add a body to the village, update barycenter
func (v *Village) addBody(b quadtree.Body) {

	v.X = ((v.X * v.M) + (b.X * b.M)) / (v.M + b.M)
	v.M += b.M
	v.NbBodies++
}
