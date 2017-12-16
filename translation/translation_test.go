package translation

import (
	"testing"
	"math"
	"fmt"
)

var epsilon float64 = 0.0000001

// testing the translation lat lng to XY function of country
func TestFRALatLng2XY(t *testing.T) {

	var fra Country
	fra.Name = "fra"
	fra.NbBodies = 34413
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

// test the translation from XY to lat lng of counry
func TestFRAXY2LatLng(t * testing.T) {


	var fra Country
	fra.Name = "fra"
	fra.NbBodies = 34413
	fra.Step = 0

	fra.Init()

	cases := []struct {
		lat, lng, x, y float64
	}{
		{40.0, -6.0, 0.0, 0.0}, // south west corner
		{40.0+12.0, -6.0+17.0, 1.0, 1.0}, // south west corner
	}
	for _, c := range cases {
		gotLat, gotLng := fra.XY2LatLng(c.x, c.y)
		if math.Abs(gotLat - c.lat) > epsilon || math.Abs(gotLng - c.lng) > epsilon  {
			t.Errorf("x %f y %f, want lat %f lng %f, got lat %f lng %f", c.x, c.y, c.lat, c.lng, gotLat, gotLng)
		}	
	}

}

// test that brest lat long has close proximity to nearest village
func TestBrestVillageProximity(t * testing.T) {

	var fra Country
	fra.Name = "fra"
	fra.NbBodies = 34413
	fra.Step = 0

	fra.Init()

	//	48° 23′ 27″ Nord 4° 29′ 08″ Ouest
	lat := 48.0 + 23.0*1.0/60.0
	lng := -4.0 - 29.0*1.0/60.0

	// fra.LatLng2XY( lat, lng)

	x, y, distance, latClosest, lngClosest, _, _, _ := fra.VillageCoordinates( lat, lng)

	deltaLat := math.Abs(latClosest-lat)
	if deltaLat > 0.1 { // we tolerate one 10th of a degree
		t.Errorf("Latitude of closest village too far, origin lat %f, village lat %f, delta lat %f, x %d, y %d, distance %f", lat, latClosest, deltaLat, x, y, distance) 
	}

	deltaLng := math.Abs(lngClosest-lng)
	if( deltaLng > 0.1) { // we tolerate one 10th of a degree
		t.Errorf("Longitude of closest village too far, origin lng %f, village lng %f, delta lng %f, x %d, y %d, distance %f", lng, lngClosest, deltaLng, x, y, distance) 
	} 
}
// test that middle of the atlantic has far distance to nearest village 

// test that the nb of bodies of each village match the total number of bodies for a country
// test that brest lat long has close proximity to nearest village
func TestBallBodiesCount(t * testing.T) {

	var fra Country
	fra.Name = "fra"
	fra.NbBodies = 154301
	fra.Step = 96962

	fra.Init()

	var totalBodies int

	for _,vRow := range fra.villages {
		for _,v := range vRow {
			totalBodies+=v.NbBodies
			fmt.Printf("%d;", v.NbBodies)
		}
		fmt.Println()
	}

	if totalBodies != fra.NbBodies {
		t.Errorf("total bodies %d not matching nb bodies of country %d", totalBodies, fra.NbBodies)
	}

}

