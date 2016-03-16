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
	"os"
	"log"
	"fmt"
	"math"
	"math/rand"
	"time"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
)

// constant to be added to the distance between bodies
// in order to compute repulsion (avoid near infinite repulsion force)
// note : declaring those variable as constant has no impact on benchmarks results
var	ETA float64 = 0.00000001

// pseudo gravitational constant to compute 
var	G float64 = 0.01
var Dt float64  = 0.1 // 0.1 second, time step
var DtRequest = Dt // new value of Dt requested by the UI. The real Dt will be changed at the end of the current step.

// velocity cannot be too high in order to stop bodies from overtaking
// each others
var MaxVelocity float64  = 0.001 // cannot make more that 1/1000 th of the unit square per second

// the barnes hut criteria 
var BN_THETA float64 = 0.5 // can use barnes if distance to COM is 5 times side of the node's box
var ThetaRequest = BN_THETA // new value of theta requested by the UI. The real BN_THETA will be changed at the end of the current step.

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

type State string

const (
	STOPPED = "STOPPED"
	RUNNING = "RUNNING"
)

//
var ConcurrentRoutines int = 100

// a simulation run
type Run struct {
	bodies * []quadtree.Body // bodies position in the quatree
	bodiesAccel * []Acc // bodies acceleration
	bodiesVel * []Vel // bodies velocity

	q quadtree.Quadtree // the supporting quadtree
	state State
	step int
	giniOverTime [][]float64 // evolution of the gini distribution over time 
	xMin, xMax, yMin, yMax float64 // coordinates of the rendering windows
}

func (r * Run) getAcc(index int) (* Acc) {
	return & (*r.bodiesAccel)[index]
}

func (r * Run) getVel(index int) (* Vel) {
	return & (*r.bodiesVel)[index]
}

func (r * Run) State() State{
	return r.state
}

func (r * Run) GetStep() int{
	return r.step
}

func (r * Run) SetState(s State) {
	r.state = s
}

func (r * Run) SetRenderingWindow( xMin, xMax, yMin, yMax float64) {
	r.xMin, r.xMax, r.yMin, r.yMax = xMin, xMax, yMin, yMax
}
func (r * Run) GiniOverTimeTransposed() [][]float64 {

	var giniOverTimeTransposed [][]float64 
	// := r.giniOverTime
	giniOverTimeTransposed = transposeFloat64( r.giniOverTime)
	return giniOverTimeTransposed
}

func (r * Run) GiniOverTime() [][]float64 {

	return r.giniOverTime
}

// init the run with an array of quadtree bodies
func (r * Run) Init( bodies * ([]quadtree.Body)) {
	r.bodies = bodies
	acc := make([]Acc, len(*bodies))
	vel := make([]Vel, len(*bodies))
	r.bodiesAccel = &acc
	r.bodiesVel = &vel
	r.q.Init(bodies)
	r.state = STOPPED
	r.SetRenderingWindow( 0.0, 0.0, 1.0, 1.0)
}

