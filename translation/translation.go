// this package provides the function for retriving villages locations, borders as well as 
// translation 
package translation

import (
)


type Translation struct {
	xMin, xMax, yMin, yMax float64 // coordinates of the rendering window (used to compute liste of villages)
	sourceCountry Country
	targetCountry Country
}


func (t * Translation) Init(sourceCountry, targetCountry Country) {

	Info.Printf("Init : Source Country is %s with step %d", sourceCountry.Name, sourceCountry.Step)
	Info.Printf("Init : Target Country is %s with step %d", sourceCountry.Name, sourceCountry.Step)

	t.sourceCountry = sourceCountry
	t.sourceCountry.Init()

	t.targetCountry = targetCountry
	t.targetCountry.Init()

}

func (t * Translation) SetRenderingWindow( xMin, xMax, yMin, yMax float64) {
	t.xMin, t.xMax, t.yMin, t.yMax = xMin, xMax, yMin, yMax
}


// 
func (t * Translation) VillageCoordinates( lat, lng float64) (x, y int, distance, latClosest, lngClosest float64) {

	// convert from lat lng to x, y in the Country 
	return t.sourceCountry.VillageCoordinates( lat, lng)
}








