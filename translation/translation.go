/*
Package translation provides functions for managing a translation between two countries.
*/
package translation

/* import (
	"github.com/thomaspeugeot/tkv/grump"
) */

// Singloton pointing to the current translation
// the singloton can be autocally initiated if it is nil
var translateCurrent Translation

// storage for all countries
var mapOfCountries map[string]*CountryWithBodies = nil

type CountrySpec struct {
	Name           string
	NbBodies, Step int
}

// should have done this array directly with CountryWithBodies but
// aint compile even with
// 	CountryWithBodies{Name: "fra", NbBodies: 934136, Step: 8725},
var countrySpecs = []CountrySpec{
	CountrySpec{Name: "fra", NbBodies: 934136, Step: 8725},
	CountrySpec{Name: "hti", NbBodies: 190948, Step: 1334},
	CountrySpec{Name: "usa", NbBodies: 1422837, Step: 2735},
	// CountrySpec{Name: "chn", NbBodies: 771973, Step: 2531},
	// CountrySpec{Name: "rus", NbBodies: 509497, Step: 3386},
}

// Singloton pattern to init the current translation
func GetTranslateCurrent() *Translation {

	// check if the current translation is void.
	if mapOfCountries == nil {

		mapOfCountries = make(map[string]*CountryWithBodies)

		for _, countrySpec := range countrySpecs {
			country := CountryWithBodies{}
			country.Name = countrySpec.Name
			country.NbBodies = countrySpec.NbBodies
			country.Step = countrySpec.Step
			country.Init()

			mapOfCountries[country.Name] = &country
		}

		translateCurrent.sourceCountry = mapOfCountries["fra"]
		translateCurrent.targetCountry = mapOfCountries["hti"]
	}

	return &translateCurrent
}

// Definition of a translation between a source and a target country
type Translation struct {
	sourceCountry *CountryWithBodies
	targetCountry *CountryWithBodies
}

func (t *Translation) GetSourceCountryName() string {
	return t.sourceCountry.Name
}

func (t *Translation) SetSourceCountry(name string) {
	t.sourceCountry = mapOfCountries[name]
}

func (t *Translation) GetTargetCountryName() string {
	return t.targetCountry.Name
}

func (t *Translation) SetTargetCountry(name string) {
	t.targetCountry = mapOfCountries[name]
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
