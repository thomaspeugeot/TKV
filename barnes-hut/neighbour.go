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

type NeighbourDico [][]Neighbour

func (r * Run) InitNeighbourDico( bodies * ([]quadtree.Body)) {
	neighbours := make(NeighbourDico, len(*bodies))
	r.bodiesNeighbours = & neighbours
	for idx,_  := range *r.bodiesNeighbours {
		(*r.bodiesNeighbours)[idx] = make( []Neighbour, NbOfNeighboursPerBody)
	}
	r.bodiesNeighbours.Reset()

	neighboursOrig := make(NeighbourDico, len(*bodies))
	r.bodiesNeighboursOrig = & neighboursOrig
	for idx,_  := range *r.bodiesNeighboursOrig {
		(*r.bodiesNeighboursOrig)[idx] = make( []Neighbour, NbOfNeighboursPerBody)
	}
	r.bodiesNeighboursOrig.Reset()
}

// reset neighbour dico
func (dico * NeighbourDico) Reset() {

	for idx,_  := range *dico {
		for n, _ := range (*dico)[idx] {
			(*dico)[idx][n].n = nil
			(*dico)[idx][n].Distance = 2.0			
		}	
	}
}

// reset neighbour dico
func (dicoTarget * NeighbourDico) Copy(dicoSource * NeighbourDico) {

	for idx,_  := range *dicoSource {
		for n, _ := range (*dicoSource)[idx] {
			(*dicoTarget)[idx][n].n = (*dicoSource)[idx][n].n 
			(*dicoTarget)[idx][n].Distance = (*dicoSource)[idx][n].Distance 			
		}	
	}
}

// insert body at index in the dico according to the distance
func (dico * NeighbourDico) Insert(index int, body * quadtree.Body, distance float64) {

	// check if body is eligible to the last rank
	// replace if last rank is nil or if distance is greater 
	if( (*dico)[index][NbOfNeighboursPerBody-1].n == nil || 
		(*dico)[index][NbOfNeighboursPerBody-1].Distance > distance ) {

		(*dico)[index][NbOfNeighboursPerBody-1].n = body
		(*dico)[index][NbOfNeighboursPerBody-1].Distance = distance
	}

	// swap from last rank to rank 0
	for rank := NbOfNeighboursPerBody-2; rank >=0; rank-- {
		if( (*dico)[index][rank].n == nil || 
			(*dico)[index][rank].Distance > (*dico)[index][rank+1].Distance ) {

			tmp := (*dico)[index][rank]
			(*dico)[index][rank] = (*dico)[index][rank+1]
			(*dico)[index][rank+1] = tmp
		}
	}
}

// compute stirring
// parse all bodies and count the number of neighbors that are still 
// the neighbours at the origin
func (dico * NeighbourDico) ComputeStirring(dicoOrig * NeighbourDico) float64 {

	numberOfKeptNeighbors := 0 
	for idx,_  := range *dico {
		for n, _ := range (*dico)[idx] {

			body := (*dico)[idx][n].n
			// parse dico orig to check is the body is present
			for nOrig , _ := range (*dicoOrig)[idx] {
				if (*dicoOrig)[idx][nOrig].n == body {
					numberOfKeptNeighbors++
					continue
				}
			}
		}	
	}

	return float64(numberOfKeptNeighbors)/float64(len(*dico))
}