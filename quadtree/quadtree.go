// Compact implementation of a static 2D quadtree
// 
// 1st Goal is to support the Barnes Hut implementation for spreading bodies on a 2D square.
// 
// 2nd Goal is to support more than 1 million bodies
// there are constraints:
//
// - X,Y coordinates are float64 between 0 & 1
//
// - quadtree architecture is static
//
// - depth is limited to 8 (256 * 256 cells at the level 8)
//
// to see the doc
//
// 		godoc -goroot=$GOPATH -http=:8080
package quadtree

// 
// Coordinate system of a node
// 
// Coordinates of a node are coded as follow
//
// 	1st byte : level (root = 0, max depth = 7) 
// 	2nd byte : X coordinate 
// 	3nd byte : Y coordinate : coded on 
//	4th byte : 0x00 (unused) 
//
// the X coordinate code depends on the level
//	level 0: there is no coordinate
//	level 1: west if 0, east is 128 (0x80)
//	level 2: quarters coordinates are 0, 64 (0x40), 128 (0x80), 192 (0x84)
//	...
type Coord uint32

// get node level from node coord c
func Level(c Coord) int {
	result := int(c >> 24)
	
	return result
}

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

