/*
Package barnes-hut is a modified barnes-hut simulation algorithm witg concurrent computation.

A the start of the simulation, bodies are spread according to the
density of the country of interest.
At the end of the simulation, bodies are spread evenly on a 2D rectangle. At the end of the simulation,
the body repartition is said to be hyperuniform (https://www.quantamagazine.org/hyperuniformity-found-in-birds-math-and-physics-20160712/)

Barnes-Hut is an embarisgly parallel algorithm. This implementation is used the concurrent model of the go langage.
Nb of conurrent routine can be set up dynamicaly.

TKV implementation starts from a Barnes-Hut implementation of the gravitation simulation
and make the following modification:

In a gravitational simulation, bodies are attracted to each others by mean of the newtownian gravitation law.
Here, repulsion is used instead of gravitational attraction.
A drag is introduced to put simulation bodies to rest (see #Run.UpdateVelocity)

In a cosmological simulation, bodies position are not limited. Here,
bodies are kept within a [0;1]*[0;1] square by having "mirror" bodies that
forbids a body from crossing the border (see #Run.UpdatePosition)
*/
package barneshut

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/thomaspeugeot/tkv/quadtree"
)

// constant to be added to the distance between bodies
// in order to compute repulsion (avoid near infinite repulsion force)
// note : declaring those variable as constant has no impact on benchmarks results
//var	ETA float64 = 1e-10
var ETA float64 = 0.0

// pseudo gravitational constant to compute
//should have no effect on the simulation since Dt is computed according to computed acceleration
var G float64 = 0.01

//var Dt float64  = 3*1e-8 // difficult to fine tune
var Dt float64 = 2.3 * 1e-10 // difficult to fine tune
var DtRequest = Dt           // new value of Dt requested by the UI. The real Dt will be changed at the end of the current step.

// velocity cannot be too high in order to stop bodies from overtaking
// each others. The variable defines when displacement has to be capped
var MaxDisplacement float64 = 0.001 // cannot make more that 1/1000 th of the unit square per second

// give the theorical max displacement, as a ratio between the min distance and the max speed
// note: 2.0 gives too much stirring
var MaxRatioDisplacement float64 = 0.5

// BN_THETA is the barnes hut criteria
// If distance to Center Of Mass (COM) of the box is more than BN_THETA * side of the box
// then, compute force from the whole box instead computing for each node of the box.
var BN_THETA float64 = 0.5

// New value of theta requested by the UI. The real BN_THETA will be changed at the end of the current step.
var BN_THETA_Request = BN_THETA

// how much drag we put (1.0 is no drag)
// tHis criteria is important because it favors bodies that moves freely against bodies that are stuck on a border
// 0.99 makes a very bumpy behavior for the Dt
var SpeedDragFactor float64 = 0.2

// used to compute speed up
var nbComputationPerStep uint64

// if true, Barnes-Hut algo is used
var UseBarnesHut bool = true

// Cutoff for influence.
// According to Bartolo, 1/r2 is very strong on a plane (integration does not convergence on the plane)
// therefore a cutoff distance for the computation of the force is wellcome
// first try at 1/10 th
var CutoffDistance float64 = 1.0

// At what step do the simulation stop
var ShutdownCriteria float64 = 0.00001

// Bodies's X,Y position coordinates are float64 between 0 & 1
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

// decides wether Dt is set manual or automaticaly
type DtAdjustModeType string

var DtAdjustMode DtAdjustModeType

// Possible values for DtAdjustModeType
const (
	AUTO   = "AUTO"
	MANUAL = "MANUAL"
)

// state of the simulation
type State string

// Possible values for simulation State
const (
	STOPPED = "STOPPED"
	RUNNING = "RUNNING"
)

// decide wether, villages borders are drawn
type RenderState string

// Possible values for RenderState
const (
	WITHOUT_BORDERS = "WITHOUT_BORDERS"
	WITH_BORDERS    = "WITH_BORDERS"
)

var ratioOfBorderVillages = 0.0 // ratio of villages that are eligible for marking a border

