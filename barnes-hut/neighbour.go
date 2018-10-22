package barneshut

import (
	"github.com/thomaspeugeot/tkv/quadtree"
)

// relative to a body of interest, the storage for a neighbour with its distance
// nota : this is used to measure the stirring of the bodies along the simulation
type Neighbour struct {
	n        *quadtree.Body // rank in the []quadtree.Body
	Distance float64
}

// the measure of stirring is computed with a finite number of neighbours
// no stirring is that the neighbours at the end of the run are the same neighbours
// that at the beginning
var NbOfNeighboursPerBody int = 10

type NeighbourDico [][]Neighbour

func (r *Run) InitNeighbourDico(bodies *([]quadtree.Body)) {
	neighbours := make(NeighbourDico, len(*bodies))
	r.bodiesNeighbours = &neighbours
	for idx := range *r.bodiesNeighbours {
		(*r.bodiesNeighbours)[idx] = make([]Neighbour, NbOfNeighboursPerBody)
	}
	r.bodiesNeighbours.Reset()

	neighboursOrig := make(NeighbourDico, len(*bodies))
	r.bodiesNeighboursOrig = &neighboursOrig
	for idx := range *r.bodiesNeighboursOrig {
		(*r.bodiesNeighboursOrig)[idx] = make([]Neighbour, NbOfNeighboursPerBody)
	}
	r.bodiesNeighboursOrig.Reset()
}

// reset neighbour dico
func (dico *NeighbourDico) Reset() {

	for idx := range *dico {
		for n := range (*dico)[idx] {
			(*dico)[idx][n].n = nil
			(*dico)[idx][n].Distance = 2.0
		}
	}
}

// reset neighbour dico
func (dicoTarget *NeighbourDico) Copy(dicoSource *NeighbourDico) {

	for idx := range *dicoSource {
		for n := range (*dicoSource)[idx] {
			(*dicoTarget)[idx][n].n = (*dicoSource)[idx][n].n
			(*dicoTarget)[idx][n].Distance = (*dicoSource)[idx][n].Distance
		}
	}
}

// insert body at index in the dico according to the distance
func (dico *NeighbourDico) Insert(index int, body *quadtree.Body, distance float64) {

	if index == 0 {
		Trace.Printf("Insert begin with distance %f", distance)
	}

	//check if body is already present
	// if yes, update distance
	for rank := NbOfNeighboursPerBody - 1; rank >= 0; rank-- {
		if (*dico)[index][rank].n == body {
			if (*dico)[index][rank].Distance > distance {
				(*dico)[index][rank].Distance = distance
			}
			return
		}
	}

	// check if body is eligible to the last rank
	// replace if last rank is nil or if distance is greater
	if (*dico)[index][NbOfNeighboursPerBody-1].n == nil ||
		(*dico)[index][NbOfNeighboursPerBody-1].Distance > distance {

		(*dico)[index][NbOfNeighboursPerBody-1].n = body
		(*dico)[index][NbOfNeighboursPerBody-1].Distance = distance
	}

	// swap from last rank to rank 0
	for rank := NbOfNeighboursPerBody - 2; rank >= 0; rank-- {
		if (*dico)[index][rank].n == nil ||
			(*dico)[index][rank].Distance > (*dico)[index][rank+1].Distance {

			tmp := (*dico)[index][rank]
			(*dico)[index][rank] = (*dico)[index][rank+1]
			(*dico)[index][rank+1] = tmp
		}
	}
}

// check integrity
// not twice the same neighbor
func (dico *NeighbourDico) Check() {

	for idx := range *dico {
		for n := range (*dico)[idx] {

			body := (*dico)[idx][n].n
			if body == nil {
				Info.Printf("nil neighbour at index %d rank %d, with distance %f", idx, n, (*dico)[idx][n].Distance)
			}

			// check to see if neighbor is present twice
			for nOrig := range (*dico)[idx] {
				if (*dico)[idx][nOrig].n == body && nOrig != n {

					Info.Printf("neighbour found twice at index %d rank %d and rank %d", idx, n, nOrig)

				}
			}
		}
	}
}

func (dico *NeighbourDico) ComputeRatioOfNilNeighbours() float64 {

	nbOfNil := 0
	for idx := range *dico {
		for n := range (*dico)[idx] {

			if (*dico)[idx][n].n == nil {
				nbOfNil++
			}
		}
	}
	return float64(nbOfNil) / float64(len(*dico))
}

// compute stirring
// parse all bodies and count the number of neighbors that are still
// the neighbours at the origin
func (dico *NeighbourDico) ComputeStirring(dicoOrig *NeighbourDico) float64 {

	numberOfNeighbors := 0 // number of neighbours in orig that are not nil
	numberOfKeptNeighbors := 0
	for idx := range *dico {
		for nOrig := range (*dicoOrig)[idx] {

			bodyOrig := (*dicoOrig)[idx][nOrig].n

			if bodyOrig != nil {
				numberOfNeighbors++

				// parse dico to check if the orig neighbour's body is present
				for n := range (*dico)[idx] {
					found := 0
					if (*dico)[idx][n].n == bodyOrig {
						numberOfKeptNeighbors++
						found++
						// continue
					}
					if found > 1 {
						Info.Printf("At index %d, Neighbours %d found twice", idx, nOrig)
					}
				}
			}
		}
	}

	Trace.Printf("numberOfNeighbors %d numberOfKeptNeighbors %d", numberOfNeighbors, numberOfKeptNeighbors)
	return float64(numberOfKeptNeighbors) / float64(numberOfNeighbors)
}
