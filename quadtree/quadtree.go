/*
Package quadtree is a compact implementation of a static 2D quadtree.

This implementation supports a Barnes Hut algorithm with more that one million bodies.
It is a compact (and fun!) implementation because node coordinates are encoded into a 4 bytes stucture
	type Coord uint32

This quatree is a hierarchical set of nodes that divide the 2D space at different
levels (1 node at level 0, 4 nodes at level 1, 16 nodes at level 2, etc...).
Each node holds sub nodes except the leaf nodes that contains all bodies that are located in the area of the node.

Bodies's X,Y coordinates are float64 between 0 & 1

A quadtree is usualy a dynmic structure. In this package implementation, the architecture is static with a depth of nodes
limited to 8 (256 * 256 cells at the level 8)
*/
package quadtree

import (
	"fmt"
	"math"
	"sort"
	"testing"
)

type QuadtreeGini [9][10]float64

// a Quadtree store Nodes. It is a an array with direct access to the Nodes with the Nodes coordinate
// see Coord
type Quadtree struct {
	Nodes         [1 << 20]Node
	bodies        *[]Body      // pointer to the body slice
	BodyCountGini QuadtreeGini // for each of the 9 levels, tencentile of bodies
}

var optim bool

func init() {
	optim = true
}

// constants used to navigate from one node to the other
const (
	NW = 0x0000
	NE = 0x0100
	SW = 0x0001
	SE = 0x0101
)

// init quadtree
func (q *Quadtree) Init(bodies *[]Body) {
	q.bodies = bodies
	q.setupNodesCoord()
	q.setupNodesLinks()
	q.updateNodesList()
	q.updateNodesCOM()
}

// updates nodes according to bodies locations
// this function should be called when bodies have been moved
func (q *Quadtree) UpdateNodesListsAndCOM() {

	q.updateNodesList()
	q.updateNodesCOM()

	var t testing.T
	q.CheckIntegrity(&t)
}

// compute quadtree Nodes for levels from 0 to 7
func (q *Quadtree) updateNodesCOMAbove8() {

	for level := 7; level >= 0; level-- {

		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)

		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {
				coord := GetCoord(level, i, j)
				node := &(q.Nodes[coord])
				node.updateCOM()
			}
		}
	}
}

// get nodes coords below
func NodesBelow(c Coord) (coordNW, coordNE, coordSW, coordSE Coord) {

	levelBelow := c.Level() + 1
	i := c.X()
	j := c.Y()
	shift := uint(8 - levelBelow)

	// to go east at the level below, we flip to 1 the bit that is significant at that level
	coordNW = Coord(uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | NW<<shift)
	coordNE = Coord(uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | NE<<shift)
	coordSW = Coord(uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | SW<<shift)
	coordSE = Coord(uint(levelBelow)<<16 | uint(i)<<8 | uint(j) | SE<<shift)

	return coordNW, coordNE, coordSW, coordSE
}

// setup node coord
func (q *Quadtree) setupNodesCoord() {
	for level := 8; level >= 0; level-- {

		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)

		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {
				coord := GetCoord(level, i, j)
				node := &(q.Nodes[coord])
				node.coord = coord
			}
		}
	}
}

// setup quadtree Nodes for levels from 7 to 0
func (q *Quadtree) setupNodesLinks() {

	Trace.Println("setupNodesLinks")
	for level := 7; level >= 0; level-- {

		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)

		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {

				coord := GetCoord(level, i, j)
				node := &(q.Nodes[coord])

				// s := fmt.Sprintf("SetupNodesLinks level %8d i %8d j %8d coord %s",
				// level, i, j, node.coord.String())
				// fmt.Println(s)

				coordNW, coordNE, coordSW, coordSE := NodesBelow(coord)

				nodeNW := &q.Nodes[coordNW]
				nodeNE := &q.Nodes[coordNE]
				nodeSW := &q.Nodes[coordSW]
				nodeSE := &q.Nodes[coordSE]

				// bodies of the nodes below are chained
				// fmt.Printf("%8x\n", coord)
				node.first = &(nodeNW.Body)
				nodeNW.Body.next = &(nodeNE.Body)
				nodeNE.Body.next = &(nodeSW.Body)
				nodeSW.Body.next = &(nodeSE.Body)
			}
		}
	}
}

// fill quadtree at level 8 with bodies
func (q *Quadtree) updateNodesList() {

	Trace.Println("updateNodesList")

	for idx := range *q.bodies {

		b := &((*q.bodies)[idx])

		// 1st phase, remove the body from its current double linked list
		// link the next body to the previous one
		if b.next != nil {
			b.next.prev = b.prev
		}
		// link the previous body to the next one
		if b.prev != nil {
			b.prev.next = b.next
		} else {
			// if body prev is nil,
			// it can be either the current first of a node or it has not been initialized
			// if it is the current first of the node,
			// the first of the node shall point to the next of the body
			if (q.Nodes[b.coord]).first == b {
				(q.Nodes[b.coord]).first = b.next
			}
		}

		// 2nd Phase
		// put body as the first body of the node
		// shift the first body if it is already there
		// compute coord of body (this is direct access)
		coord := b.getCoord8()
		node := &(q.Nodes[coord])
		initialFirstBody := node.first
		if (initialFirstBody != nil) && (initialFirstBody != b) {
			// double link body to the current node's first
			b.next = initialFirstBody
			initialFirstBody.prev = b
		}

		// body b is the new node's first
		node.first = b
		b.prev = nil

		// setup new coord
		b.coord = coord

		if b.next == b {
			s := fmt.Sprintf("updateNodesList: Node linked to itself coord : idx %d, %s", idx, b.coord.String())
			panic(s)
		}
	}
}