// decide wether to display the original configuration or the running configruation
type RenderChoice string

// Possible  values for RenderChoice
const (
	ORIGINAL_CONFIGURATION = "ORIGINAL_CONFIGURATION"
	RUNNING_CONFIGURATION  = "RUNNING_CONFIGURATION"
)

// set the number of concurrent routine for the physic calculation
// this value can be set interactively during the run
var ConcurrentRoutines int = 100

// number of village per X or Y axis. For 10 000 villages, this number is 100
// this value can be set interactively during the run
var nbVillagePerAxe int = 100

// a simulation run
type Run struct {
	bodies               *[]quadtree.Body // bodies position in the quatree
	bodiesOrig           *[]quadtree.Body // original bodies position in the quatree
	bodiesAccel          *[]Acc           // bodies acceleration
	bodiesVel            *[]Vel           // bodies velocity
	bodiesEnergy         *[]float64       // bodies energy
	bodiesNeighbours     *NeighbourDico   // storage for neighbour of all bodies
	bodiesNeighboursOrig *NeighbourDico   // storage for neighbour of all bodies at init

	q       quadtree.Quadtree // the supporting quadtree
	country string            // the country of interest
	state   State
	step    int

	giniFileLog            *os.File
	giniOverTime           [][]float64 // evolution of the gini distribution over time
	xMin, xMax, yMin, yMax float64     // coordinates of the rendering windows
	renderState            RenderState
	renderChoice           RenderChoice

	fieldRendering bool // if true, render the repulsion field
	gridFieldNb    int  // nb of ticks for the field render area

	minInterBodyDistance    float64           // computed at each step (to compute optimal DT value)
	maxRepulsiveForce       MaxRepulsiveForce // computed at each step (to compute optimal DT value)
	maxVelocity             float64           // max velocity
	dtOptim                 float64           // optimal dt
	ratioOfBodiesWithCapVel float64           // ratio of bodies where the speed has been capped
	energy                  float64           // total repulsive energy
	energyDecreaseRatio     float64           // energy decrease ratio. Is used as a shutdown criteria

	status string // status of the run

	borderHasBeenMet bool // compute wether the border has been met (see issue#4)

	OutputDir     string // output dir for the run
	StatusFileLog *os.File

	CaptureGifStep int // simulaton steps between gif generation
}

// (in order to solve issue "over accumulation of bodies at border slows dow spreading #5")
var maxMinInterBodyDistance float64 // this variable store the max of the previous value

// RunSimulation is the main entry to the simulation
// It call for one step of simulation until the
// energy decrease ratio is met
func (r *Run) RunSimulation() {

	Info.Printf("Energy decrease %f, ShutdownCriteria %f crit %t", r.energyDecreaseRatio, ShutdownCriteria, (r.energyDecreaseRatio > ShutdownCriteria))
	for r.energyDecreaseRatio > ShutdownCriteria {
		// if state is STOPPED, pause
		for r.state == STOPPED {
			time.Sleep(100 * time.Millisecond)
		}
		r.OneStep()
	}

	r.state = STOPPED
	r.CaptureConfig()
	// r.CreateMovieFromGif()
	os.Exit(0)
}

func (r *Run) SetCountry(country string) {
	r.country = country
}

func (r *Run) SetGridFieldNb(v int) {

	renderingMutex.Lock()
	r.gridFieldNb = v
	// Trace.Printf("r.gridFieldNb %d", r.gridFieldNb)
	renderingMutex.Unlock()
}

func NewRun() *Run {
	var r Run
	r.state = STOPPED
	r.gridFieldNb = 10
	bodies := make([]quadtree.Body, 0)

	// create output directory and cwd to it
	// https://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format
	r.OutputDir = time.Now().Local().Format("2006_01_02_150405")
	Info.Printf("Output dir %s", r.OutputDir)
	os.Mkdir(r.OutputDir, 0777)

	// init the file storing the gini distribution over time
	filename := fmt.Sprintf(r.OutputDir + "/gini_out.csv")
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	r.giniFileLog = file

	// init the file storing status of the run at all steps
	filename = fmt.Sprintf(r.OutputDir + "/status_out.csv")
	file, err = os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	r.StatusFileLog = file

	r.Init(&bodies)

	return &r
}

