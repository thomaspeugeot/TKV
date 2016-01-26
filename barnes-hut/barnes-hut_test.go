package barnes_hut

import (
	"testing"
	"fmt"
)

// init 

// benchmark init
func BenchmarkInitRun1000000(b * testing.B) {
	
	bodies := make([]Body, 1000000)

	fmt.Printf("\n%#v", bodies[0])
	
	var r Run
	for i := 0; i<b.N;i++ {
		r.Init( & bodies)
	}
}
