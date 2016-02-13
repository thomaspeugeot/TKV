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
	"github.com/thomaspeugeot/tkv/quadtree"
	"image"
	"image/color"
	"image/gif"
	"io"
	"fmt"
	"math"
)


	
// constant to be added to the distance between bodies
// in order to compute repulsion (avoid near infinite repulsion force)
// note : declaring those variable as constant has no impact on benchmarks results
var	ETA float64 = 0.0000001

// pseudo gravitational constant to compute 
var	G float64 = 0.000001
var Dt float64  = 1.0 // 1 second, time step

// velocity cannot be too high in order to stop bodies from overtaking
// each others
var MaxVelocity float64  = 0.001 // cannot make more that 1/1000 th of the unit square per second

// the barnes hut criteria 
var BN_THETA float64 = 0.5 // can use barnes if distance to COM is 5 times side of the node's box

// used to compute speed up
var nbComputationPerStep int

// if true, Barnes-Hut algo is used
var UseBarnesHut bool = true

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


var palette = []color.Color{color.White, color.Black}
const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

// a simulation run
type Run struct {
	bodies * []quadtree.Body // bodies position in the quatree
	bodiesAccel * []Acc // bodies acceleration
	bodiesVel * []Vel // bodies velocity

	q quadtree.Quadtree // the supporting quadtree
}

func (r * Run) getAcc(index int) (* Acc) {
	return & (*r.bodiesAccel)[index]
}

func (r * Run) getVel(index int) (* Vel) {
	return & (*r.bodiesVel)[index]
}

// init the run with an array of quadtree bodies
func (r * Run) Init( bodies * ([]quadtree.Body)) {
	r.bodies = bodies
	acc := make([]Acc, len(*bodies))
	vel := make([]Vel, len(*bodies))
	r.bodiesAccel = &acc
	r.bodiesVel = &vel
	r.q.Init(bodies)
}

func (r * Run) oneStep() {

	nbComputationPerStep =0
	
	// compute the quadtree from the bodies
	r.q.UpdateNodesListsAndCOM()
	
	// compute repulsive forces & acceleration
	r.ComputeRepulsiveForceConcurrent( 20)
	
	// compute velocity
	r.UpdateVelocity()
		
	// compute new position
	r.UpdatePosition()
}


// compute repulsive forces by spreading the calculus
// among nbRoutine go routines
func (r * Run) ComputeRepulsiveForceConcurrent(nbRoutine int) {

	sliceLen := len(*r.bodies)
	done := make( chan struct{})

	// breakdown slice
	for i:=0; i<nbRoutine; i++ {
	
		startIndex := (i*sliceLen)/nbRoutine
		endIndex := ((i+1)*sliceLen)/nbRoutine -1
		go func() { 
			r.ComputeRepulsiveForceSubSet( startIndex, endIndex)
			done <- struct{}{} 
		}()
	}

	// wait for return
	for i:=0; i<nbRoutine; i++ {
		<- done
	}

}

// compute repulsive forces
func (r * Run) ComputeRepulsiveForce() {
	
	r.ComputeRepulsiveForceSubSet( 0, len(*r.bodies))
}

// compute repulsive forces for a sub part of the bodies
func (r * Run) ComputeRepulsiveForceSubSet( startIndex, endIndex int) {

	// parse all bodies
	bodiesSubSet := (*r.bodies)[startIndex:endIndex]
	for idx, _ := range  bodiesSubSet {
		
		// index in the original slice
		origIndex := idx+startIndex
		
		if( UseBarnesHut ) {
			r.computeAccelerationOnBodyBarnesHut( origIndex)
		} else {
			r.computeAccelerationOnBody( origIndex)
		}
	}
}

// parse all other bodies to compute acceleration
func (r * Run) computeAccelerationOnBody(origIndex int) {

	body := (*r.bodies)[origIndex]

	// reset acceleration
	acc := &((*r.bodiesAccel)[origIndex])
	acc.X = 0
	acc.Y = 0
	
	// parse all other bodies for repulsions
	// accumulate repulsion on acceleration
	for idx2, _ := range (*r.bodies) {
		
		if( idx2 != origIndex) {
			body2 := (*r.bodies)[idx2]
			
			x, y := getRepulsionVector( &body, &body2)
			
			acc.X += x
			acc.Y += y
		}
	}
	
}

// parse all other bodies to compute acceleration
// with the barnes-hut algorithm
func (r * Run) computeAccelerationOnBodyBarnesHut(idx int) {

	// reset acceleration
	acc := &((*r.bodiesAccel)[idx])
	acc.X = 0
	acc.Y = 0
	
	// Coord is initialized at the Root coord
	var rootCoord quadtree.Coord
	
	r.computeAccelationWithNodeRecursive( idx, rootCoord)
}

