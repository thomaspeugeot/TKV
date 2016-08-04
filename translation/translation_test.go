package translation

import (
	"testing"
	"math"
)

var epsilon float64 = 0.0000001

// testing the translation function of country
func TestFRALatLng2XY(t *testing.T) {

	var fra Country
	fra.Name = "fra"
	fra.NbBodies = 222317
	fra.Step = 0

	fra.Init()

	cases := []struct {
		lat, lng, x, y float64
	}{
		{40.0, -6.0, 0.0, 0.0}, // south west corner
		{40.0+12.0, -6.0+17.0, 1.0, 1.0}, // south west corner
	}
	for _, c := range cases {
		gotX, gotY := fra.LatLng2XY(c.lat, c.lng)
		if math.Abs(gotX - c.x) > epsilon || math.Abs(gotY - c.y) > epsilon  {
			t.Errorf("lat %f lng %f, want x %f y %f, got x %f y %f", c.lat, c.lng, gotX, gotY, c.x, c.y)
		}	
	}

}
