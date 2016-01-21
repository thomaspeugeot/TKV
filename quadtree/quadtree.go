// compact implementation of a static 2D quadtree
// 
// goal is to support the Barnes Hut implementation for spreading bodies on a 2D rectangle.
package quadtree

// the coordinate of a node
// this is coded on a uint32 for speed sake
// the first byte encode the depth of the node (256 levels seems OK)
//           0 is the root node
//			 1 is the level 1 node
//           ...
// second to fourth byte encode the coordinate of the node by block of 2 bits by level
// 00 is for north west
// 01 is for north east
// 10 is for south west
// 11 is for south east
// for instance 
// 0x00 0x00 0x00 0x00 is the root node
// 0x01 0x00 0x00 0x00 is the level 1 node on north west
//
type Coord uint32

// get node level from node coord c
func Level(c Coord) int {
	result := int(c >> 24)
	
	return result
}

// initiate a node coordinate from level
// and set of coordinates
func set( c *Coord, level int, coords []int) {
}

