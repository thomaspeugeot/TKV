package barneshut

import (
	"io/ioutil"
	"testing"
)

func OldTestRepulsionFieldInit(t *testing.T) {

	r := NewRun()
	r.LoadConfig("conf-fra-00000.bods")

	// get pointer on quadtree
	q := &(r.q)
	Info.Printf("TestRepulsionFieldInit pointer on quadtree %p", q)
	r.gridFieldNb = 4

	f := NewRepulsionField(0.3, 0.5,
		0.4, 0.6,
		r.gridFieldNb,
		q, // quadtree
		0.00001)
	f.ComputeField()
	r.fieldRendering = true
	Info.Printf("TestRepulsionFieldInit value at 1 1 %e", f.values[1][1])

	r.RenderGif(ioutil.Discard, false)

	cases := make([]struct {
		i, j         int
		wantX, wantY float64
	}, 1)

	cases[0].i = 1
	cases[0].j = 2
	cases[0].wantX = 0.3375 // 0.3 + 0.125 * (1 + 2*1)
	cases[0].wantY = 0.5625 // 0.5 + 0.125 * (1 + 2*2)

	for _, c := range cases {
		gotX, gotY := f.XY(c.i, c.j)
		if (gotX != c.wantX) && (gotY != c.wantY) {
			t.Errorf("i %d j %d == %f %f, want %f %f", c.i, c.j, gotX, gotY, c.wantX, c.wantY)
		}
	}
}
