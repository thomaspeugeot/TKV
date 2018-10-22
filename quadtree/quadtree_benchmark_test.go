package quadtree

import (
	"testing"
)

func BenchmarkSetLevel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var c Coord
		c.SetLevel(6)
	}
}

func BenchmarkSetXHexaLevel8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var c Coord
		c.setXHexaLevel8(6)
	}
}

func BenchmarkSetXHexa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var c Coord
		c.setXHexa(6, 5)
	}
}

func BenchmarkSetYHexaLevel8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var c Coord
		c.setYHexaLevel8(6)
	}
}

func BenchmarkGetCoord8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b Body
		b.getCoord8()
	}
}

func BenchmarkUpdateNodesList_10M(b *testing.B) {
	var q Quadtree
	var bodies []Body

	InitBodiesUniform(&bodies, 10000000)
	q.Init(&bodies)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.updateNodesList()
	}
}

func BenchmarkUpdateNodesCOM_10M(b *testing.B) {
	var q Quadtree
	var bodies []Body

	InitBodiesUniform(&bodies, 10000000)
	q.Init(&bodies)
	q.updateNodesList()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.updateNodesCOM()
	}
}

func BenchmarkInitQuadtree(b *testing.B) {
	var bodies []Body

	for i := 0; i < b.N; i++ {
		InitBodiesUniform(&bodies, 1000000)
	}
}
