// Compact implementation of a static 2D quadtree
//
// Caution : Work In Progress
// 
// 1st Goal is to support a Barnes Hut algorithm implementation. 
// This variation of the BH is not for cosmology but for the problem of bodies on a 2D square that you want 
// to put the most apart (like dancers in a crowded night club).
// 
// 2nd Goal is to support more than 1 million bodies
//
// A quatree is a hierarchical set of nodes that divide the 2D space. 
// Each node holds the bodies that are located in its area.
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
	"bytes"
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
	
	// bodies of a node are linked together
	// Some quadtree use an alternative choice : store bodies of a node in a slice attached
	// to the node. This alternative implies memory allocation which one tries to avoid.
	// number of memory allocation are benchmaked with 
	//	go test -bench=BenchmarkUpdateNodesList_10M -benchmem
	prev, next * Body 
	
}

// a node is a body (
type Node struct {
	// bodies of the node 
	//  at level 8, this is the list of bodies pertaining in the bounding box of the node
	// for level 7 to 0, this is the list of the four bodies of the nodes at the level below (or +1)
	
	// Barycenter with mass of all the bodies of the node
	// this body is linked with the bodies at his level in the node
	Body 
	first * Body  // link to the bodies below
}


// a Quadtree store Nodes. It is a an array with direct access to the Nodes with the Nodes coordinate
// see Coord
type Quadtree [1<<20]Node

var optim bool

func init() {
	optim = true
} 

// node level of a node coord c
// is between 0 and 8 and coded on 2nd byte of the Coord c
func (c Coord) Level() int { return int( c >> 16) }
func (c * Coord) setLevel(level int) { 
	
	*c = *c & 0x0000FFFF // reset level but bytes for x & y are preserved
	var pad uint32
	pad = (uint32(level) << 16) 
	*c = *c | Coord(pad) // set level
	
}

// x coord
func (c Coord) X() int { return int((c & 0x0000FFFF) >> 8) }
// set X coordinate of node in Hexa from 0 to 255
func (c * Coord) setXHexaLevel8(x int) { 

	*c = *c & 0x00FF00FF // reset x bytes
	
	var pad uint32
	pad = (uint32(x) << 8) 

	*c = *c | Coord(pad)
}
// set X coordinate in Hexa according to level
// x is between 0 and 1<<(level-1)
func (c * Coord) setXHexa(x int, level int) {
	c.setXHexaLevel8( x << (8- uint(level)))
}

// y coord
func (c Coord) Y() int { return int( c & 0x000000FF) }
func (c * Coord) setYHexaLevel8(y int) { 
	*c = *c & 0x00FFFF00 // reset y bytes
	
	var pad uint32
	pad = uint32(y) 
	*c = *c | Coord(pad)
}
// set Y coordinate in Hexa according to level
// y is between 0 and 1<<(level-1)
func (c * Coord) setYHexa(y int, level int) {
	c.setYHexaLevel8( y << (8-uint(level)))
}

// get Node coordinates at level 8
func ( b Body) getCoord8() Coord {
	var c Coord
	
	c.setLevel( 8)
	c.setXHexaLevel8( int(b.X * 256.0) )
	c.setYHexaLevel8( int(b.Y * 256.0) )
	return c
}

// check encoding of c
func ( c * Coord) checkIntegrity() bool {
	
	//	byte 0 is null
	if res := 0xFF000000 & (*c); res != 0x00 {
		return false
	}
	
	// check level is below or equal to 8
	if c.Level() > 8 {
		return false
	}
	
	// check x coord is encoded acoording to the level
	if (false) { fmt.Printf( "y (0xFF >> uint( setLevel(%d))) %08b\n", c.Level(), 0xFF >> uint( c.Level())) }
	if (0xFF >> uint( c.Level())) & c.X() != 0x00 {
		return false
	}

	// check y coord
	if (0xFF >> uint( c.Level())) & c.Y() != 0x00 {
		return false
	}
	
	return true
}

// constants used to navigate from one node to the other
const (
	NW = 0x0000
	NE = 0x0100
	SW = 0x0001
	SE = 0x0101
)

