package quadtree

import (
	"fmt"
)

// a node is a body
type Node struct {
	// bodies of the node
	//  at level 8, this is the list of bodies pertaining in the bounding box of the node
	// for level 7 to 0, this is the list of the four bodies of the nodes at the level below (or +1)

	// Barycenter with mass of all the bodies of the node
	// this body is linked with the bodies at his level in the node
	Body
	first    *Body // link to the bodies below
	coord    Coord // the coordinate of the Node
	nbBodies int   // number of bodies in the node
}

// link to the first body of the bodies chain belonging to the node
func (n *Node) First() *Body { return n.first }

func (n *Node) Coord() Coord { return n.coord }

// update COM of a node (reset the current COM before)
func (n *Node) updateCOM() {

	n.M = 0.0
	n.X = 0.0
	n.Y = 0.0

	// fmt.Println("updateCOM ", n.coord.String())

	// parse bodies of the node
	rank := 0
	for b := n.first; b != nil; b = b.next {

		// fmt.Printf("updateCOM body adress %x\n", &b)
		if b.next == b {
			s := fmt.Sprintf("Node linked to itself coord : rank %d, %s", rank, n.coord.String())
			panic(s)
		}

		n.M += b.M
		n.X += b.X * b.M
		n.Y += b.Y * b.M

		rank++
	}

	// divide by total mass to get the barycenter
	if n.M > 0 {
		n.X /= n.M
		n.Y /= n.M
	}
}
