// this the file for the test (benchmarks are in another file)
package quadtree

import (
	"fmt"
	"testing"
)

// test function Level()
func TestLevel(t *testing.T) {

	cases := []struct {
		in   Coord
		want int
	}{
		{0, 0},
		{1 << 16, 1},
		{1<<16 + 34, 1},
		{21<<16 + 1<<12, 21},
	}

	for _, c := range cases {
		got := c.in.Level()
		if got != c.want {
			t.Errorf("Level(%32b) == %d, want %d", c.in, got, c.want)
			t.Errorf("Level(|%8b|%8b|%8b|%8b|) == %d, want %d",
				0x000000FF&(c.in>>24), 0x000000FF&(c.in>>16), 0x000000FF&(c.in>>8), 0x000000FF&c.in,
				got, c.want)
		}
	}
}

// test function "check integrity" of Coords
func TestCoordIntegrity(t *testing.T) {

	cases := []struct {
		in   Coord
		want bool
	}{
		{0x00, true},
		{0x000000FF, false}, // at level 0, no bits are allowed for x or y
		{0x000a0001, false}, // level a is above 8
		{0x00070010, true},  // level 7 is OK
		{0x00070001, false}, // the last bit shall be 0
		{0x00070101, false}, // the last bit of x shall be 0
		{0x00080001, true},  // the last bit can be 1
		{0x000A0000, false},
		{0x0A0A0000, false}, // byte 0 shall be null
	}
	for rank, c := range cases {
		got := c.in.checkIntegrity()
		if got != c.want {
			// t.Errorf("checkIntegrity(%b) == %t, want %t", c.in, got, c.want)
			t.Errorf("case %d - checkIntegrity of |%8b|%8b|%8b|%8b|, %8x, level %d == %t, want %t",
				rank,
				0x000000FF&(c.in>>24), 0x000000FF&(c.in>>16), 0x000000FF&(c.in>>8), 0x000000FF&c.in,
				c.in, c.in.Level(),
				got, c.want)
		}
	}
}

func TestSetXHexaLevel8(t *testing.T) {
	cases := []struct {
		in   Coord
		inX  int
		want Coord
	}{
		{0x00080000, 8, 0x00080800},
		{0x00080001, 8, 0x00080801},
	}
	for _, c := range cases {
		coord := c.in
		coord.setXHexaLevel8(c.inX)
		got := coord
		if coord != c.want {
			t.Errorf("\n setXHexa(%d)\nin   %s\ngot  %s, \nwant %s",
				c.inX, &c.in, &got, &c.want)
		}
	}
}

func TestCoordString(t *testing.T) {

	c := Coord(0x00080A0B)

	want := "{|       0|    1000|    1010|    1011|    80a0b}"

	got := fmt.Sprintf("%s", &c)

	if got != want {
		t.Errorf("\ngot  %s\nwant %s", got, want)
	}
}

func TestSetYHexaLevel8(t *testing.T) {
	cases := []struct {
		in   Coord
		inY  int
		want Coord
	}{
		{0x00080000, 8, 0x00080008},
	}
	for _, c := range cases {
		got := c.in
		got.setYHexaLevel8(c.inY)
		if got != c.want {
			t.Errorf("%#v setYHexa(%d) == |%8b|%8b|%8b|%8b|, got %8x, want %8x", c.in, c.inY,
				0x000000FF&(got>>24), 0x000000FF&(got>>16),
				0x000000FF&(got>>8), 0x000000FF&got, got, c.want)
		}
	}
}

