package quadtree

import (
	"testing"
	"math/rand"
)

func BenchmarkSetLevel(b * testing.B) {
	for i := 0; i<b.N;i++ {		var c Coord;	c.setLevel(6) }
}

func BenchmarkSetX(b * testing.B) {
	for i := 0; i<b.N;i++ {		var c Coord;	c.setX(6) }
}

func BenchmarkSetY(b * testing.B) {
	for i := 0; i<b.N;i++ {		var c Coord;	c.setY(6) }
}

func BenchmarkGetCoord8(b * testing.B) {
	for i := 0; i<b.N;i++ {		var b Body; b.getCoord8()}
}

func BenchmarkComputeLevel8(b * testing.B) {
	var q Quadtree
	var bodies []Body
		
	initQuadtree( &q, &bodies, 1000000)
	
	for i := 0; i<b.N;i++ {	q.updateNodesList( bodies)}
}

func BenchmarkUpdateNodesCOM(b * testing.B) {
	var q Quadtree
	var bodies []Body
		
	initQuadtree( &q, &bodies, 1000000)
	q.updateNodesList( bodies)
	
	for i := 0; i<b.N;i++ {	q.updateNodesCOM()}
}

func BenchmarkUpdateNodesAbove8(b * testing.B) {
	var q Quadtree
	var bodies []Body
		
	initQuadtree( &q, &bodies, 1000000)
	q.SetupNodeLinks()
	q.updateNodesList( bodies)
	
	for i := 0; i<b.N;i++ {	q.updateNodesAbove8()}
}

// init a quadtree with random position
func initQuadtree( q * Quadtree, bodies * []Body, nbBodies int) {
	
	// var q Quadtree
	*bodies = make([]Body, nbBodies)
	
	// init bodies
	for i := 0; i < nbBodies; i++ {
		(*bodies)[i].X = rand.Float64()
		(*bodies)[i].Y = rand.Float64()
		(*bodies)[i].M = rand.Float64()
	}
}

func BenchmarkInitQuadtree(b * testing.B) {
	for i := 0; i<b.N;i++ {
		var q Quadtree
		var bodies []Body
		initQuadtree( &q , &bodies, 1000000)
	}
}