// init the run with an array of quadtree bodies
func (r *Run) Init(bodies *([]quadtree.Body)) {

	Trace.Printf("Init begin")

	r.bodies = bodies

	makeBodiesMemory := func(varAddress **[]quadtree.Body) {
		tmp := make([]quadtree.Body, len(*bodies))
		*varAddress = &tmp
	}

	// create a reference of the bodies
	makeBodiesMemory(&r.bodiesOrig)
	copy(*r.bodiesOrig, *r.bodies)

	acc := make([]Acc, len(*bodies))
	vel := make([]Vel, len(*bodies))
	energy := make([]float64, len(*bodies))
	r.bodiesAccel = &acc
	r.bodiesVel = &vel
	r.bodiesEnergy = &energy

	// init quatrees
	r.q.Init(bodies)

	// init neighbour array
	r.InitNeighbourDico(bodies)

	r.state = STOPPED
	r.SetRenderingWindow(0.0, 1.0, 0.0, 1.0)
	r.renderState = WITH_BORDERS           // we draw borders
	r.renderChoice = RUNNING_CONFIGURATION // we draw borders
	r.fieldRendering = false

	r.energy = math.MaxFloat64 // very high
	r.energyDecreaseRatio = 1.0

	DtAdjustMode = AUTO

	Trace.Printf("Init end")
}

// rendering the data set can be done only outside the load config xxx function
var renderingMutex sync.Mutex

func (r *Run) GetMaxRepulsiveForce() MaxRepulsiveForce {
	return r.maxRepulsiveForce
}

func (r *Run) RatioOfBodiesWithCapVel() float64 {
	// log.Output( 1, fmt.Sprintf( "ratioOfBodiesWithCapVel %f ", r.ratioOfBodiesWithCapVel))
	return r.ratioOfBodiesWithCapVel
}

func (r *Run) getAcc(index int) *Acc {
	return &(*r.bodiesAccel)[index]
}

func (r *Run) getVel(index int) *Vel {
	return &(*r.bodiesVel)[index]
}

func (r *Run) State() State {
	return r.state
}

func (r *Run) GetStep() int {
	return r.step
}

func (r *Run) SetState(s State) {
	r.state = s
}