// compute COM of quadtree from level 8 to level 0
func (q *Quadtree) updateNodesCOM() {

	Trace.Println("updateNodesCOM")
	// compute is bottom up
	for level := 8; level >= 0; level-- {

		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)

		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {

				coord := GetCoord(level, i, j)
				node := &(q.Nodes[coord])

				// fmt.Println("updateNodesCOM ", q.Nodes[coord].coord.String())
				// s := fmt.Sprintf("updateNodesCOM level %8d i %8d j %8d coord %s",
				// level, i, j, node.coord.String())
				// fmt.Println(s)
				node.updateCOM()
			}
		}
	}
}

// check integrity of the quadtree by performing
// all kinds of test
func (q *Quadtree) CheckIntegrity(t *testing.T) {

	Trace.Printf("CheckIntegrity")
	nbBodies := 0

	// perform some tests on the links of each nodes
	for level := 8; level >= 8; level-- {

		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)

		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {

				coord := GetCoord(level, i, j)
				node := &(q.Nodes[coord])

				// test that the node coord is correct
				if q.Nodes[coord].coord != coord {
					s := fmt.Sprintf("node coord = %s, want %s",
						q.Nodes[coord].coord.String(), coord.String())
					t.Errorf(s)
				}

				// test that the node first body
				// has a nil previous body
				if node.first != nil && node.first.prev != nil {
					s := fmt.Sprintf("node coord = %s, has first body with non nil prev",
						q.Nodes[coord].coord.String())
					t.Errorf(s)
				}

				// test for each body of the chain of bodies
				// - that the next body previous body is the body
				rank := 0
				for b := node.first; b != nil; b = b.next {

					if b.next != nil && b.next.prev != b {
						s := fmt.Sprintf("node coord = %s, has %d nth body with next body not point to him for prev",
							q.Nodes[coord].coord.String(), rank)
						t.Errorf(s)
					}
					nbBodies++
					rank++
				}
			}
		}
	}

	// check that all bodies are accounted for
	if nbBodies != len(*q.bodies) {
		t.Errorf("Nb bodies do not match expected %d, got %d", len(*q.bodies), nbBodies)
	}
}

// compute number of bodies per node
// update the counting of bodies per node for all levels
func (q *Quadtree) ComputeNbBodiesPerNode() {

	// perform some tests on the links of each nodes
	for level := 8; level >= 0; level-- {

		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)

		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {

				coord := GetCoord(level, i, j)
				node := &(q.Nodes[coord])
				q.updateBodiesNb(node)
			}
		}
	}
}

// compute the gini of body density par node at level 8
func (q *Quadtree) ComputeQuadtreeGini() {

	Trace.Printf("ComputeQuadtreeGini begin")
	q.ComputeNbBodiesPerNode()

	// perform some tests on the links of each nodes
	for level := 8; level >= 0; level-- {

		rank := 0

		// nb of nodes for the current level
		nbNodesX := 1 << uint(level)
		nbNodesY := 1 << uint(level)

		// var bodyCount []int
		bodyCountPerLevel := make([]int, nbNodesX*nbNodesY)

		// parse nodes of level
		for i := 0; i < nbNodesX; i++ {
			for j := 0; j < nbNodesY; j++ {

				coord := GetCoord(level, i, j)
				node := &(q.Nodes[coord])
				bodyCountPerLevel[rank] += node.nbBodies
				rank++
				// fmt.Println( fmt.Sprintf("i %d j %d: %d", i, j, nbBodies))
			}
		}
		sort.Ints(bodyCountPerLevel)

		for tencile := 0; tencile < 10; tencile++ {
			lowIndex := int(math.Floor(float64(nbNodesX*nbNodesY) * float64(tencile) / 10.0))
			highIndex := int(math.Floor(float64(nbNodesX*nbNodesY) * float64(tencile+1) / 10.0))

			nbBodiesInTencile := 0
			for _, nbBodies := range bodyCountPerLevel[lowIndex:highIndex] {
				nbBodiesInTencile += nbBodies
			}
			q.BodyCountGini[level][tencile] = float64(nbBodiesInTencile) // /float64(len(*q.bodies))
		}
	}
	Trace.Printf("ComputeQuadtreeGini end")
}

// consolidate the number of bodies attached to the node
// at level 8, this is the number of bodies
// above level 8, this is an aggregate of the number of bodies at the level below
func (q *Quadtree) updateBodiesNb(n *Node) {
	n.nbBodies = 0

	if n.coord.Level() == 8 {
		for b := n.first; b != nil; b = b.next {
			n.nbBodies++
		}
	} else {
		coordNW, coordNE, coordSW, coordSE := NodesBelow(n.coord)

		nodeNW := &q.Nodes[coordNW]
		nodeNE := &q.Nodes[coordNE]
		nodeSW := &q.Nodes[coordSW]
		nodeSE := &q.Nodes[coordSE]

		n.nbBodies += nodeNW.nbBodies
		n.nbBodies += nodeNE.nbBodies
		n.nbBodies += nodeSW.nbBodies
		n.nbBodies += nodeSE.nbBodies
	}
}
