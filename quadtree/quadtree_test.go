package quadtree

import (
	"testing"
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
		got := Level(c.in)
		if( got != c.want) {
			t.Errorf("Level(%32b) == %d, want %d", c.in, got, c.want)
			t.Errorf("Level(|%8b|%8b|%8b|%8b|) == %d, want %d", 
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
		{ 0x000000FF, false},
		{ 0x000a0001, false},
		{ 0x00070001, true},
		{ 0x000A0000, false},
		{ 0x0A0A0000, false},
	}
	for rank, c := range cases {
		got := checkIntegrity(c.in)
		if( got != c.want) {
			// t.Errorf("checkIntegrity(%b) == %t, want %t", c.in, got, c.want)
			t.Errorf("case %d - checkIntegrity of |%8b|%8b|%8b|%8b|, %8x, level %d == %t, want %t", 
					rank,
					0x000000FF & (c.in >> 24), 0x000000FF & (c.in >> 16), 0x000000FF & (c.in >> 8), 0x000000FF & c.in, 
					c.in, Level( c.in),
					got, c.want)
		}	
	}
}

func TestSet(t * testing.T) {

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