func TestUpdateNodesList(t *testing.T) {

	var q Quadtree
	var bodies []Body

	// fmt.Printf("TestUpdateNodesList before initQuadtree\n")
	InitBodiesUniform(&bodies, 1000000)
	q.Init(&bodies)

	var coord Coord = 0x0007546e
	if q.Nodes[coord].coord != coord {
		t.Errorf("coord not set up want %s, got %s", coord.String(), q.Nodes[coord].coord.String())
	}

	coord = 0x00069430
	if q.Nodes[coord].coord != coord {
		t.Errorf("coord not set up want %s, got %s", coord.String(), q.Nodes[coord].coord.String())
	}

	// fmt.Printf("TestUpdateNodesList before updateNodesList\n")
	q.updateNodesList()
	q.CheckIntegrity(t)
	q.updateNodesList()
	q.CheckIntegrity(t)
	q.updateNodesList()
	q.CheckIntegrity(t)
}

// check computation of nodes below
func TestNodesBelow(t *testing.T) {
	// var q Quadtree
	var c Coord
	c.SetLevel(1)
	c.setXHexa(1, 1)
	c.setYHexa(1, 1)

	if !c.checkIntegrity() {
		t.Errorf("invalid input %s", &c)
	}

	// coordNW, coordNE, coordSW, coordSE := NodesBelow( c)
	// fmt.Printf("\nTestNodesBelow\nin %s\nnw %s\nne %s\nsw %s\nse %s", &c, &coordNW, &coordNE, &coordSW, &coordSE)

	// n_coordNW, coordNE, coordSW, coordSE := NodesBelow( coordNW)
	// fmt.Printf("\nTestNodesBelow\n\nin %s\nnw %s\nne %s\nsw %s\nse %s", &coordNW, &n_coordNW, &coordNE, &coordSW, &coordSE)
}

func TestUpdateNodesCOM(t *testing.T) {

	var q Quadtree

	var bodies []Body

	// fmt.Printf("TestUpdateNodesCOM before initQuadtree\n")
	InitBodiesUniform(&bodies, 10000)

	// fmt.Printf("TestUpdateNodesCOM before updateNodesList\n")
	q.Init(&bodies)

	q.updateNodesList()
	q.updateNodesCOM()

	q.CheckIntegrity(t)
}

func TestComputeGini(t *testing.T) {

	var q Quadtree
	var bodies []Body
	InitBodiesUniform(&bodies, 1000000)

	// fmt.Printf("TestUpdateNodesCOM before updateNodesList\n")
	q.Init(&bodies)

	q.updateNodesList()
	q.updateNodesCOM()
	q.CheckIntegrity(t)

	q.ComputeQuadtreeGini()
	ratioPoor := float32(q.BodyCountGini[8][0]) / float32(len(bodies))
	if ratioPoor < 0.001 {
		t.Errorf("too low value for poor tencile got : %f, want %f", ratioPoor, 0.001)
	}
	if q.BodyCountGini[8][0] == 0 {
		t.Errorf("zero value for poor tencile")
	}

}

func TestGetCoord8(t *testing.T) {
	cases := []struct {
		in   Body
		want Coord
	}{
		{Body{BodyXY{0.0, 0.0}, 0.0, 0x0, nil, nil}, 0x00080000},
		{Body{BodyXY{0.0, 255.999}, 255.999, 0x0, nil, nil}, 0x0008FFFF},
	}
	for _, c := range cases {
		got := c.in.getCoord8()
		if got != c.want {
			t.Errorf("getCoord8(%#v) == |%8b|%8b|%8b|%8b|, %8x, want %8x", c.in,
				0x000000FF&(got>>24), 0x000000FF&(got>>16),
				0x000000FF&(got>>8), 0x000000FF&got, got, c.want)
		}
	}
}

func TestSetupNodesLinks(t *testing.T) {
	var q Quadtree

	// t.Errorf("TestSetupNodesLinks")
	var c Coord
	c.SetLevel(1)
	c.setXHexaLevel8(128)
	c.setYHexaLevel8(128)

	// coordNW, coordNE, coordSW, coordSE := NodesBelow( c)

	q.setupNodesLinks()
	n := &(q.Nodes[c])
	b := n.first

	// fmt.Printf("%8x\n", c)
	if b == nil {
		t.Errorf("first is nil")
	} else if b.next == nil {
		t.Errorf("first body of node has no next")
	}
}
