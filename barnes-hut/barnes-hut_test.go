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

	bodies := make([]quadtree.Body, 100000)
	spreadOnCircle( & bodies)
	
	var r Run
	r.Init( & bodies)
	
	var output *os.File
	output, _ = os.Create("essai.gif")
	
	r.outputGif( output)
	// visual verification
}

func Test20Steps(t *testing.T) {
	
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

func spreadOnCircle(bodies * []quadtree.Body) {
	for idx, _ := range *bodies {
		
		body := &((*bodies)[idx])
		
		radius := 0.6 * rand.Float64()
		angle := 2.0 * math.Pi * rand.Float64()
		
		if idx%2 == 0 {
			body.X += 0.4
			body.Y += 0.7
			radius *= 0.5
		} else {
			body.X -= 0.2
			body.Y -= 0.1
			radius *= 1.3
		}
		
		body.M =0.1000000
		body.X += math.Cos( angle) * radius
		body.Y += math.Sin( angle) * radius
	}
}
