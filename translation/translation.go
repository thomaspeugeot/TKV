// this package provides the function for retriving villages locations, borders as well as 
// translation 
package translation

import (
)

type Translation struct {
	xMin, xMax, yMin, yMax float64 // coordinates of the rendering window
}


func (t * Translation) Init(country Country) {

	country.Init()

	Info.Printf("Country is %s with step %d", country.Name, country.Step)

	// rowLatWidth := 0.0083333333333
	// colLngWidth := 0.0083333333333

	// load final config
	// load initial config
}

func (t * Translation) SetRenderingWindow( xMin, xMax, yMin, yMax float64) {
	t.xMin, t.xMax, t.yMin, t.yMax = xMin, xMax, yMin, yMax
}



