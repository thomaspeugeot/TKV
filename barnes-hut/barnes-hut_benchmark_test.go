package barnes_hut

import (
	"github.com/thomaspeugeot/tkv/quadtree"
	"testing"
	"fmt"
	"os"
	"math/rand"
)

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

func BenchmarkComputeRepulsiveForcesOnHalfSet_1K(b * testing.B ) {

	bodies := make([]quadtree.Body, 1000)
	spreadOnCircle( & bodies)
	var r Run
	r.Init( & bodies)
	endIndex := len(bodies)/2
	for i := 0; i<b.N;i++ { 
		r.ComputeRepulsiveForceSubSet(0, endIndex)
	}
}

func BenchmarkComputeRepulsiveForcesConcurrent20_30K(b * testing.B ) {

	bodies := make([]quadtree.Body, 30000)
	spreadOnCircle( & bodies)
	var r Run
	r.Init( & bodies)
	for i := 0; i<b.N;i++ { 
		r.ComputeRepulsiveForceConcurrent( 20)
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
func BenchmarkInitRun_1M(b * testing.B) {
	
	bodies := make([]quadtree.Body, 1000 * 1000)

	if false { fmt.Printf("\n%#v", bodies[0]) }
	
	spreadOnCircle( & bodies)
	
	var r Run
	for i := 0; i<b.N;i++ {
		r.Init( & bodies)
	}
}

// benchmark gif output
func BenchmarkOutputGif_10KBody_2KSteps(b * testing.B) {

	bodies := make([]quadtree.Body, 20000)
	spreadOnCircle( & bodies)
	
	var r Run
	r.Init( & bodies)
	
	var output *os.File
	output, _ = os.Create("essai1Kbody_1Ksteps.gif")
	
	for i := 0; i<b.N;i++ {
		r.outputGif( output, 50)
	}
}
