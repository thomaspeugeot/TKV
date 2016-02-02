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
	bodies * []quadtree.Body // bodies
	q quadtree.Quadtree // the supporting quadtree
}

func (r * Run) Init( bodies * ([]quadtree.Body)) {
	r.bodies = bodies
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
		
		// parse all repulsions
		for idx2, _ := range (*r.bodies) {
			body2 := (*r.bodies)[idx2]
			
		}
	}
}


func (r * Run) outputGif(out io.Writer) {
	const (
		size    = 500   // image canvas covers [-size..+size]
		nframes = 1    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
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
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