func (r * Run) OneStep() {

	nbComputationPerStep =0

	// update Dt according to request
	Dt = DtRequest
	BN_THETA = ThetaRequest
	
	// compute the quadtree from the bodies
	r.q.UpdateNodesListsAndCOM()
	
	// compute repulsive forces & acceleration
	r.ComputeRepulsiveForceConcurrent( ConcurrentRoutines)
	
	// compute velocity
	r.UpdateVelocity()
		
	// compute new position
	r.UpdatePosition()

	// update the step
	r.step++

	fmt.Printf("step %d speedup %f low 10 %f high 5 %f high 10 %f\n",
		r.step, 
		float64(len(*r.bodies)*len(*r.bodies))/float64(nbComputationPerStep),
		r.q.BodyCountGini[8][0],
		r.q.BodyCountGini[8][5],
		r.q.BodyCountGini[8][9])
		
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

			// fmt.Printf("computeAccelerationOnBody idx2 %3d x %9.3f y %9.3f \n", idx2, x, y)
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

	
	// avoid node with zero mass
	if( node.M == 0) {
		return
	}
	
	// fmt.Printf("computeAccelationWithNodeRecursive index %d at coord %#v level %d boxSize %f mass %f\n", idx, coord, level, boxSize, node.M)

	// check if the COM of the node can be used
	if (boxSize / dist) < BN_THETA {
	
		x, y := getRepulsionVector( &body, &(node.Body))
			
		acc.X += x
		acc.Y += y

		// fmt.Printf("computeAccelationWithNodeRecursive at node %#v x %9.3f y %9.3f\n", node.Coord(), x, y)

	} else {		
		if( level < 8) {
			// parse sub nodes
			// fmt.Printf("computeAccelationWithNodeRecursive go down at node %#v\n", node.Coord())
			coordNW, coordNE, coordSW, coordSE := quadtree.NodesBelow( coord)
			r.computeAccelationWithNodeRecursive( idx, coordNW)
			r.computeAccelationWithNodeRecursive( idx, coordNE)
			r.computeAccelationWithNodeRecursive( idx, coordSW)
			r.computeAccelationWithNodeRecursive( idx, coordSE)		
		} else {
		
			// parse bodies of the node
			rank := 0
			for b := node.First() ; b != nil; b = b.Next() {
				if( *b != body) {
					x, y := getRepulsionVector( &body, b)
			
					acc.X += x
					acc.Y += y
					rank++
					// fmt.Printf("computeAccelationWithNodeRecursive at leaf %#v rank %d x %9.3f y %9.3f\n", b.Coord(), rank, x, y)
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

func (r * Run) RenderGif(out io.Writer) {
	const (
		size    = 600   // image canvas 
		delay   = 4    // delay between frames in 10ms units
		nframes = 0
	)
	anim := gif.GIF{LoopCount: nframes}
	rect := image.Rect(0, 0, size+1, size+1)
	img := image.NewPaletted(rect, palette)
		
	for idx, _ := range (*r.bodies) {
	
		body := (*r.bodies)[idx]
	
		if false { fmt.Printf("Encoding body %d %f %f\n", idx, body.X, body.Y) }
	
		// take into account rendering window
		if( (body.X > r.xMin) && (body.X < r.xMax) && (body.Y > r.yMin) && (body.Y < r.yMax) ) {

			imX := (body.X - r.xMin)/(r.xMax-r.xMin)
			imY := (body.Y - r.yMin)/(r.yMax-r.yMin)

			img.SetColorIndex(
			int(imX*size+0.5), 
			int(imY*size+0.5),
			blackIndex)
		}
	}
	anim.Delay = append(anim.Delay, delay)
	anim.Image = append(anim.Image, img)
	var b bytes.Buffer
	gif.EncodeAll(&b, &anim)
	encodedB64 := base64.StdEncoding.EncodeToString([]byte(b.Bytes()))
	out.Write( []byte(encodedB64))

}

// output position of bodies of the Run into a GIF representation
func (r * Run) OutputGif(out io.Writer, nbStep int) {

	for r.step < nbStep  {

		// if state is STOPPED, pause
		for r.state == STOPPED {
			time.Sleep(100 * time.Millisecond)
		}
		r.q.ComputeQuadtreeGini()

		// append the new gini elements
		// create the array
		giniArray := make( []float64, 10)
		copy( giniArray, r.q.BodyCountGini[8][:])
		r.giniOverTime = append( r.giniOverTime, giniArray)

		r.OneStep()
	}
}

// serialize bodies's state vector into a file
// convention is "step-xxxx.bod"
// return true if operation was successfull 
// works only if state is STOPPED
func (r * Run) CaptureConfig() bool {
	if r.state == STOPPED {
		filename := fmt.Sprintf("conf-TST-%05d.bods", r.step)
		file, err := os.Create(filename)
		if( err != nil) {
			log.Fatal(err)
			return false
		}
		jsonBodies, _ := json.MarshalIndent( r.bodies, "","\t")
		file.Write( jsonBodies)
		file.Close()
		return true
	} else {
		return false
	}
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
	
	// repulsion is proportional to mass
	massCombined := A.M * B.M
	x *= massCombined
	y *= massCombined

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

// function used to spread bodies randomly on 
// the unit square
func SpreadOnCircle(bodies * []quadtree.Body) {
	for idx, _ := range *bodies {
		
		body := &((*bodies)[idx])
		
		radius := rand.Float64()
		angle := 2.0 * math.Pi * rand.Float64()
		
		if idx%2 == 0 {
			body.X = 0.2
			body.Y = 0.7
			radius *= 0.15
		} else {
			body.X = 0.6
			body.Y = 0.4
			radius *= 0.25
		}
		
		body.M =0.1000000
		body.X += math.Cos( angle) * radius
		body.Y += math.Sin( angle) * radius
	}
}

func (r * Run) BodyCountGini() quadtree.QuadtreeGini {
	return r.q.BodyCountGini
}

var CurrentCountry = "TST"

// return the list of available configuration
func (r * Run) DirConfig() []string {

	// open the current working directory
	cwd, error := os.Open(".")

	if( error != nil ) {
		panic( "not able to open current working directory")
	}

	// get files with their names
	names, err := cwd.Readdirnames(0)

	if( err != nil ) {
		panic( "cannot read names in current working directory")
	}

	// parse the list of names and pick the ones that match the 
	var result []string

	for _, dirname := range(names) {
		if strings.Contains( dirname, CurrentCountry) {
			result = append( result, dirname)
		}
	}

	return result
}