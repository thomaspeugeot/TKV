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
	"image/color"
	"fmt"
	"math"
	"math/rand"
	"time"
	"sort"
	"sync/atomic"
	"sync"
	)

// constant to be added to the distance between bodies
// in order to compute repulsion (avoid near infinite repulsion force)
// note : declaring those variable as constant has no impact on benchmarks results
var	ETA float64 = 1e-10

// pseudo gravitational constant to compute 
var	G float64 = 0.01
//var Dt float64  = 3*1e-8 // difficult to fine tune
var Dt float64  = 2.3*1e-10 // difficult to fine tune
var DtRequest = Dt // new value of Dt requested by the UI. The real Dt will be changed at the end of the current step.

// velocity cannot be too high in order to stop bodies from overtaking
// each others
var MaxDisplacement float64  = 0.001 // cannot make more that 1/1000 th of the unit square per second

// the barnes hut criteria 
var BN_THETA float64 = 0.5 // can use barnes if distance to COM is 5 times side of the node's box
var ThetaRequest = BN_THETA // new value of theta requested by the UI. The real BN_THETA will be changed at the end of the current step.

// how much drag we put (1.0 is no drag)
// tis criteria is important because it favors bodies that moves freely against bodies that are stuck on a border
var SpeedDragFactor float64 = 0.3 // 0.99 makes a very bumpy behavior for the Dt

// used to compute speed up
var nbComputationPerStep uint64

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


// var palette = []color.Color{color.White, color.Black}
var palette = []color.Color{color.White, color.Black, color.RGBA{255,0,0,255}, color.RGBA{0,255,0,255}, }
const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
	redIndex = 2 // next color in palette
	blueIndex = 3 // next color in palette
)

// decides wether Dt is set manual or automaticaly
type DtAdjustModeType string
var DtAdjustMode DtAdjustModeType
const (
	AUTO = "AUTO"
	MANUAL = "MANUAL"
)

// state of the simulation
type State string
const (
	STOPPED = "STOPPED"
	RUNNING = "RUNNING"
)

// decide wether, villages borders are drawn
type RenderState string
const (
	WITHOUT_BORDERS = "WITHOUT_BORDERS"
	WITH_BORDERS = "WITH_BORDERS"
)
var ratioOfBorderVillages = 0.1 // ratio of villages that are eligible for marking a border 

// decide wether to display the original configuration or the running configruation
type RenderChoice string
const (
	ORIGINAL_CONFIGURATION = "ORIGINAL_CONFIGURATION"
	RUNNING_CONFIGURATION = "RUNNING_CONFIGURATION"
)

// set the number of concurrent routine for the physic calculation
// this value can be set interactively during the run
var ConcurrentRoutines int = 100

// number of village per X or Y axis. For 10 000 villages, this number is 100
// this value can be set interactively during the run
var nbVillagePerAxe int = 100 

// a simulation run
type Run struct {
	bodies * []quadtree.Body // bodies position in the quatree
	bodiesOrig * []quadtree.Body // original bodies position in the quatree
	bodiesAccel * []Acc // bodies acceleration
	bodiesVel * []Vel // bodies velocity

	q quadtree.Quadtree // the supporting quadtree
	state State
	step int
	giniOverTime [][]float64 // evolution of the gini distribution over time 
	xMin, xMax, yMin, yMax float64 // coordinates of the rendering windows
	renderState RenderState
	renderChoice RenderChoice

	minInterBodyDistance float64 // computed at each step (to compute optimal DT value)
	maxVelocity float64 // max velocity
	dtOptim float64 // optimal dt
	ratioOfBodiesWithCapVel float64 // ratio of bodies where the speed has been capped

	status string // status of the run
}

// rendering the data set can be done only outside the load config xxx function
var rendering sync.Mutex						