// compute quadtree Nodes for levels from 0 to 7
func (q * Quadtree) updateNodesCOMAbove8() {
	
	for level := 7; level >= 0; level-- {
	
		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)
		
		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {
				
				var coord Coord 
				coord.setLevel( level)
				coord.setXHexa(i, level)
				coord.setYHexa(j, level)
				
				node := &(q[coord])
				node.updateCOM()
			}
		}
	}
}

// get nodes coords below
func (q * Quadtree) NodesBelow(c Coord)  (coordNW, coordNE, coordSW, coordSE Coord) {

	levelBelow := c.Level() + 1
	i := c.X()
	j := c.Y()
	shift := uint( 8-levelBelow)

	// to go east at the level below, we flip to 1 the bit that is significant at that level 
	coordNW = Coord( uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | NW << shift)
	coordNE = Coord( uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | NE << shift)
	coordSW = Coord( uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | SW << shift)
	coordSE = Coord( uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | SE << shift)
	
	return coordNW, coordNE, coordSW, coordSE
}

// print a coord
func (c * Coord) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
		
	fmt.Fprintf(&buf, "|%8b|%8b|%8b|%8b| %8x",
		0x000000FF & (*c >> 24), 
		0x000000FF & (*c >> 16), 
		0x000000FF & (*c >> 8), 
		0x000000FF & *c, 
		*c)
	
	buf.WriteByte('}')
	return buf.String()

}

// setup quadtree Nodes for levels from 7 to 0
func (q * Quadtree) SetupNodesLinks() {
	
	for level := 7; level >= 0; level-- {
	
		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)

		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {
				
				var coord Coord 
				coord.setLevel( level)
				coord.setXHexa(i, level)
				coord.setYHexa(j, level)
				
				node := &(q[coord])
				coordNW, coordNE, coordSW, coordSE := q.NodesBelow(coord)
				
				nodeNW := &q[coordNW]
				nodeNE := &q[coordNE]
				nodeSW := &q[coordSW]
				nodeSE := &q[coordSE]
	
				// bodies of the nodes below are chained
				// fmt.Printf("%8x\n", coord)
				node.first = & (nodeNW.Body)
				nodeNW.Body.next = & (nodeNE.Body)
				nodeNE.Body.next = & (nodeSW.Body)
				nodeSW.Body.next = & (nodeSE.Body)
			}
		}
	}
}

// fill quadtree at level with bodies 
func (q * Quadtree) updateNodesList (bodies * []Body) {

	for idx, _ := range (*bodies) {
	
		b := &((*bodies)[idx])
	
		// link the next body to the previous one
		if( b.next != nil) {
			b.next.prev = b.prev
		}
		// link the previous body to the next one
		if (b.prev != nil) {
			b.prev = b.next
		}
		
		// get coord of body (this is direct access)
		coord := b.getCoord8()

		// put body as the first body of the node
		// shift the first body if it is already there
		if( q[coord].first != nil) {
			// double link body to the current node's first
			b.next = q[coord].first
			q[coord].first.prev = b
		}
		
		// body b is the new node's first
		q[coord].first = b
		b.prev = q[coord].first
	}		
}

// compute COM of quadtree from level 8 to level 0
func (q * Quadtree) updateNodesCOM () {

	// compute is bottom up
	for level := 8; level >= 0; level-- {
	
		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)
		
		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {
				
				var coord Coord 
				coord.setLevel(level)
				coord.setXHexa(i, level)
				coord.setYHexa(j, level)
				
				node := &(q[coord])
				node.updateCOM()
			}
		}
	}	
}

func (q * Quadtree) updateNodesListsAndCOM (bodies * []Body) {

	q.updateNodesList( bodies)
	q.updateNodesCOM()
}


// update COM of a node (reset the current COM before)
func (n * Node) updateCOM() {
	
	n.M = 0.0
	n.X = 0.0
	n.Y = 0.0
	
	// parse bodies of the node
	for b := n.first ; b != nil; b = b.next {
		// fmt.Printf("updateCOM body adress %x\n", &b)
		if( b.next == b) { panic("linked to itself")}
		
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