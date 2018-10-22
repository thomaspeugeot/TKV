package barneshut

import (
	"fmt"
	"math"
	"sort"
)

type MaxRepulsiveForce struct {
	AccX, AccY, X, Y float64
	Norm             float64 // direction and norm
	Idx              int     // body index where repulsive vector is max
}

func (r *Run) ComputeMaxRepulsiveForce() {
	r.maxRepulsiveForce.Norm = 0.0

	// parse bodies
	for idx := range *r.bodies {
		acc := &((*r.bodiesAccel)[idx])
		norm := math.Sqrt(acc.X*acc.X + acc.Y*acc.Y)
		if norm > r.maxRepulsiveForce.Norm {
			r.maxRepulsiveForce.AccX = acc.X
			r.maxRepulsiveForce.AccY = acc.Y
			r.maxRepulsiveForce.Idx = idx
			r.maxRepulsiveForce.Norm = norm
			r.maxRepulsiveForce.X = (*r.bodies)[idx].X
			r.maxRepulsiveForce.Y = (*r.bodies)[idx].Y
		}
	}
}

// compute the density per village and return the density per village
func (r *Run) ComputeDensityTencilePerTerritoryString() [10]string {
	var densityString [10]string
	density := r.ComputeDensityTencilePerTerritory()
	for tencile := range density {
		densityString[tencile] = fmt.Sprintf("%3.2f", density[tencile])
	}
	return densityString
}
func (r *Run) ComputeDensityTencilePerTerritory() [10]float64 {

	// parse all bodies
	// prepare the village
	villages := make([][]int, nbVillagePerAxe)
	for x := range villages {
		villages[x] = make([]int, nbVillagePerAxe)
	}

	// parse bodies
	for _, b := range *r.bodies {
		// compute village coordinate (from 0 to nbVillagePerAxe-1)
		x := int(math.Floor(float64(nbVillagePerAxe) * b.X))
		y := int(math.Floor(float64(nbVillagePerAxe) * b.Y))

		villages[x][y]++
	}

	// var bodyCount []int
	nbVillages := nbVillagePerAxe * nbVillagePerAxe
	bodyCountPerVillage := make([]int, nbVillages)
	for x := range villages {
		for y := range villages[x] {
			bodyCountPerVillage[y+x*nbVillagePerAxe] = villages[x][y]
		}
	}
	sort.Ints(bodyCountPerVillage)

	var density [10]float64
	for tencile := range density {
		lowIndex := int(math.Floor(float64(nbVillages) * float64(tencile) / 10.0))
		highIndex := int(math.Floor(float64(nbVillages) * float64(tencile+1) / 10.0))
		// log.Output( 1, fmt.Sprintf( "tencile %d ", tencile))
		// log.Output( 1, fmt.Sprintf( "lowIndex %d ", lowIndex))
		// log.Output( 1, fmt.Sprintf( "highIndex %d ", highIndex))

		nbBodiesInTencile := 0
		for _, nbBodies := range bodyCountPerVillage[lowIndex:highIndex] {
			nbBodiesInTencile += nbBodies
		}
		density[tencile] = float64(nbBodiesInTencile) / float64(len(bodyCountPerVillage[lowIndex:highIndex]))

		// we compare with then average bodies per villages
		density[tencile] /= float64(len(*r.bodies)) / float64(nbVillages)

		// we round the density to 0.01 precision, and put it in percentage point
		density[tencile] *= 100.0 * 100.0
		intDensity := math.Floor(density[tencile])
		density[tencile] = float64(intDensity) / 100.0
	}

	return density
}
