// Compact implementation of a static 2D quadtree
//
// Caution : Work In Progress
// 
// 1st Goal is to support the Barnes Hut implementation for spreading bodies on a 2D square.
// 
// 2nd Goal is to support more than 1 million bodies
//
// A quatree is a set of nodes holding bodies.
//
// This implementation put constraints on inputs :
//
//	Bodies's X,Y coordinates are float64 between 0 & 1
//	Quadtree architecture is static
//	Depth of nodes is limited to 8 (256 * 256 cells at the level 8)
//
// to see the doc
//
// 		godoc -goroot=$GOPATH -http=:8080
package quadtree

import (
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

// a body is a position & a mass
type Body struct {
	X float64
	Y float64
	M float64
}

// a node is a body (
type Node struct {
	// bodies of the node 
	//  at level 8, this is the list of bodies pertaining in the bounding box of the node
	// for level 7 to 0, this is the list of the four bodies of the nodes at the level below (or +1)
	Bodies [](*Body)
	Body // barycenter with mass of all the bodies of the node
}


// a Quadtree store Nodes. It is a an array with direct access to the Nodes with the Nodes coordinate
// see Coord
type Quadtree [16*256*256]Node

var optim bool

func init() {
	optim = true
	fmt.Printf("Size of quadtree %d\n", 8*256*256)
} 

// node level of a node coord c
// is between 0 and 8 and coded on 2nd byte of the Coord c
func (c Coord) getLevel() int { return int( c >> 16) }
func (c * Coord) setLevel(level int) { 
	
	*c = *c & 0x0000FFFF // reset level but bytes for x & y are preserved
	var pad uint32
	pad = (uint32(level) << 16) 
	*c = *c | Coord(pad) // set level
	
}

// x coord
func (c Coord) getX() int { return int((c & 0x0000FFFF) >> 8) }
func (c * Coord) setX(x int) { 
	// fmt.Printf( "SetX c before reset x %8x\n", *c)

	*c = *c & 0x00FF00FF // reset x bytes
	// fmt.Printf( "SetX c after reset x %8x\n", *c)
	
	var pad uint32
	pad = (uint32(x) << 8) 
	// fmt.Printf( "SetX pad %8x\n", pad)


	*c = *c | Coord(pad)

	// fmt.Printf( "SetX c after x input %8x\n", *c)
	// if !checkIntegrity( *c) { panic("set X failed")}
}

// y coord
func (c Coord) getY() int { return int( c & 0x000000FF) }
func (c * Coord) setY(y int) { 
	*c = *c & 0x00FFFF00 // reset y bytes
	
	var pad uint32
	pad = uint32(y) 
	*c = *c | Coord(pad)
}



// get Node coordinates at level 8
func ( b Body) getCoord8() Coord {
	var c Coord
	
	c.setLevel( 8)
	c.setX( int(b.X * 256.0) )
	c.setY( int(b.Y * 256.0) )
	return c
}

// check encoding of c
func checkIntegrity( c Coord) bool {
	
	//	byte 0 is null
	if res := 0xFF000000 & c; res != 0x00 {
		return false
	}
	
	// check level is below or equal to 8
	if c.getLevel() > 8 {
		return false
	}
	
	// check x coord is encoded acoording to the level
	if (false) { fmt.Printf( "y (0xFF >> uint( setLevel(%d))) %08b\n", c.getLevel(), 0xFF >> uint( c.getLevel())) }
	if (0xFF >> uint( c.getLevel())) & c.getX() != 0x00 {
		return false
	}

	// check y coord
	if (0xFF >> uint( c.getLevel())) & c.getY() != 0x00 {
		return false
	}
	
	return true
}

type Direction uint
const (
	NW = 0x0000
	NE = 0x0100
	SW = 0x0001
	SE = 0x0101
)

// compute quadtree Nodes for levels from 0 to 7
func (q * Quadtree) updateNodesAbove8() {
	
	for level := 7; level >= 0; level-- {
	
		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)
		
		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {
				
				coord := Coord( uint(level)<<16 | uint(i)<<8 | uint(j))
				node := q[coord]
				node.updateNode()
			}
		}
	}
}

// setup quadtree Nodes for levels from 0 to 7
func (q * Quadtree) setupNodeLinks() {
	
	for level := 7; level >= 0; level-- {
	
		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)

		// level below has a higher number (this goes against elevator common sense)
		levelBelow := level+1
		
		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {
				
				coord := Coord( uint(level)<<16 | uint(i)<<8 | uint(j))
				node := q[coord]
				node.Bodies = make([]*Body, 4)
				shift := uint( 8-levelBelow)
				
				// to go east at the level below, we flip to 1 the bit that is significant at that level 
				coordNW := Coord( uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | NW << shift)
				coordNE := Coord( uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | NE << shift)
				coordSW := Coord( uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | SW << shift)
				coordSE := Coord( uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | SE << shift)
				
				nodeNW := q[coordNW]
				nodeNE := q[coordNE]
				nodeSW := q[coordSW]
				nodeSE := q[coordSE]
				
				node.Bodies[0] = & nodeNW.Body
				node.Bodies[1] = & nodeNE.Body
				node.Bodies[2] = & nodeSW.Body
				node.Bodies[3] = & nodeSE.Body
			}
		}
	}
}

// reset a Quadtree level 8
// clear all nodes bodies & set masse to 0
func (q * Quadtree) resetLevel8() {
	
	for _, n := range q {
		//		if i >= (1<<16) { return }
		n.Bodies = make([]*Body, 0) // this put the slice into garbage collection
	}
}

// fill quadtree at level with bodies 
func (q * Quadtree) computeLevel8 (bodies []Body) {

	q.resetLevel8()
	for _, b := range bodies {
	
		// get coord of body (this is direct access)
		coord := b.getCoord8()
	
		if (false) { fmt.Printf("computeLevel8 coordinate %d\n", coord) }
		
		nodeBodies := q[coord].Bodies
		nodeBodies = append( nodeBodies, &b)
	}		
}

// compute COM of quadtree at level 8 
func (q * Quadtree) computeCOMAtLevel8 () {

	q.resetLevel8()
	for _, n := range q {
		
		n.M = 0.0
		// get all bodies of the node
		for _, b	:= range n.Bodies {
			n.M += b.M
			n.X += b.X*b.M
			n.Y += b.Y*b.M
		}	
		// divide by total mass to get the barycenter
		if n.M > 0 {
			n.X /= n.M
			n.Y /= n.M
		}
	}		
}

func (n * Node) updateNode() {
	
	n.M = 0.0
	// get all bodies of the node
	for _, b	:= range n.Bodies {
		n.M += b.M
		n.X += b.X*b.M
		n.Y += b.Y*b.M
	}	
	
	// divide by total mass to get the barycenter
	if n.M > 0 {
		n.X /= n.M
		n.Y /= n.M
	}
}