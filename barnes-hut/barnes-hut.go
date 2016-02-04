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
	r.q.SetupNodesLinks()
}

func (r * Run) oneStep() {

	// compute the quadtree from the bodies
	r.q.UpdateNodesListsAndCOM( r.bodies)
	
	// compute repulsive forces & acceleration
	
	// compute velocity
	
	// compute new position
	
}

func (r * Run) ComputeRepulsiveForce() {

	// parse all bodies
	for idx, _ := range (*r.bodies) {
		
		body := (*r.bodies)[idx]
		
		// reset acceleration
		acc := &((*r.bodiesAccel)[idx])
		acc.X = 0
		acc.Y = 0
		
		// parse all other bodies for repulsions
		// accumulate repulsion on acceleration
		for idx2, _ := range (*r.bodies) {
			
			if( idx2 != idx) {
				body2 := (*r.bodies)[idx2]
				
				x, y := getRepulsionVector( &body, &body2)
				
				acc.X += x
				acc.Y += y
			}
			
		}
	}
}

func (r * Run) UpdateVelocity() {

	// parse all bodies
	for idx, _ := range (*r.bodies) {

		// update velocity (to be completed with dt)
		acc := r.getAcc(idx)
		vel := r.getVel(idx)
		vel.X += acc.X / 10000000
		vel.Y += acc.Y / 10000000
		
		// put some drag
		vel.X *= 0.9
		vel.Y *= 0.9
	}
}

func (r * Run) UpdatePosition() {

	// parse all bodies
	for idx, _ := range (*r.bodies) {
		
		body := &((*r.bodies)[idx])
		
		// updatePos
		vel := r.getVel(idx)
		body.X += vel.X
		body.Y += vel.Y
		
		if body.X > 1.0 { 
			body.X = 1.0 - (body.X - 1.0) 
			vel.X = -vel.X
		}
		if body.X < 0.0 { 
			body.X = - body.X 
			vel.X = -vel.X
		}
		if body.Y > 1.0 { 
			body.Y = 1.0 - (body.Y - 1.0) 
			vel.Y = -vel.Y
		}
		if body.Y < 0.0 { 
			body.Y = - body.Y 
			vel.Y = -vel.Y
		}
	}
}

// output position of bodies of the Run into a GIF representation
func (r * Run) outputGif(out io.Writer, nbStep int) {
	const (
		size    = 200   // image canvas covers [-size..+size]
		delay   = 50     // delay between frames in 10ms units
	)
	var nframes = nbStep    // number of animation frames
	
	anim := gif.GIF{LoopCount: nframes}
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)

		for idx, _ := range (*r.bodies) {
		
			body := (*r.bodies)[idx]
		
			if false { fmt.Printf("Encoding body %d %f %f\n", idx, body.X, body.Y) }
		
			img.SetColorIndex(
				size+int(body.X*size+0.5), 
				size+int(body.Y*size+0.5),
				blackIndex)
		}
		
		// encode time step into the image
		for j:= 0; j< i; j++ {
			img.SetColorIndex(
				j+1, 
				10,
				blackIndex)
		}
		
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
		
		r.ComputeRepulsiveForce()
		r.UpdateVelocity()
		r.UpdatePosition()
		
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

// compute repulsion force vector between body A and body B
// applied to body A
// proportional to the inverse of the distance squared
func getRepulsionVector( A, B *quadtree.Body) (x, y float64) {

	x = getModuloDistance( B.X, A.X)
	y = getModuloDistance( B.Y, A.Y)

	distQuared := (x*x + y*y)
	distPow3 := math.Pow( distQuared, 1.5)
	
	return x/distPow3, y/distPow3
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