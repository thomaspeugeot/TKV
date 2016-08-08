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

// test that brest lat long has close proximity to nearest village
func TestBrestVillageProximity(t * testing.T) {

	var fra Country
	fra.Name = "fra"
	fra.NbBodies = 222317
	fra.Step = 0

	fra.Init()

	//	48° 23′ 27″ Nord 4° 29′ 08″ Ouest
	lat := 48.0 + 23.0*1.0/60.0
	lng := -4.0 - 29.0*1.0/60.0

	// fra.LatLng2XY( lat, lng)

	x, y, distance, latClosest, lngClosest := fra.VillageCoordinates( lat, lng)

	deltaLat := math.Abs(latClosest-lat)
	if deltaLat > 0.1 { // we tolerate one 10th of a degree
		t.Errorf("Latitude of closest village too far, origin lat %f, village lat %f, delta lat %f, x %f, y %f, distance %d", lat, latClosest, deltaLat, x, y, distance) 
	}

	deltaLng := math.Abs(lngClosest-lng)
	if( deltaLng > 0.1) { // we tolerate one 10th of a degree
		t.Errorf("Longitude of closest village too far, origin lng %f, village lng %f, delta lng %f, x %f, y %f, distance %d", lng, lngClosest, deltaLng, x, y, distance) 
	} 
}
// test that middle of the atlantic has far distance to nearest village 
