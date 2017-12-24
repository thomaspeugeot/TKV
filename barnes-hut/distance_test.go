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
	}, 1)
	
	bodies := make([]quadtree.Body, 2)
	bodies[1].X = 0.4
	bodies[1].Y = 0.3
	
	cases[0].B = bodies[1]
	cases[0].x = 0
	cases[0].wantX = 0.4
	cases[0].wantY = 0.3
	
	for _, c := range cases {
		gotX, gotY := getVectorBetweenBodiesWithMirror( & c.A, & c.B, 0, 0)
		if( (gotX != c.wantX) && (gotY != c.wantY)) {
			t.Errorf("A %#v B %#v == %f %f, want %f %f", c.A, c.B, gotX, gotY, c.wantX, c.wantY )
		}	
	}
}