func (r *Run) SetRenderingWindow(xMin, xMax, yMin, yMax float64) {
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

func SetRatioBorderBodies(ratioOfBorderVillages_p float64) {
	ratioOfBorderVillages = ratioOfBorderVillages_p
}

func (r *Run) GetMinInterBodyDistance() float64 {
	return r.minInterBodyDistance
}

func (r *Run) GiniOverTimeTransposed() [][]float64 {

	var giniOverTimeTransposed [][]float64
	// := r.giniOverTime
	giniOverTimeTransposed = transposeFloat64(r.giniOverTime)
	return giniOverTimeTransposed
}

func (r *Run) GiniOverTime() [][]float64 {

	return r.giniOverTime
}

func (r *Run) ToggleRenderChoice() {
	if r.renderChoice == RUNNING_CONFIGURATION {
		r.renderChoice = ORIGINAL_CONFIGURATION
	} else {
		r.renderChoice = RUNNING_CONFIGURATION
	}
}

func (r *Run) ToggleFieldRendering() {
	Info.Printf("ToggleFieldRendering new state %v", !r.fieldRendering)
	r.fieldRendering = !r.fieldRendering
}

func (r *Run) ToggleManualAuto() {
	if DtAdjustMode == MANUAL {
		DtAdjustMode = AUTO
	} else {
		DtAdjustMode = MANUAL
	}
}

func (r *Run) OneStep() {
	// renderingMutex.Lock()
	r.OneStepOptional(true)
	// renderingMutex.Unlock()
}
func (r *Run) OneStepOptional(updatePosition bool) {

	// serialize into a file the gif
	if r.step%r.CaptureGifStep == 0 {
		r.CaptureGif()
	}

	t0 := time.Now()

	// compute gini distribution
	r.q.ComputeQuadtreeGini()

	// append the new gini elements
	// create the array
	giniArray := make([]float64, 10)
	copy(giniArray, r.q.BodyCountGini[8][:])
	r.giniOverTime = append(r.giniOverTime, giniArray)

	//
	w := csv.NewWriter(r.giniFileLog)
	// w.Comma = '\t'

	recordStr := make([]string, 10)
	// get gini distribution at level 8
	density := r.ComputeDensityTencilePerTerritory()
	// for idx, record := range r.q.BodyCountGini[8][:] {
	for idx, record := range density {
		recordStr[idx] = fmt.Sprintf("%f", record)

	}
	if errCsv := w.Write(recordStr); errCsv != nil {
		log.Fatalln("error writing record to csv:", errCsv)
	}
	w.Flush()

	nbComputationPerStep = 0
	r.maxVelocity = 0.0

	BN_THETA = BN_THETA_Request

	// compute the quadtree from the bodies
	r.q.UpdateNodesListsAndCOM()

	// compute repulsive forces & acceleration
	r.ComputeRepulsiveForceConcurrent(ConcurrentRoutines)
	r.ComputeMaxRepulsiveForce()

	// Trace.Printf("MaxRepulsiveForce %#v", r.maxRepulsiveForce)

	// compute optimal Dt, where we want the move to be
	// half of the minimum distance between bodies
	// with initial speed at 0, the speed will increase to Dt*Acc and
	// the displacement will be Dx = Dt*Dt*Acc. Thefore Dt
	r.dtOptim = math.Sqrt(MaxRatioDisplacement * r.minInterBodyDistance / r.maxRepulsiveForce.Norm)

	// update Dt according to request or according to computing optimal Dt
	if DtAdjustMode == MANUAL {
		Dt = DtRequest
	} else {
		if r.dtOptim > 0.0 {
			Dt = r.dtOptim
		}
	}

	// compute velocity
	r.UpdateVelocity()

	// compute new position
	if updatePosition {
		r.UpdatePosition()
	}

	// init neighbours original
	if r.step == 0 {
		r.bodiesNeighboursOrig.Copy(r.bodiesNeighbours)
	}

	// compute total energy
	// parse all bodies
	lastEnergy := r.energy
	r.energy = 0.0
	for idx := range *r.bodies {
		e := &((*r.bodiesEnergy)[idx])
		r.energy += *e
	}
	r.energyDecreaseRatio = (lastEnergy - r.energy) / lastEnergy

	// compute stirring
	stirring := r.bodiesNeighbours.ComputeStirring(r.bodiesNeighboursOrig)
	ratioOfNil := r.bodiesNeighbours.ComputeRatioOfNilNeighbours()

	// update the step
	r.step++

	t1 := time.Now()
	StepDuration = float64((t1.Sub(t0)).Nanoseconds())
	Gflops = float64(nbComputationPerStep) / StepDuration

	r.status = fmt.Sprintf("%s step %5d speedup %5.2f Dur %4.2f E %e MinD %e MaxMinD %e MaxV %e Dt Opt %e Dt %e F/A %e stirring %f nils %f e decr ratio %1.10f\n",
		time.Now().Local().Format("2006 01 02 15 04 05"),
		r.step,
		float64(len(*r.bodies)*len(*r.bodies))/float64(nbComputationPerStep), //speedup
		StepDuration/1000000000, // duration in seconds
		r.energy,                // energy
		r.minInterBodyDistance,
		maxMinInterBodyDistance,
		r.maxVelocity,
		r.dtOptim,
		Dt,
		r.ratioOfBodiesWithCapVel,
		stirring,
		ratioOfNil,
		r.energyDecreaseRatio)

	fmt.Fprintf(r.StatusFileLog, r.Status())
	fmt.Printf(r.Status())
}

var Gflops float64
var StepDuration float64

func (r *Run) Status() string {

	return r.status
}

// compute repulsive forces by spreading the calculus
// among nbRoutine go routines
//
// return minInterbodyDistance
func (r *Run) ComputeRepulsiveForceConcurrent(nbRoutine int) float64 {

	Trace.Println("ComputeRepulsiveForceConcurrent")
	sliceLen := len(*r.bodies)
	minInterbodyDistanceChan := make(chan float64)

	// breakdown slice
	for i := 0; i < nbRoutine; i++ {

		startIndex := (i * sliceLen) / nbRoutine
		endIndex := ((i + 1) * sliceLen) / nbRoutine

		// if nbRoutine is above sliceLen
		if endIndex < startIndex {
			endIndex = startIndex
		}

		// log.Printf( "started routine %3d\n", i)
		go r.ComputeRepulsiveForceSubSetMinDist(startIndex, endIndex, minInterbodyDistanceChan)
	}

	// wait for return and compute the min distance across all routines
	r.minInterBodyDistance = 2.0 // cannot be in a 1.0 by 1.0 square
	for i := 0; i < nbRoutine; i++ {
		// log.Printf( "waiting routine %3d\n", i)

		minInterbodyDistanceRoutine := <-minInterbodyDistanceChan
		// log.Printf( "routine %3d minInterbodyDistance by mutex %e, by concurency %e\n", i, r.minInterBodyDistance, minInterbodyDistanceRoutine)

		if minInterbodyDistanceRoutine < r.minInterBodyDistance {
			r.minInterBodyDistance = minInterbodyDistanceRoutine
		}
	}
	// log.Printf( "minInterbodyDistance by mutex %e, by concurency %e\n", r.minInterBodyDistance, minInterbodyDistance)

	// update maxMinInterBodyDistance
	if maxMinInterBodyDistance == 0 {
		maxMinInterBodyDistance = r.minInterBodyDistance
	}
	if maxMinInterBodyDistance < r.minInterBodyDistance {
		maxMinInterBodyDistance = r.minInterBodyDistance
	}

	return r.minInterBodyDistance
}

// compute repulsive forces
func (r *Run) ComputeRepulsiveForce() {

	r.ComputeRepulsiveForceSubSet(0, len(*r.bodies))
}

// compute repulsive forces for a sub part of the bodies
//
// send the computed min distance through minInterbodyDistanceChan
func (r *Run) ComputeRepulsiveForceSubSetMinDist(startIndex, endIndex int, minInterbodyDistanceChan chan<- float64) {

	_minInterbodyDistance := r.ComputeRepulsiveForceSubSet(startIndex, endIndex)
	// log.Printf( "minInterbodyDistance by mutex %e, by concurency %e\n", r.minInterBodyDistance, _minInterbodyDistance)

	minInterbodyDistanceChan <- _minInterbodyDistance
}

// compute repulsive forces for a sub part of the bodies
// return the minimal distance between the bodies sub set
func (r *Run) ComputeRepulsiveForceSubSet(startIndex, endIndex int) float64 {

	// Trace.Printf("ComputeRepulsiveForceSubSet %d %d", startIndex, endIndex)
	minInterbodyDistance := 2.0

	// parse all bodies
	bodiesSubSet := (*r.bodies)[startIndex:endIndex]
	for idx := range bodiesSubSet {

		// index in the original slice
		origIndex := idx + startIndex

		minInterbodyDistanceSubSet := 2.0
		if UseBarnesHut {
			minInterbodyDistanceSubSet = r.computeAccelerationOnBodyBarnesHut(origIndex)
		} else {
			minInterbodyDistanceSubSet = r.computeAccelerationOnBody(origIndex)
		}
		if minInterbodyDistanceSubSet < minInterbodyDistance {
			minInterbodyDistance = minInterbodyDistanceSubSet
		}
	}
	return minInterbodyDistance
}

// parse all other bodies to compute acceleration
func (r *Run) computeAccelerationOnBody(origIndex int) float64 {

	minInterbodyDistance := 2.0

	body := (*r.bodies)[origIndex]

	// reset acceleration
	acc := &((*r.bodiesAccel)[origIndex])
	acc.X = 0
	acc.Y = 0

	// reset energy
	energy := &((*r.bodiesEnergy)[origIndex])
	*energy = 0.0

	// parse all other bodies for repulsions
	// accumulate repulsion on acceleration
	for idx2 := range *r.bodies {

		if idx2 != origIndex {
			body2 := (*r.bodies)[idx2]

			dist := getDistanceBetweenBodiesWithMirror(&body, &body2, 0, 0)

			if dist == 0.0 {
				log.Fatal("distance is 0.0 between ", body, " and ", body2)
			}

			if dist < minInterbodyDistance {
				minInterbodyDistance = dist
			}

			x, y, e := getRepulsionVector(&body, &body2, 0, 0)

			acc.X += x
			acc.Y += y
			*energy += e
			// Trace.Printf("computeAccelerationOnBody idx2 %3d x %9.3f y %9.3f \n", idx2, x, y)
		}
	}
	return minInterbodyDistance
}

// with body index idx, parse all other bodies to compute acceleration
// with the barnes-hut algorithm
//
// return min distance between body and other bodies
func (r *Run) computeAccelerationOnBodyBarnesHut(idx int) float64 {

	// reset acceleration
	acc := &((*r.bodiesAccel)[idx])
	acc.X = 0
	acc.Y = 0

	// reset energy
	energy := &((*r.bodiesEnergy)[idx])
	*energy = 0.0

	// Coord is initialized at the Root coord
	var rootCoord quadtree.Coord

	result := 2.0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {

			resultTmp := r.computeAccelationWithNodeRecursive(idx, rootCoord, i, j)
			if resultTmp < result {
				result = resultTmp
			}
		}
	}

	return result
}

