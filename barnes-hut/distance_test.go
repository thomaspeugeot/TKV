package barneshut

import (
	"math"
	"testing"

	"github.com/thomaspeugeot/tkv/quadtree"
)

// test gif output
func TestMirror(t *testing.T) {

	cases := make([]struct {
		A, B         quadtree.Body
		x, y         int
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
		gotX, gotY := getVectorBetweenBodiesWithMirror(&c.A, &c.B, c.x, c.y)
		if (gotX != c.wantX) && (gotY != c.wantY) {
			t.Errorf("vect mirror x %d y %d A x %f y %f B %f %f == %f %f, want %f %f", c.x, c.y, c.A.X, c.A.Y, c.B.X, c.B.Y, gotX, gotY, c.wantX, c.wantY)
		}
	}
}

func TestMirrorDistance(t *testing.T) {

	cases := make([]struct {
		A, B   quadtree.Body
		xM, yM int
		wantD  float64
	}, 3)

	bodies := make([]quadtree.Body, 5)
	bodies[1].X = 0.4
	bodies[1].Y = 0.3

	bodies[2].X = 0
	bodies[2].Y = 0

	bodies[3].X = 0
	bodies[3].Y = 1.0

	bodies[4].X = 1.0
	bodies[4].Y = 1.0

	cases[0].B = bodies[1]
	cases[0].xM = 0
	cases[0].yM = 0
	cases[0].wantD = 0.5

	cases[1].A = bodies[2]
	cases[1].B = bodies[3]
	cases[1].xM = 0
	cases[1].yM = 0
	cases[1].wantD = 1.0

	cases[2].A = bodies[2]
	cases[2].B = bodies[4]
	cases[2].xM = 0
	cases[2].yM = 0
	cases[2].wantD = math.Sqrt(2.0)

	for _, c := range cases {
		gotD := getDistanceBetweenBodiesWithMirror(&c.A, &c.B, c.xM, c.yM)
		if gotD != c.wantD {
			t.Errorf("vect mirror x %d y %d A x %f y %f B %f %f == %f, want %f", c.xM, c.yM, c.A.X, c.A.Y, c.B.X, c.B.Y, gotD, c.wantD)
		}
	}
}
