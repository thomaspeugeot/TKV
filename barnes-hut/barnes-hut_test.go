package barnes_hut

import (
	"os"
	"github.com/thomaspeugeot/tkv/quadtree"
	"testing"
	"math"
)

// test gif output
func TestOutputGif(t *testing.T) {

	bodies := make([]quadtree.Body, 2000)
	SpreadOnCircle( & bodies)
	
	var r Run
	r.Init( & bodies)
	
	var output *os.File
	output, _ = os.Create("essai.gif")
	
	r.SetState( RUNNING)
	r.OutputGif( output, 0)
	// visual verification
}

func TestOneStep(t *testing.T) {
	bodies := make([]quadtree.Body, 2000)
	SpreadOnCircle( & bodies)
	
	var r Run
	r.Init( & bodies)
	// r.q.CheckIntegrity( t)
	r.OneStep()
	r.OneStep()
	r.q.CheckIntegrity( t)
}

func TestGetModuloDistance(t *testing.T) {
	
	cases := []struct {
		x1, x2, want float64
	}{
		{0.0, 0.1, 0.1},
		{0.0, 0.0, 0.0},
	}
	for _, c := range cases {
		got := getModuloDistance( c.x1, c.x2)
		if( got != c.want) {
			t.Errorf("x1 %f x2 %f == %f, want %f", c.x1, c.x2, got, c.want)
		}	
	}
}

func TesGetRepulsionVector(t *testing.T) {
	
	cases := make( []struct {
		A, B quadtree.Body
		wantX, wantY float64
	}, 1)
	
	bodies := make([]quadtree.Body, 2)
	bodies[1].X = 0.4
	bodies[1].Y = 0.3
	
	cases[0].B = bodies[1]
	cases[0].wantX = -3.2
	cases[0].wantY = -2.4
	
	for _, c := range cases {
		gotX, gotY := getRepulsionVector( & c.A, & c.B)
		if( (gotX != c.wantX) && (gotY != c.wantY)) {
			t.Errorf("A %#v B %#v == %f %f, want %f %f", c.A, c.B, gotX, gotY, c.wantX, c.wantY )
		}	
	}
}

// test the step
func TesComputeRepulsiveForces(t *testing.T) {
	
	cases := make( []struct {
		r Run
		want [2]Acc }, 1) // 1 case
	
	bodies := make([]quadtree.Body, 2)
	bodies[1].X = 0.4
	bodies[1].Y = 0.3
	cases[0].r.Init( & bodies)
	
	cases[0].want[0] = Acc{-3.2, -2.4}
	cases[0].want[1] = Acc{3.2, 2.4}
	
	for _, c := range cases {
		c.r.ComputeRepulsiveForce()
		if( *(c.r.getAcc(0)) != c.want[0] && *(c.r.getAcc(1)) != c.want[1]) {
			t.Errorf("\ngot %#v %#v\nwant %#v %#v", c.r.getAcc(0), c.r.getAcc(1), c.want[0], c.want[1])
		}
	}
}

// func test the concurrent version is the same as the serial version
func TestComputeRepulsiveForcesConcurrent(t *testing.T) {
	
	bodies := make([]quadtree.Body, 10 * 10)
	bodies2 := make([]quadtree.Body, 10 * 10)
	SpreadOnCircle( & bodies)
	copy( bodies2, bodies)
	var r, r2 Run
	r.Init( & bodies)
	r2.Init( & bodies2)
	r.ComputeRepulsiveForce()
	r2.ComputeRepulsiveForceConcurrent(13)

	same := true
	for idx, _ := range *r.bodies {
		if (*r.bodies)[idx].X != (*r2.bodies)[idx].X { same = false}
		if (*r.bodies)[idx].Y != (*r2.bodies)[idx].Y { same = false}
	}
	if ! same {
		t.Errorf("different results")
	}
}

// test that the barnes hut computation of repulsive
// forces is close to the classic computation 
func TestComputeAccelerationOnBodyBarnesHut(t *testing.T) {

	bodies := make([]quadtree.Body, 2000000)
	SpreadOnCircle( & bodies)
	// quadtree.InitBodiesUniform( &bodies, 200)
	var r Run
	r.Init( & bodies)
	
	r.computeAccelerationOnBody( 0)
	accReference := (*r.bodiesAccel)[0]
	
	r.computeAccelerationOnBodyBarnesHut( 0)
	accBH := (*r.bodiesAccel)[0]

	r.q.CheckIntegrity(t)
	
	// measure the difference between reference and BH
	accReferenceLength := math.Hypot( accReference.X, accReference.Y)
	diff := math.Hypot( (accReference.X - accBH.X), (accReference.Y - accBH.Y))
	
	relativeError := diff/accReferenceLength
	
	// tolerance is arbitrary set
	tolerance := 0.02 // 5%
	if( relativeError > tolerance) {
		t.Errorf("different results, accel ref x %f y %f, got x %f y %f", accReference.X, accReference.Y, accBH.X, accBH.Y)	
		t.Errorf("different results, expected less than %f, got %f", tolerance, relativeError)
	}
	
}

// test wether the computation of min distance is equal between
// a mutex approach or a concurrent approach
//
// reference failure :
// barnes-hut_test.go:178: different results for concurrent computation 6.064784e-05, with mutex 1.448228e-04
func TestConcurrentMinDistanceCompute( t *testing.T) {
	bodies := make([]quadtree.Body, 2000)
	SpreadOnCircle( & bodies)
	
	var r Run
	r.Init( & bodies)
	
	// init
	r.minInterBodyDistance = 2.0 // cannot be in a 1.0 by 1.0 square
	r.q.UpdateNodesListsAndCOM()
	minDistance := r.ComputeRepulsiveForceConcurrent(20)
	if( r.minInterBodyDistance == 0) {
		t.Errorf("minInterBodyDistance is 0.0")
	}
	
	if( minDistance == 0) {
		t.Errorf("minDistance is 0.0")
	}
	if( minDistance != r.minInterBodyDistance) {
		t.Errorf("different results for concurrent computation %e, with mutex %e", minDistance, r.minInterBodyDistance)
	}

}