func (r * Run) RatioOfBodiesWithCapVel() float64 {
	// log.Output( 1, fmt.Sprintf( "ratioOfBodiesWithCapVel %f ", r.ratioOfBodiesWithCapVel))
	return r.ratioOfBodiesWithCapVel
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

func NbVillagePerAxe() int {
	return nbVillagePerAxe
}

func SetNbVillagePerAxe(nbVillagePerAxe_p int) {
	nbVillagePerAxe = nbVillagePerAxe_p
}

func SetNbRoutines(nbRoutines_p int) {
	ConcurrentRoutines = nbRoutines_p
}

func SetRatioBorderBodies( ratioOfBorderVillages_p float64) {
	ratioOfBorderVillages = ratioOfBorderVillages_p
}

func (r * Run) GetMinInterBodyDistance() float64 {
	return r.minInterBodyDistance
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

	// create a reference of the bodies
	copySliceOfBodies := make( []quadtree.Body, len(*bodies))
	r.bodiesOrig = &copySliceOfBodies
	copy(  *r.bodiesOrig, *r.bodies)

	acc := make([]Acc, len(*bodies))
	vel := make([]Vel, len(*bodies))
	r.bodiesAccel = &acc
	r.bodiesVel = &vel
	r.q.Init(bodies)
	r.state = STOPPED
	r.SetRenderingWindow( 0.0, 0.0, 1.0, 1.0)
	r.renderState = WITH_BORDERS // we draw borders
	r.renderChoice = RUNNING_CONFIGURATION // we draw borders

	DtAdjustMode = MANUAL

	// init measures
	// r.OneStepOptional( false)
}

func (r * Run) ToggleRenderChoice() {
	if r.renderChoice == RUNNING_CONFIGURATION {
		r.renderChoice = ORIGINAL_CONFIGURATION
	} else {
		r.renderChoice = RUNNING_CONFIGURATION
	}
}

func (r * Run) ToggleManualAuto() {
	if DtAdjustMode == MANUAL {
		DtAdjustMode = AUTO
	} else {
		DtAdjustMode = MANUAL
	}
}

// compute the density per village and return the density per village
func (r * Run) ComputeDensityTencilePerVillageString() [10]string {
	var densityString [10]string
	density := r.ComputeDensityTencilePerVillage()
	for tencile,_ := range density {
		densityString[tencile] = fmt.Sprintf( "%3.2f", density[tencile])
	}
	return densityString
}
func (r * Run) ComputeDensityTencilePerVillage() [10]float64 {

	// log.Output( 1, fmt.Sprintf( "ComputeDensityTencilePerVillage %d ", nbVillagePerAxe))

	// parse all bodies
	// prepare the village
	villages := make([][]int, nbVillagePerAxe)
	for x,_  := range villages {
		villages[x] = make([]int, nbVillagePerAxe)
	}

	// parse bodies
	for _,b := range *r.bodies {
		// compute village coordinate (from 0 to nbVillagePerAxe-1)
		x := int( math.Floor(float64( nbVillagePerAxe) * b.X))
		y := int( math.Floor(float64( nbVillagePerAxe) * b.Y))

		villages[x][y]++
	}

	// var bodyCount []int
	nbVillages := nbVillagePerAxe*nbVillagePerAxe
	bodyCountPerVillage := make([]int, nbVillages)
	for x,_  := range villages {
		for y,_  := range villages[x] {
			bodyCountPerVillage[y + x*nbVillagePerAxe] = villages[x][y]
		}
	}
	sort.Ints(bodyCountPerVillage)


	var density [10]float64
	for tencile,_ := range density {
		lowIndex  := int(math.Floor(float64(nbVillages) * float64(tencile)/10.0))
		highIndex := int(math.Floor(float64(nbVillages) * float64(tencile+1)/10.0))
		// log.Output( 1, fmt.Sprintf( "tencile %d ", tencile))
		// log.Output( 1, fmt.Sprintf( "lowIndex %d ", lowIndex))
		// log.Output( 1, fmt.Sprintf( "highIndex %d ", highIndex))
		
		nbBodiesInTencile := 0
		for _, nbBodies := range bodyCountPerVillage[lowIndex:highIndex] {
			nbBodiesInTencile += nbBodies
		}
		density[tencile] = float64(nbBodiesInTencile) / float64(len(bodyCountPerVillage[lowIndex:highIndex]))

		// we compare with then average bodies per villages
		density[tencile] /=	float64(len( *r.bodies)) / float64( nbVillages)	

		// we round the density to 0.01 precision, and put it in percentage point
		density[tencile] *= 100.0 * 100.0
		intDensity := math.Floor( density[tencile] )
		density[tencile] = float64( intDensity) / 100.0
	}

	return density
}

func (r * Run) OneStep() {
	r.OneStepOptional( true)
}
func (r * Run) OneStepOptional( updatePosition bool) {

	t0 := time.Now()

	nbComputationPerStep = 0
	r.minInterBodyDistance = 2.0 // cannot be in a 1.0 by 1.0 square
	r.maxVelocity = 0.0

	// update Dt according to request
	if DtAdjustMode == MANUAL {
		Dt = DtRequest
	} else {
		if( r.dtOptim > 0.0) {	Dt = r.dtOptim }
	} 
	// adjust Dt


	BN_THETA = ThetaRequest
	
	// compute the quadtree from the bodies
	r.q.UpdateNodesListsAndCOM()
	
	// compute repulsive forces & acceleration
	r.ComputeRepulsiveForceConcurrent( ConcurrentRoutines)
	
	// compute velocity
	r.UpdateVelocity()
		
	// compute new position
	if( updatePosition) { r.UpdatePosition() }

	// update the step
	r.step++

	//	fmt.Printf("step %d speedup %f low 10 %f high 5 %f high 10 %f MFlops %f Dur (s) %f MinDist %f Max Vel %f Optim Dt %f Dt %f ratio %f \n",
	r.status = fmt.Sprintf("step %d speedup %f MFlops %f Dur (s) %e MinDist %e Max Vel %e Dt Opt %e Dt %e ratio %e \n",
		r.step, 
		float64(len(*r.bodies)*len(*r.bodies))/float64(nbComputationPerStep),
		// r.q.BodyCountGini[8][0],
		// r.q.BodyCountGini[8][5],
		// r.q.BodyCountGini[8][9],
		Gflops*1000.0,
		StepDuration/1000000000	,
		r.minInterBodyDistance,
		r.maxVelocity,
		r.dtOptim,
		Dt,
		r.ratioOfBodiesWithCapVel)

	
	fmt.Printf( r.Status())
	
	t1 := time.Now()
	StepDuration = float64((t1.Sub(t0)).Nanoseconds())
	Gflops = float64( nbComputationPerStep) /  StepDuration
}
var Gflops float64
var StepDuration float64

func (r * Run) Status() string {

	return r.status
}

// compute repulsive forces by spreading the calculus
// among nbRoutine go routines
func (r * Run) ComputeRepulsiveForceConcurrent(nbRoutine int) {

	sliceLen := len(*r.bodies)
	misDistance := make( chan float64)

	// breakdown slice
	for i:=0; i<nbRoutine; i++ {
	
		startIndex := (i*sliceLen)/nbRoutine
		endIndex := ((i+1)*sliceLen)/nbRoutine -1
		go r.ComputeRepulsiveForceSubSetMinDist( startIndex, endIndex, misDistance)
	}

	// wait for return
	for i:=0; i<nbRoutine; i++ {
		<- misDistance
	}

}

// compute repulsive forces
func (r * Run) ComputeRepulsiveForce() {
	
	r.ComputeRepulsiveForceSubSet( 0, len(*r.bodies))
}

// compute repulsive forces for a sub part of the bodies
func (r * Run) ComputeRepulsiveForceSubSetMinDist( startIndex, endIndex int, minDistance chan<- float64) {
	
	var _minDistance float64
	r.ComputeRepulsiveForceSubSet( startIndex, endIndex)
	minDistance <- 	_minDistance
}
// compute repulsive forces for a sub part of the bodies
// return the minimal distance between the bodies sub set
func (r * Run) ComputeRepulsiveForceSubSet( startIndex, endIndex int) float64 {

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
	return 0.0
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
// return min distance
func (r * Run) computeAccelerationOnBodyBarnesHut(idx int) float64 {

	// reset acceleration
	acc := &((*r.bodiesAccel)[idx])
	acc.X = 0
	acc.Y = 0
	
	// Coord is initialized at the Root coord
	var rootCoord quadtree.Coord
	
	return r.computeAccelationWithNodeRecursive( idx, rootCoord)
}

// given a body and a node in the quadtree, compute the repulsive force
func (r * Run) computeAccelationWithNodeRecursive( idx int, coord quadtree.Coord) float64 {
	
	body := (*r.bodies)[idx]
	acc := &((*r.bodiesAccel)[idx])
	
	// compute the node box size
	level := coord.Level()
	boxSize := 1.0 / math.Pow( 2.0, float64(level)) // if level = 0, this is 1.0
	
	node := & (r.q.Nodes[coord])
	dist := getModuloDistanceBetweenBodies( &body, &(node.Body))

	
	// avoid node with zero mass
	if( node.M == 0) {
		return 0.0
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
	
					dist := getModuloDistanceBetweenBodies( &body, b)
					if dist < r.minInterBodyDistance { 
						var m sync.Mutex
						m.Lock()
						r.minInterBodyDistance = dist 
						m.Unlock()	
					}
					
					x, y := getRepulsionVector( &body, b)
			
					acc.X += x
					acc.Y += y
					rank++
					// fmt.Printf("computeAccelationWithNodeRecursive at leaf %#v rank %d x %9.3f y %9.3f\n", b.Coord(), rank, x, y)
				}
			}
		}
	}
	return 0.0
}

