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

var palette = []color.Color{color.White, color.Black}
const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

// a simulation run
type Run struct {
	bodies * []quadtree.Body // bodies position in the quatree
	bodiesAccel * []Acc // bodies acceleration

	q quadtree.Quadtree // the supporting quadtree
}

func (r * Run) getAcc(index int) Acc {

	return (*r.bodiesAccel)[index]
}

// init the run with an array of quadtree bodies
func (r * Run) Init( bodies * ([]quadtree.Body)) {
	r.bodies = bodies
	acc := make([]Acc, len(*bodies))
	r.bodiesAccel = &acc
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
			body2 := (*r.bodies)[idx2]
			
			x, y := getRepulsionVector( &body, &body2)
			acc.X += x
			acc.Y += y
		}
	}
}

func (r * Run) UpdatePosition() {

	// parse all bodies
	for idx, _ := range (*r.bodies) {
		
		body := &((*r.bodies)[idx])
		
		// updatePos
		acc := (*r.bodiesAccel)[idx]
		body.X += acc.X / 10000
		body.Y += acc.Y / 10000
	}
}

// output position of bodies of the Run into a GIF representation
func (r * Run) outputGif(out io.Writer, nbStep int) {
	const (
		size    = 500   // image canvas covers [-size..+size]
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
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
		
		r.ComputeRepulsiveForce()
		r.UpdatePosition()
		
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

// compute repulsion force vector between body A and body B
// applied to body A
func getRepulsionVector( A, B *quadtree.Body) (x, y float64) {

	x = getModuloDistance( B.X, A.X)
	y = getModuloDistance( B.Y, A.Y)

	return x, y
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