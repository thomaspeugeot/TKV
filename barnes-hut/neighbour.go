package barnes_hut

import (
	"github.com/thomaspeugeot/tkv/quadtree"
)

// relative to a body of interest, the storage for a neighbour with its distance 
// nota : this is used to measure the stiring of the bodies along the simulation
type Neighbour struct {
	n * quadtree.Body // rank in the []quadtree.Body
	Distance float64
}
// the measure of stiring is computed with a finite number of neighbours
// no stiring is that the neighbours at the end of the run are the same neighbours
// that at the begining
var NbOfNeighboursPerBody int = 10

type NeighbourDico * [][]Neighbour

func (r * Run) InitNeighbourDico( bodies * ([]quadtree.Body)) {
	neighbours := make([][]Neighbour, len(*bodies))
	r.bodiesNeighbours = &neighbours
	for idx,_  := range *r.bodiesNeighbours {
		(*r.bodiesNeighbours)[idx] = make( []Neighbour, NbOfNeighboursPerBody)
	}
	neighboursOrig := make([][]Neighbour, len(*bodies))
	r.bodiesNeighboursOrig = &neighboursOrig
	for idx,_  := range *r.bodiesNeighboursOrig {
		(*r.bodiesNeighboursOrig)[idx] = make( []Neighbour, NbOfNeighboursPerBody)
	}
} 