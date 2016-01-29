package quadtree

import (
	"testing"
	"math/rand"
)

func BenchmarkSetLevel(b * testing.B) {
	for i := 0; i<b.N;i++ {		var c Coord;	c.setLevel(6) }
}

func BenchmarkSetXHexaLevel8(b * testing.B) {
	for i := 0; i<b.N;i++ {		var c Coord;	c.setXHexaLevel8(6) }
}

func BenchmarkSetXHexa(b * testing.B) {
	for i := 0; i<b.N;i++ {		var c Coord;	c.setXHexa(6, 5) }
}

func BenchmarkSetYHexaLevel8(b * testing.B) {
	for i := 0; i<b.N;i++ {		var c Coord;	c.setYHexaLevel8(6) }
}

func BenchmarkGetCoord8(b * testing.B) {
	for i := 0; i<b.N;i++ {		var b Body; b.getCoord8()}
}

func BenchmarkUpdateNodesList_10M(b * testing.B) {
	var q Quadtree
	var bodies []Body
		
	initBodies( &bodies, 10000000)
	b.ResetTimer()
	
	for i := 0; i<b.N;i++ {	q.updateNodesList( bodies)}
}

func BenchmarkUpdateNodesCOM_10M(b * testing.B) {
	var q Quadtree
	var bodies []Body
		
	initBodies( &bodies, 10000000)
	q.updateNodesList( bodies)
	
	b.ResetTimer()
	
	for i := 0; i<b.N;i++ {	q.updateNodesCOM()}
}

// init a quadtree with random position
func initBodies( bodies * []Body, nbBodies int) {
	
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
	var bodies []Body


	for i := 0; i<b.N;i++ {
		initBodies( &bodies, 1000000)
	}
}
