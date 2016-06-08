package barnes_hut

import (
	"testing"
)

func TestRepulsionFieldInit(t *testing.T) {

	r := NewRun()
	r.LoadConfig("conf-fra-00000.bods")

	// get pointer on quadtree
	q := & (r.q)
	Info.Printf( "TestRepulsionFieldInit pointer on quadtree %p", q)

	f := NewRepulsionField( 0.3, 0.5, 
							0.4, 0.6, 
							4,
							q) // quadtree
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
