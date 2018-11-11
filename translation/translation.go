/*
Package translation provides functions for managing a translation between two countries.
*/
package translation

// Singloton pointing to the current translation
// the singloton can be autocally initiated if it is nil
var translateCurrent Translation

// Singloton pattern to init the current translation
func GetTranslateCurrent() *Translation {

	// check if the current translation is void.
	if translateCurrent.sourceCountry == nil {
		var sourceCountry Country
		var targetCountry Country

		sourceCountry.Name = "fra"
		sourceCountry.NbBodies = 934136
		sourceCountry.Step = 8725

		targetCountry.Name = "hti"
		targetCountry.NbBodies = 927787
		targetCountry.Step = 8564

		translateCurrent.Init(sourceCountry, targetCountry)
	}

	return &translateCurrent
}

// Definition of a translation between a source and a target country
type Translation struct {
	xMin, xMax, yMin, yMax float64 // coordinates of the rendering window (used to compute liste of villages)
	sourceCountry          *Country
	targetCountry          *Country
}

func (t *Translation) GetSourceCountryName() string {
	return t.sourceCountry.Name
}

func (t *Translation) GetTargetCountryName() string {
	return t.targetCountry.Name
}

// Init source & target countries of the translation
func (t *Translation) Init(sourceCountry, targetCountry Country) {

	Info.Printf("Init : Source Country is %s with nbBodies %d at simulation step %d", sourceCountry.Name, sourceCountry.NbBodies, sourceCountry.Step)
	Info.Printf("Init : Target Country is %s with nbBodies %d at simulation step %d", targetCountry.Name, targetCountry.NbBodies, targetCountry.Step)

	t.sourceCountry = &sourceCountry
	t.sourceCountry.Init()

	t.targetCountry = &targetCountry
	t.targetCountry.Init()
}

// Swap source & target
func (t *Translation) Swap() {
	tmp := t.sourceCountry
	t.sourceCountry = t.targetCountry
	t.targetCountry = tmp
}

func (t *Translation) SetRenderingWindow(xMin, xMax, yMin, yMax float64) {
	t.xMin, t.xMax, t.yMin, t.yMax = xMin, xMax, yMin, yMax
}

// from lat, lng in source country, find the closest body in source country
func (t *Translation) BodyCoordsInSourceCountry(lat, lng float64) (distance, latClosest, lngClosest, xSpread, ySpread float64, closestIndex int) {

	// convert from lat lng to x, y in the Country
	return t.sourceCountry.ClosestBodyInOriginalPosition(lat, lng)
}

// from lat, lng in source country, find the closest body in source country
func (t *Translation) BodyCoordsInTargetCountry(lat, lng float64) (distance, latClosest, lngClosest, xSpread, ySpread float64, closestIndex int) {

	// convert from lat lng to x, y in the Country
	return t.targetCountry.ClosestBodyInOriginalPosition(lat, lng)
}

// from x, y get closest body lat/lng in target country
func (t *Translation) LatLngToXYInTargetCountry(x, y float64) (latTarget, lngTarget float64) {

	return t.targetCountry.XYToLatLng(x, y)
}

// from a coordinate in source coutry, get border
func (t *Translation) TargetBorder(x, y float64) PointList {

	return t.targetCountry.XYtoTerritoryBodies(x, y)
}

func (t *Translation) SourceBorder(lat, lng float64) PointList {

	Info.Printf("Source Border for lat %f lng %f", lat, lng)

	points := t.sourceCountry.LatLngToTerritoryBorder(lat, lng)

	Info.Printf("Source Border nb of points %d", len(points))

	return points
}