func (r * Run) UpdateVelocity() {

	var nbVelCapping int64
	// parse all bodies
	for idx, _ := range (*r.bodies) {

		// update velocity (to be completed with Dt)
		acc := r.getAcc(idx)
		vel := r.getVel(idx)
		vel.X += acc.X * G * Dt
		vel.Y += acc.Y * G * Dt
		
		// put some drag
		vel.X *= SpeedDragFactor
		vel.Y *= SpeedDragFactor
		
		// if velocity is above
		velocity := math.Sqrt( vel.X*vel.X + vel.Y*vel.Y)
		
		if velocity > r.maxVelocity {
			var m sync.Mutex
			m.Lock()
			r.maxVelocity = velocity
			m.Unlock()
		}

		if velocity*Dt > MaxDisplacement { 
			vel.X *= MaxDisplacement/(velocity*Dt)
			vel.Y *= MaxDisplacement/(velocity*Dt)
			nbVelCapping += 1
		}
	}
	r.ratioOfBodiesWithCapVel = float64(nbVelCapping) / float64(len(*r.bodies))
}

func (r * Run) UpdatePosition() {

	// compute optimal Dt, where we want the move to be
	// half of the minimum distance between bodies 
	r.dtOptim = 0.5 * (r.minInterBodyDistance / r.maxVelocity)

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
// return x, y of repulsion vector and distance between A & B
func getRepulsionVector( A, B *quadtree.Body) (x, y float64) {

	atomic.AddUint64( &nbComputationPerStep, 1)

	x = getModuloDistance( B.X, A.X)
	y = getModuloDistance( B.Y, A.Y)

	distQuared := (x*x + y*y)
	absDistance := math.Sqrt( distQuared + ETA )
	
	distPow3 := (distQuared + ETA) * absDistance
	
	if false { 
		distPow3 := math.Pow( distQuared, 1.5) 
		distQuared /= distPow3
	}
	
	// repulsion is proportional to mass
	massCombined := A.M * B.M
	x *= massCombined
	y *= massCombined

	// repulsion is inversly proportional to the square of the distance (1/r2)
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

// get modulo distance between alpha and beta in a given village.
//
// alpha and beta are between left and rigth
// the modulo distance cannot be above (rigth-left) /2
func getModuloDistanceLocal( alpha, beta, left, rigth float64) (dist float64) {

	dist = beta-alpha
	maxDist := rigth-left
	halfMaxDist := maxDist/2.0
	if( dist > halfMaxDist ) { dist -= maxDist }
	if( dist < -halfMaxDist ) { dist += maxDist }
		
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

var CurrentCountry = "bods"
