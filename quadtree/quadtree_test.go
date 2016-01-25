package quadtree

import (
	"testing"
	"math/rand"
)


func TestLevel(t *testing.T) {

	cases := []struct {
		in Coord
		want int
	}{
		{0, 0},
		{1<<16, 1},
		{1<<16 + 34, 1},
		{21<<16 + 1 <<12, 21},
	}
		
	for _, c := range cases {
		got := c.in.getLevel()
		if( got != c.want) {
			t.Errorf("getLevel(%32b) == %d, want %d", c.in, got, c.want)
			t.Errorf("getLevel(|%8b|%8b|%8b|%8b|) == %d, want %d", 
					0x000000FF & (c.in >> 24), 0x000000FF & (c.in >> 16), 0x000000FF & (c.in >> 8), 0x000000FF & c.in, 
					got, c.want)
		}	
	}
}



func TestIntegrity( t *testing.T) {
	
	cases := []struct {
		in Coord
		want bool
	}{
		{ 0x00, true},
		{ 0x000000FF, false}, // at level 0, no bits are allowed for x or y
		{ 0x000a0001, false}, // level a is above 8
		{ 0x00070010, true}, // level 7 is OK
		{ 0x00070001, false}, // the last bit shall be 0
		{ 0x00070101, false}, // the last bit of x shall be 0
		{ 0x00080001, true}, // the last bit can be 1
		{ 0x000A0000, false},
		{ 0x0A0A0000, false}, // byt0 shall be null
	}
	for rank, c := range cases {
		got := checkIntegrity(c.in)
		if( got != c.want) {
			// t.Errorf("checkIntegrity(%b) == %t, want %t", c.in, got, c.want)
			t.Errorf("case %d - checkIntegrity of |%8b|%8b|%8b|%8b|, %8x, level %d == %t, want %t", 
					rank,
					0x000000FF & (c.in >> 24), 0x000000FF & (c.in >> 16), 0x000000FF & (c.in >> 8), 0x000000FF & c.in, 
					c.in, c.in.getLevel(),
					got, c.want)
		}	
	}
}

// init a quadtree with random position
//
func initQuadtree( qa * Quadtree, nbBodies int) {
	
	var q Quadtree
	b := make([]Body, nbBodies)
	
	// init bodies
	for i := 0; i < nbBodies; i++ {
		b[i].X = rand.Float64()
		b[i].Y = rand.Float64()
		b[i].M = rand.Float64()
	}
	
	//
	qa = &q
	
}

func BenchmarkInitQuadtree(b * testing.B) {
	for i := 0; i<b.N;i++ {
		var q Quadtree
		initQuadtree( &q , 1000000)
	}

}

func BenchmarkSetLevel(b * testing.B) {
	for i := 0; i<b.N;i++ {		var c Coord;	c.setLevel(6) }
}

func BenchmarkSetX(b * testing.B) {
	for i := 0; i<b.N;i++ {		var c Coord;	c.setX(6) }
}

func BenchmarkSetY(b * testing.B) {
	for i := 0; i<b.N;i++ {		var c Coord;	c.setY(6) }
}

func TestGetCoord8(t * testing.T) {
	cases := []struct {
		in Body
		want Coord
	}{
		{ Body{0.0, 0.0, 0.0}, 0x00080000 }, 
	}
	for _, c := range cases {
		got := c.in.getCoord8()
		if( got != c.want) {
			t.Errorf("getCoord8(%#v) == |%8b|%8b|%8b|%8b|, %8x, want %8x", c.in, 
			0x000000FF & (got >> 24), 0x000000FF & (got >> 16), 
			0x000000FF & (got >> 8), 0x000000FF & got, got, c.want) 
		}
	}
}

func TesSet(t * testing.T) {
	cases := []struct {
		inCoord Coord
		inLevel int
		inCoords PosSuite
		want bool
	}{
		{	0, 			2, 			PosSuite{1, 2}, 			true}, // level is OK
		{	0, 			2, 			PosSuite{1, 4}, 			false}, // one pos is above 3
		{	0, 			2, 			PosSuite{1, -1}, 			false}, // one pos is below 0
		{	0, 			1, 			PosSuite{1, 2}, 			false}, // level is not good
 	}
		
	for _, c := range cases {
		got := set( &c.inCoord, c.inLevel, c.inCoords)
		if( got != c.want) {
			t.Errorf("Set(%b, %d, %q) == %t, want %t", c.inCoord, c.inLevel, c.inCoords, got, c.want)
		}	
	}
}