// given a body at index idx, and a node at coordinate coord in the q quadtree,
// compute the repulsive force and update the accelation at index idx
// x and y and the mirror configuration (xM==yM==0 means no)
//
// return the minditance between the body and bodies in the quadtree node
func (r *Run) computeAccelationWithNodeRecursive(idx int, coord quadtree.Coord, xM, yM int) float64 {

	// can be useful to debug  -->
	// time.Sleep(300 * time.Millisecond)
	// Info.Printf("computeAccelationWithNodeRecursive idx %d coord %#v", idx, coord)
	minInterbodyDistance := 2.0

	body := (*r.bodies)[idx]
	acc := &((*r.bodiesAccel)[idx])
	energy := &((*r.bodiesEnergy)[idx])

	// compute the node box size
	level := coord.Level()
	boxSize := 1.0 / math.Pow(2.0, float64(level)) // if level = 0, this is 1.0

	// fetch node in the quadtree
	node := &(r.q.Nodes[coord])
	distToNode := getDistanceBetweenBodiesWithMirror(&body, &(node.Body), xM, yM)

	// Info.Printf("computeAccelationWithNodeRecursive distance to quadtree node %f", distToNode)

	// avoid node with zero mass
	if node.M == 0 {
		return 2.0
	}

	// check if the COM of the node can be used
	if (boxSize / distToNode) < BN_THETA {

		x, y, e := getRepulsionVector(&body, &(node.Body), xM, yM)

		acc.X += x
		acc.Y += y
		*energy += e

		// Info.Printf("computeAccelationWithNodeRecursive at node %#v x %9.3f y %9.3f\n", node.Coord(), x, y)

	} else {
		if level < 8 {
			// parse sub nodes
			// Info.Printf("computeAccelationWithNodeRecursive go down at node %#v\n", node.Coord())
			coordNW, coordNE, coordSW, coordSE := quadtree.NodesBelow(coord)
			dist := 2.0
			dist = r.computeAccelationWithNodeRecursive(idx, coordNW, xM, yM)
			if dist < minInterbodyDistance {
				minInterbodyDistance = dist
			}

			dist = r.computeAccelationWithNodeRecursive(idx, coordNE, xM, yM)
			if dist < minInterbodyDistance {
				minInterbodyDistance = dist
			}

			dist = r.computeAccelationWithNodeRecursive(idx, coordSW, xM, yM)
			if dist < minInterbodyDistance {
				minInterbodyDistance = dist
			}

			dist = r.computeAccelationWithNodeRecursive(idx, coordSE, xM, yM)
			if dist < minInterbodyDistance {
				minInterbodyDistance = dist
			}

		} else {

			// parse bodies of the node
			rank := 0
			rankOfBody := -1
			for b := node.First(); b != nil; b = b.Next() {
				if *b != body {

					// Info.Printf("computeAccelationWithNodeRecursive at leaf %#v rank %d", b.Coord(), rank)
					dist := getDistanceBetweenBodiesWithMirror(&body, b, xM, yM)

					r.bodiesNeighbours.Insert(idx, b, dist)

					if dist == 0.0 {
						var t testing.T
						r.q.CheckIntegrity(&t)

						Error.Printf("Problem body x %f y %f to x %f y %f", body.X, body.Y, b.X, b.Y)
						Error.Printf("Problem at rank %d for body of rank %d on node %#v ",
							rank, rankOfBody, *node)
						os.Exit(3)
					}
					if dist < minInterbodyDistance {
						minInterbodyDistance = dist
					}

					x, y, e := getRepulsionVector(&body, b, xM, yM)
					// Info.Printf("computeAccelationWithNodeRecursive at leaf %#v rank %d x %9.3f y %9.3f\n", b.Coord(), rank, x, y)

					acc.X += x
					acc.Y += y
					*energy += e
					rank++
					rankOfBody++
				} else {
					rankOfBody = rank
				}
			}
		}
	}
	return minInterbodyDistance
}

