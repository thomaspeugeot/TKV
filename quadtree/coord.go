package quadtree

import (
	"bytes"
	"fmt"
)

//
// Coordinate system of a node
//
// Situation : most quadtree implementation are dynamic (the nodes can be created and deleted after the quadtree initialisation).
// This is an optimal solution if bodies are sparesely located (as in cosmology). Access to node is in (log(n))
//
// Current case is different because bodies are uniformly spread on a 2D square.
//
// Problem : the dynamic implementation is not necessary for uniformly spread bodies
//
// Solution : a static quadtree with the following requirements
//
//	Node coordinates are their rank in a direct access table.
//	Node coordinates of a body are computed directly from body's X,Y position
//
// Coordinates of a node are coded as follow
//
//	byte 0 : 0x00 (unused)
// 	byte 1 : level (root = 0, max depth = 7)
// 	byte 2 : X coordinate
// 	byte 3 : Y coordinate : coded on
//
// the X coordinate code depends on the level
//	level 0: there is no coordinate
//	level 1: west if 0, east is 128 (0x80)
//	level 2: node coordinates are 0, 64 (0x40), 128 (0x80), 192 (0x84)
//	...
//	level 8: node coordinates are encoded on the full 8 bits, from 0 to 0xFF (255)
type Coord uint32

// node level of a node coord c
// is between 0 and 8 and coded on 2nd byte of the Coord c
func (c Coord) Level() int { return int(c >> 16) }
func (c *Coord) SetLevel(level int) {

	*c = *c & 0x0000FFFF // reset level but bytes for x & y are preserved
	var pad uint32
	pad = (uint32(level) << 16)
	*c = *c | Coord(pad) // set level

}

// x coord
func (c Coord) X() int { return int((c & 0x0000FFFF) >> 8) }

// set X coordinate of node in Hexa from 0 to 255
func (c *Coord) setXHexaLevel8(x int) {

	*c = *c & 0x00FF00FF // reset x bytes

	var pad uint32
	pad = (uint32(x) << 8)

	*c = *c | Coord(pad)
}

// set X coordinate in Hexa according to level
// x is between 0 and 1<<(level-1)
func (c *Coord) setXHexa(x int, level int) {
	c.setXHexaLevel8(x << (8 - uint(level)))
}

// y coord
func (c Coord) Y() int { return int(c & 0x000000FF) }
func (c *Coord) setYHexaLevel8(y int) {
	*c = *c & 0x00FFFF00 // reset y bytes

	var pad uint32
	pad = uint32(y)
	*c = *c | Coord(pad)
}

// set Y coordinate in Hexa according to level
// y is between 0 and 1<<(level-1)
func (c *Coord) setYHexa(y int, level int) {
	c.setYHexaLevel8(y << (8 - uint(level)))
}

// check encoding of c
func (c *Coord) checkIntegrity() bool {

	//	byte 0 is null
	if res := 0xFF000000 & (*c); res != 0x00 {
		return false
	}

	// check level is below or equal to 8
	if c.Level() > 8 {
		return false
	}

	// check x coord is encoded acoording to the level
	if false {
		fmt.Printf("y (0xFF >> uint( SetLevel(%d))) %08b\n", c.Level(), 0xFF>>uint(c.Level()))
	}
	if (0xFF>>uint(c.Level()))&c.X() != 0x00 {
		return false
	}

	// check y coord
	if (0xFF>>uint(c.Level()))&c.Y() != 0x00 {
		return false
	}

	return true
}

func GetCoord(level, i, j int) Coord {
	var coord Coord
	coord.SetLevel(level)
	coord.setXHexa(i, level)
	coord.setYHexa(j, level)

	return coord
}

// print a coord
func (c *Coord) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')

	fmt.Fprintf(&buf, "|%8b|%8b|%8b|%8b| %8x",
		0x000000FF&(*c>>24),
		0x000000FF&(*c>>16),
		0x000000FF&(*c>>8),
		0x000000FF&*c,
		*c)

	buf.WriteByte('}')
	return buf.String()
}
