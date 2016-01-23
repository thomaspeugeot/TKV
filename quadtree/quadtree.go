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
//	level 2: quarters coordinates are 0, 64 (0x40), 128 (0x80), 192 (0x84)
//	...
type Coord uint32

// node level of a node coord c
func Level(c Coord) int { return int( c >> 16) }

// x coord node coord c
func x(c Coord) int { return int((c & 0x0000FFFF) >> 8) }

// check encoding of c
func checkIntegrity( c Coord) bool {
	
	// check byte 4 is null
	if res := 0x000000FF & c; res != 0x00 {
		return false
	}
	
	return true
}

// suite of integer between 0 & 3 denoting position of the 
// node within the node above level.
type PosSuite []int

// initiate a node coordinate from level and coords
// coords is an array of int with length level
func set( c *Coord, level int, coords PosSuite) bool {
	result := true
	
	if len(coords) != level { 
		result = false
	}
	
	
	return result
}