func (r *Run) UpdateVelocity() {

	Trace.Println("UpdateVelocity")

	var nbVelCapping int64
	// parse all bodies
	for idx := range *r.bodies {

		// put some drag on initial speed
		vel := r.getVel(idx)
		vel.X *= SpeedDragFactor
		vel.Y *= SpeedDragFactor

		// update velocity (to be completed with Dt)
		acc := r.getAcc(idx)
		vel.X += acc.X * Dt
		vel.Y += acc.Y * Dt

		// if velocity is above
		velocity := math.Sqrt(vel.X*vel.X + vel.Y*vel.Y)

		if velocity > r.maxVelocity {
			var m sync.Mutex
			m.Lock()
			r.maxVelocity = velocity
			m.Unlock()
		}

		if velocity*Dt > MaxDisplacement {
			vel.X *= MaxDisplacement / (velocity * Dt)
			vel.Y *= MaxDisplacement / (velocity * Dt)
			nbVelCapping += 1
		}
	}
	r.ratioOfBodiesWithCapVel = float64(nbVelCapping) / float64(len(*r.bodies))
}

func (r *Run) UpdatePosition() {

	r.borderHasBeenMet = false

	// parse all bodies
	for idx := range *r.bodies {

		body := &((*r.bodies)[idx])

		// updatePos
		vel := r.getVel(idx)

		body.X += vel.X * Dt
		body.Y += vel.Y * Dt

		if body.X >= 1.0 {
			body.X = 1.0 - (body.X - 1.0)
			vel.X = -vel.X
			r.borderHasBeenMet = true
		}
		if body.X <= 0.0 {
			body.X = -body.X
			vel.X = -vel.X
			r.borderHasBeenMet = true
		}
		if body.Y >= 1.0 {
			body.Y = 1.0 - (body.Y - 1.0)
			vel.Y = -vel.Y
			r.borderHasBeenMet = true
		}
		if body.Y <= 0.0 {
			body.Y = -body.Y
			vel.Y = -vel.Y
			r.borderHasBeenMet = true
		}
	}
}

func (r *Run) BodyCountGini() quadtree.QuadtreeGini {
	return r.q.BodyCountGini
}

var CurrentCountry = "bods"
