package barnes_hut

import (
	"github.com/thomaspeugeot/tkv/quadtree"
	"testing"
	"fmt"
)

// init 

// benchmark init
func BenchmarkInitRun1000000(b * testing.B) {
	
	bodies := make([]quadtree.Body, 1000000)

	fmt.Printf("\n%#v", bodies[0])
	
	var r Run
	for i := 0; i<b.N;i++ {
		r.Init( & bodies)
	}
}
