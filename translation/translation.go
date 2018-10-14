// this package provides the function for retriving villages locations, borders as well as
// translation
package translation

// singloton pointing to the current translation
// the singloton can be autocally initiated if it is nil
var translateCurrent Translation

// singloton pattern to init the current translation
func GetTranslateCurrent() *Translation {

	// check if the current translation is void.
	if translateCurrent.GetSourceCountryName() != "fra" {
		var sourceCountry Country
		var targetCountry Country

		sourceCountry.Name = "fra"
		sourceCountry.NbBodies = 697529
		sourceCountry.Step = 4723

		targetCountry.Name = "hti"
		targetCountry.NbBodies = 927787
		targetCountry.Step = 8564

		translateCurrent.Init(sourceCountry, targetCountry)
	}

	return &translateCurrent
}

type Translation struct {
	xMin, xMax, yMin, yMax float64 // coordinates of the rendering window (used to compute liste of villages)
	sourceCountry          Country
	targetCountry          Country
}

func (t *Translation) GetSourceCountryName() string {
	return t.sourceCountry.Name
}

func (t *Translation) Init(sourceCountry, targetCountry Country) {

	Info.Printf("Init : Source Country is %s with nbBodies %d at simulation step %d", sourceCountry.Name, sourceCountry.NbBodies, sourceCountry.Step)
	Info.Printf("Init : Target Country is %s with nbBodies %d at simulation step %d", targetCountry.Name, targetCountry.NbBodies, targetCountry.Step)

	t.sourceCountry = sourceCountry
	t.sourceCountry.Init()

	t.targetCountry = targetCountry
	t.targetCountry.Init()

}

func (t *Translation) SetRenderingWindow(xMin, xMax, yMin, yMax float64) {
	t.xMin, t.xMax, t.yMin, t.yMax = xMin, xMax, yMin, yMax
}

// from lat, lng in source country, find the closest body in source country
func (t *Translation) ClosestBodyInOriginalPosition(lat, lng float64) (x, y, distance, latClosest, lngClosest, xSpread, ySpread float64, closestIndex int) {

	// convert from lat lng to x, y in the Country
	return t.sourceCountry.ClosestBodyInOriginalPosition(lat, lng)
}

// from x, y corrdinates in spread, get closest body lat/lng in target country
func (t *Translation) XYSpreadToLatLngInTargetCountry(xSpread, ySpread float64) (latTarget, lngTarget float64) {

	Info.Printf("XYSpreadToLatLngInTargetCountry input xSpread %f ySpread %f", xSpread, ySpread)

	latTarget, lngTarget = t.targetCountry.XYSpreadToLatLngOrig(xSpread, ySpread)

	Info.Printf("XYSpreadToLatLngInTargetCountry output lat %f lng %f", latTarget, lngTarget)

	return latTarget, lngTarget
}

// from a coordinate in source coutry, get border
func (t *Translation) TargetBorder(xSpread, ySpread float64) PointList {

	Info.Printf("TargetBorder input xSpread %f ySpread %f", xSpread, ySpread)

	points := t.targetCountry.XYSpreadToTerritoryBorder(xSpread, ySpread)

	Info.Printf("Target Border nb of points %d", len(points))

	return points
}

func (t *Translation) SourceBorder(lat, lng float64) PointList {

	Info.Printf("Source Border for lat %f lng %f", lat, lng)

	points := t.sourceCountry.LatLngToTerritoryBorder(lat, lng)

	Info.Printf("Source Border nb of points %d", len(points))

	return points
}
