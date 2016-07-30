// this package provides the function for retriving villages locations, borders as well as 
// translation 
package translation

import (
)

const GrumpSpacing float64 = 0.0083333333333


type Translation struct {
	xMin, xMax, yMin, yMax float64 // coordinates of the rendering window (used to compute liste of villages)
	country Country
}


func (t * Translation) Init(country Country) {

	t.country = country
	t.country.Init()

	Info.Printf("Country is %s with step %d", country.Name, country.Step)

	// rowLatWidth := 0.0083333333333
	// colLngWidth := 0.0083333333333

	// load final config
	// load initial config
}

func (t * Translation) SetRenderingWindow( xMin, xMax, yMin, yMax float64) {
	t.xMin, t.xMax, t.yMin, t.yMax = xMin, xMax, yMin, yMax
}


// 
func (t * Translation) VillageCoordinates( lat, lng float64) (x, y int) {

	// we work for france only 
	// convert from lat lng to x, y in the Country 
	return t.country.VillageCoordinates( lat, lng)
}