// given a body and a node in the quadtree, compute the repulsive force
func (r * Run) computeAccelationWithNodeRecursive( idx int, coord quadtree.Coord) {
	
	body := (*r.bodies)[idx]
	acc := &((*r.bodiesAccel)[idx])
	
	// compute the node box size
	level := coord.Level()
	boxSize := 1.0 / math.Pow( 2.0, float64(level)) // if level = 0, this is 1.0
	
	node := & (r.q.Nodes[coord])
	dist := getModuloDistanceBetweenBodies( &body, &(node.Body))

	// fmt.Printf("computeAccelationWithNodeRecursive index %d at coord %#v level %d boxSize %f mass %f\n", idx, coord, level, boxSize, node.M)

	// avoid node with zero mass
	if( node.M == 0) {
		return
	}
	
	// check if the COM of the node can be used
	if (boxSize / dist) < BN_THETA {
	
		// fmt.Printf("computeAccelationWithNodeRecursive at node %#v\n", node)
		x, y := getRepulsionVector( &body, &(node.Body))
			
		acc.X += x
		acc.Y += y
	} else {		
		if( level < 8) {
			// parse sub nodes
			// fmt.Printf("computeAccelationWithNodeRecursive go down at node %#v\n", node)
			coordNW, coordNE, coordSW, coordSE := r.q.NodesBelow( coord)
			r.computeAccelationWithNodeRecursive( idx, coordNW)
			r.computeAccelationWithNodeRecursive( idx, coordNE)
			r.computeAccelationWithNodeRecursive( idx, coordSW)
			r.computeAccelationWithNodeRecursive( idx, coordSE)		
		} else {
		
			// parse bodies of the node
			rank := 0
			for b := node.First() ; b != nil; b = b.Next() {
				if( *b != body) {
					// fmt.Printf("computeAccelationWithNodeRecursive at leaf %#v rank %d\n", b, rank)
					x, y := getRepulsionVector( &body, b)
			
					acc.X += x
					acc.Y += y
					rank++
				}
			}
		}
	}
}

func (r * Run) UpdateVelocity() {

	// parse all bodies
	for idx, _ := range (*r.bodies) {

		// update velocity (to be completed with Dt)
		acc := r.getAcc(idx)
		vel := r.getVel(idx)
		vel.X += acc.X * G * Dt
		vel.Y += acc.Y * G * Dt
		
		// put some drag
		vel.X *= 0.75
		vel.Y *= 0.75
		
		// if velocity is above
		velocity := math.Sqrt( vel.X*vel.X + vel.Y*vel.Y)
		
		if velocity > MaxVelocity { 
			vel.X *= MaxVelocity/velocity
			vel.Y *= MaxVelocity/velocity
		}
	}
}

func (r * Run) UpdatePosition() {

	// parse all bodies
	for idx, _ := range (*r.bodies) {
		
		body := &((*r.bodies)[idx])
		
		// updatePos
		vel := r.getVel(idx)
		body.X += vel.X * Dt
		body.Y += vel.Y * Dt
		
		if body.X >= 1.0 { 
			body.X = 1.0 - (body.X - 1.0) 
			vel.X = -vel.X
		}
		if body.X <= 0.0 { 
			body.X = - body.X 
			vel.X = -vel.X
		}
		if body.Y >= 1.0 { 
			body.Y = 1.0 - (body.Y - 1.0) 
			vel.Y = -vel.Y
		}
		if body.Y <= 0.0 { 
			body.Y = - body.Y 
			vel.Y = -vel.Y
		}
	}
}

// output position of bodies of the Run into a GIF representation
func (r * Run) outputGif(out io.Writer, nbStep int) {
	const (
		size    = 600   // image canvas 
		delay   = 4    // delay between frames in 10ms units
	)
	var nframes = nbStep    // number of animation frames
	
	anim := gif.GIF{LoopCount: nframes}
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, size+1, size+1)
		img := image.NewPaletted(rect, palette)

		for idx, _ := range (*r.bodies) {
		
			body := (*r.bodies)[idx]
		
			if false { fmt.Printf("Encoding body %d %f %f\n", idx, body.X, body.Y) }
		
			img.SetColorIndex(
				int(body.X*size+0.5), 
				int(body.Y*size+0.5),
				blackIndex)
		}
		
		// encode time step into the image
		progress := float64(i) / float64 (nframes)
		for j:= 0; j < int( size*progress); j++ {
			img.SetColorIndex(
				j+1, 
				10,
				blackIndex)
		}
		fmt.Printf("\rProgress %f speedup %f", progress, 
			float64(len(*r.bodies)*len(*r.bodies))/float64(nbComputationPerStep))
		
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
		
		r.oneStep()
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

// compute modulo distance
func getModuloDistanceBetweenBodies( A, B *quadtree.Body) float64 {

	x := getModuloDistance( B.X, A.X)
	y := getModuloDistance( B.Y, A.Y)

	distQuared := (x*x + y*y)
	
	return math.Sqrt( distQuared )
}

// compute repulsion force vector between body A and body B
// applied to body A
// proportional to the inverse of the distance squared
func getRepulsionVector( A, B *quadtree.Body) (x, y float64) {

	x = getModuloDistance( B.X, A.X)
	y = getModuloDistance( B.Y, A.Y)

	distQuared := (x*x + y*y) + ETA
	
	distPow3 := distQuared * math.Sqrt( distQuared )
	
	if false { 
		distPow3 := math.Pow( distQuared, 1.5) 
		distQuared /= distPow3
	}
	
	nbComputationPerStep++
	
	return x/distPow3, y/distPow3

	// return x / distQuared, y / distQuared
}

// get modulo distance between alpha and beta.
//
// alpha and beta are between 0.0 and 1.0
// the modulo distance cannot be above 0.5
func getModuloDistance( alpha, beta float64) (dist float64) {

	dist = beta-alpha
	if( dist > 0.5 ) { dist -= 1.0 }
	if( dist < -0.5 ) { dist += 1.0 }
		
	return dist
}