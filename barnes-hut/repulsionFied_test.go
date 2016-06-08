package barnes_hut

import (
	"github.com/thomaspeugeot/tkv/quadtree"
	"testing"
)

func TestRepulsionFieldInit(t *testing.T) {

	bodies := make([]quadtree.Body, 10 * 10)
	SpreadOnCircle( & bodies)
	var r Run
	r.Init( & bodies)

	f := NewRepulsionField( 0.3, 0.5, 
							0.4, 0.6, 
							4,
							& (r.q)) // quadtree
	f.ComputeField()

	cases := make( []struct {
		i, j int
		wantX, wantY float64
	}, 1)
	
	cases[0].i = 1
	cases[0].j = 2
	cases[0].wantX = 0.3375 // 0.3 + 0.125 * (1 + 2*1)
	cases[0].wantY = 0.5625 // 0.5 + 0.125 * (1 + 2*2)

	for _, c := range cases {
		gotX, gotY := f.XY( c.i, c.j)
		if( (gotX != c.wantX) && (gotY != c.wantY)) {
			t.Errorf("i %d j %d == %f %f, want %f %f", c.i, c.j, gotX, gotY, c.wantX, c.wantY )
		}	
	}
}
