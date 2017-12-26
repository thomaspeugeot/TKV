package barnes_hut

import (
	"github.com/thomaspeugeot/tkv/quadtree"
	"testing"
	"fmt"
	"os"
	"math/rand"
)


// MacBook-Pro-de-Thomas:barnes-hut thomaspeugeot$ go test -bench=.
// PASS
// BenchmarkComputeRepulsiveForces_1K-8             	      50	  23658474 ns/op
// BenchmarkComputeRepulsiveForces_10K-8            	       5	 292513993 ns/op
// BenchmarkComputeRepulsiveForcesOnHalfSet_1K-8    	     100	  11466750 ns/op
// BenchmarkComputeRepulsiveForcesConcurrent20_30K-8	       5	 253429700 ns/op
// BenchmarkGetModuleDistance-8                     	2000000000	         1.23 ns/op
// BenchmarkGetRepulsionVector-8                    	100000000	        15.5 ns/op
// BenchmarkInitRun_1M-8                            	       5	 304567105 ns/op
// BenchmarkOutputGif_1MBody_1KSteps-8              	       1	11655204392 ns/op

func BenchmarkComputeRepulsiveForces_1K(b * testing.B ) {

	bodies := make([]quadtree.Body, 1000)
	SpreadOnCircle( & bodies)
	var r Run
	r.Init( & bodies)
	for i := 0; i<b.N;i++ { 
		r.ComputeRepulsiveForce()
	}
}
func BenchmarkComputeRepulsiveForces_10K(b * testing.B ) {

	bodies := make([]quadtree.Body, 10000)
	SpreadOnCircle( & bodies)
	var r Run
	r.Init( & bodies)
	for i := 0; i<b.N;i++ { 
		r.ComputeRepulsiveForce()
	}
}

func BenchmarkComputeRepulsiveForcesOnHalfSet_1K(b * testing.B ) {

	bodies := make([]quadtree.Body, 1000)
	SpreadOnCircle( & bodies)
	var r Run
	r.Init( & bodies)
	endIndex := len(bodies)/2
	for i := 0; i<b.N;i++ { 
		r.ComputeRepulsiveForceSubSet(0, endIndex)
	}
}

func BenchmarkComputeRepulsiveForcesConcurrent20_30K(b * testing.B ) {

	bodies := make([]quadtree.Body, 30000)
	SpreadOnCircle( & bodies)
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

func BenchmarkGetVectorBetweenBodiesWithMirror(b * testing.B ) {
	
	bodies := make([]quadtree.Body, 2)
	bodies[1].X = 0.5
	bodies[1].Y = 0.5
		
	for i := 0; i<b.N;i++ { 
		getVectorBetweenBodiesWithMirror( &(bodies[0]), &(bodies[1]), 0, 0)
	}
}

func BenchmarkGetDistanceBetweenBodiesWithMirror(b * testing.B ) {
	
	bodies := make([]quadtree.Body, 2)
	bodies[1].X = 0.5
	bodies[1].Y = 0.5
		
	for i := 0; i<b.N;i++ { 
		getDistanceBetweenBodiesWithMirror( &(bodies[0]), &(bodies[1]), 0, 0)
	}
}

func BenchmarkGetRepulsionVector(b * testing.B ) {
	
	bodies := make([]quadtree.Body, 2)
	bodies[1].X = 0.5
	bodies[1].Y = 0.5
		
	for i := 0; i<b.N;i++ { 
		getRepulsionVector( &(bodies[0]), &(bodies[1]), 0, 0)
	}
}

// benchmark init
func BenchmarkInitRun_1M(b * testing.B) {
	
	bodies := make([]quadtree.Body, 1000 * 1000)

	if false { fmt.Printf("\n%#v", bodies[0]) }
	
	SpreadOnCircle( & bodies)
	
	var r Run
	for i := 0; i<b.N;i++ {
		r.Init( & bodies)
	}
}

// benchmark gif output
func BenchmarkOutputGif_1MBody_1KSteps(b * testing.B) {

	var bodies []quadtree.Body
	quadtree.InitBodiesUniform( &bodies, 2000)

	SpreadOnCircle( & bodies)
	
	var r Run
	r.Init( & bodies)
	
	var output *os.File
	output, _ = os.Create("essai30Kbody_6Ksteps.gif")
	
	for i := 0; i<b.N;i++ {
		r.SetState( RUNNING)
		r.OutputGif( output, 600)
	}
}
