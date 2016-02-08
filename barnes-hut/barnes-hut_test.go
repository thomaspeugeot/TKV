package barnes_hut

import (
	"os"
	"github.com/thomaspeugeot/tkv/quadtree"
	"testing"
	"math"
	"math/rand"
)

// test gif output
func TestOutputGif(t *testing.T) {

	bodies := make([]quadtree.Body, 100)
	spreadOnCircle( & bodies)
	
	var r Run
	r.Init( & bodies)
	
	var output *os.File
	output, _ = os.Create("essai.gif")
	
	r.outputGif( output, 20)
	// visual verification
}

func Test20Steps(t *testing.T) {
	
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
	spreadOnCircle( & bodies)
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
// function used to spread bodies randomly on 
// the unit square
func spreadOnCircle(bodies * []quadtree.Body) {
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
