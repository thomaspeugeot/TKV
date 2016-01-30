// compact implementation of a modified barnes-hut algorithm
//
// goal is to spread evenly bodies on a 2D rectangle
// 
// TKV implementation starts from a Barnes-Hut implementation of the gravitation simulation and make the following modification:
//
// - keep bodies within the canvas: bodies "bumps" on bodders (see updatePos)
// - for spreading, use repulsion instead of gravitational attraction and add friction (see updateVel)
// - use a ring topology instead of a linear topology (think of spreading bodies on a ring, see getDist), modification of metric
package barnes_hut

import (
	"tkv/quadtree"
)

//	Bodies's X,Y position coordinates are float64 between 0 & 1
type Pos struct {
	X float64
	Y float64
}

// Velocity
type Vel struct {
	X float64
	Y float64
}

// Acceleration
type Acc struct {
	X float64
	Y float64
}

// definition of a body
type Body struct {
	Pos
	Vel
	Acc
}

// a simulation run
type Run struct {
	bodies []Body // nb of bodies
	q quadtree.Quadtree // the supporting quadtree
}

func (r * Run) Init( bodies * ([]Body)) {
	r.bodies = *bodies
	r.q.SetupNodesLinks()
}

func (r * Run) oneStep( bodies * ([]Body)) {

	r.updateNodesListsAndCOM()
	
}