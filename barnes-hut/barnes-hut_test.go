package barnes_hut

import (
	"os"
	"github.com/thomaspeugeot/tkv/quadtree"
	"testing"
	"fmt"
	"math"
	"math/rand"
)

// init 
func TestOutputGif(t *testing.T) {

	bodies := make([]quadtree.Body, 1000)
	spreadOnCircle( & bodies)
	
	var r Run
	r.Init( & bodies)
	
	var output *os.File
	output, _ = os.Create("essai.gif")
	
	r.outputGif( output, 100)
	// visual verification
}

func Test20Steps(t *testing.T) {
	
}

func TestGetModuloDistance(t *testing.T) {
	
	cases := []struct {
		x1, x2, want float64
	}{
		{0.0, 0.0, 0.0},
	}
	for _, c := range cases {
		got := getModuloDistance( c.x1, c.x2)
		if( got != c.want) {
			t.Errorf("x1 %f x2 %f == %f, want %f", c.x1, c.x2, got, c.want)
		}	
	}
}

func TestGetRepulsionVector(t *testing.T) {
	
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
func TestComputeRepulsiveForces(t *testing.T) {
	
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

func BenchmarkComputeRepulsiveForces_1K(b * testing.B ) {

	bodies := make([]quadtree.Body, 1000)
	spreadOnCircle( & bodies)
	var r Run
	r.Init( & bodies)
	for i := 0; i<b.N;i++ { 
		r.ComputeRepulsiveForce()
	}
}
func BenchmarkComputeRepulsiveForces_10K(b * testing.B ) {

	bodies := make([]quadtree.Body, 10000)
	spreadOnCircle( & bodies)
	var r Run
	r.Init( & bodies)
	for i := 0; i<b.N;i++ { 
		r.ComputeRepulsiveForce()
	}
}


func BenchmarkGetModuleDistance(b * testing.B ) {
	
	x := rand.Float64()
	y := rand.Float64()
		
	b.ResetTimer()
		
	for i := 0; i<b.N;i++ { 
		getModuloDistance( x, y)
	}
}

func BenchmarkGetRepulsionVector(b * testing.B ) {
	
	bodies := make([]quadtree.Body, 2)
	bodies[1].X = 0.5
	bodies[1].Y = 0.5
		
	for i := 0; i<b.N;i++ { 
		getRepulsionVector( &(bodies[0]), &(bodies[1]))
	}
}

// benchmark init
func BenchmarkInitRun1000000(b * testing.B) {
	
	bodies := make([]quadtree.Body, 1000000)

	if false { fmt.Printf("\n%#v", bodies[0]) }
	
	spreadOnCircle( & bodies)
	
	var r Run
	for i := 0; i<b.N;i++ {
		r.Init( & bodies)
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
