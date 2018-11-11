package quadtree

import (
	"fmt"
	"math/rand"
)

type BodyXY struct {
	X float64
	Y float64
}

// a body is a position & a mass
type Body struct {
	BodyXY
	M float64

	// coordinate in the quadtree
	coord Coord

	prev, next *Body
}

// bodies of a node are linked together
// Some quadtree use an alternative choice : store bodies of a node in a slice attached
// to the node. This alternative implies memory allocation which one tries to avoid.
// number of memory allocation are benchmaked with
//	go test -bench=BenchmarkUpdateNodesList_10M -benchmem
func (b *Body) Next() *Body { return b.next }

func (b *Body) Coord() Coord { return b.coord }

// get Node coordinates at level 8
func (b Body) getCoord8() Coord {
	var c Coord

	c.SetLevel(8)
	c.setXHexaLevel8(int(b.X * 256.0))
	c.setYHexaLevel8(int(b.Y * 256.0))

	if c.checkIntegrity() == false {
		s := fmt.Sprintf("getCoord8 invalid coord %s", c.String())
		panic(s)
	}
	return c
}

// init a quadtree with random position
func InitBodiesUniform(bodies *[]Body, nbBodies int) {

	// var q Quadtree
	*bodies = make([]Body, nbBodies)

	// init bodies
	for idx := range *bodies {
		(*bodies)[idx].X = rand.Float64()
		(*bodies)[idx].Y = rand.Float64()
		(*bodies)[idx].M = rand.Float64()
	}
}
