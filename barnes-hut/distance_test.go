// distance computation functions
package barnes_hut

import (
	"github.com/thomaspeugeot/tkv/quadtree"
	"testing"
)

// test gif output
func TestMirror(t *testing.T) {

	cases := make( []struct {
		A, B quadtree.Body
		x, y int
		wantX, wantY float64
	}, 3)
	
	bodies := make([]quadtree.Body, 2)
	bodies[1].X = 0.4
	bodies[1].Y = 0.3
	
	cases[0].B = bodies[1]
	cases[0].x = 0
	cases[0].y = 0
	cases[0].wantX = 0.4
	cases[0].wantY = 0.3
	
	cases[1].B = bodies[1]
	cases[1].x = -1
	cases[1].y = 0
	cases[1].wantX = -0.4
	cases[1].wantY = 0.3
	
	cases[2].B = bodies[1]
	cases[2].x = -1
	cases[2].y = -1
	cases[2].wantX = -0.4
	cases[2].wantY = -0.3
	
	for _, c := range cases {
		gotX, gotY := getVectorBetweenBodiesWithMirror( & c.A, & c.B, c.x, c.y)
		if( (gotX != c.wantX) && (gotY != c.wantY)) {
			t.Errorf("vect mirror x %d y %d A x %f y %f B %f %f == %f %f, want %f %f", c.x, c.y, c.A.X, c.A.Y, c.B.X, c.B.Y, gotX, gotY, c.wantX, c.wantY )
		}	
	}
}