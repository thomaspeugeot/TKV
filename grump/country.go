/*
Package grump contains function to process data of a GRUMP file describing the country density
into a body file that can be fed into the simulator.

GRUMP is the acronym for "Global Rural-Urban Mapping Project"; this is a freely available data source
of population density.
*/
package grump

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// GrumpSpacing is side length in degrees of the input data unit square
const GrumpSpacing float64 = 0.0083333333333

// Country stores country meta data
type Country struct {
	Name                               string
	NCols, NRows int
	XllCorner, YllCorner float64
}

// Row2Lat converts from row index to lat
func (country *Country) Row2Lat(row int) (lat float64) {
	// lat := float64( country.YllCorner) + (float64( country.NRows - row)*rowLatWidth)
	lat = float64(country.YllCorner) + float64(row)*GrumpSpacing
	return lat
}

// Serialize into a coord file
func (country *Country) Serialize() {

	filename := fmt.Sprintf("conf-%s.coord", country.Name)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return
	}
	jsonCountry, _ := json.MarshalIndent(country, "", "\t")
	file.Write(jsonCountry)
	file.Close()
}

// Unserialize inits struct from a coord file
func (country *Country) Unserialize() {

	filename := fmt.Sprintf("conf-%s.coord", country.Name)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(country); err != nil {
		log.Fatal(fmt.Sprintf("parsing config file %s", err.Error()))
	}

	Info.Printf("(Grump) Unserialize country %s", country.Name)

	Info.Printf("(Grump) Init Country orig lat %f lng %f size lat %f lng %f ",
		float64(country.YllCorner),
		float64(country.XllCorner),
		float64(country.NRows)*GrumpSpacing,
		float64(country.NCols)*GrumpSpacing)

	file.Close()
}

// LatLng2XY gives from lat/lng, the relative coordinate within the country
func (country *Country) LatLng2XY(lat, lng float64) (x, y float64) {

	// compute relative coordinates within the square
	x = (lng - float64(country.XllCorner)) / (float64(country.NCols) * GrumpSpacing)
	y = (lat - float64(country.YllCorner)) / (float64(country.NRows) * GrumpSpacing) // y is 0 at northest point and 1.0 at southest point

	return x, y
}

// XY2LatLng gives from lat/lng, the relative coordinate within the country
func (country *Country) XY2LatLng(x, y float64) (lat, lng float64) {

	lat = float64(country.YllCorner) + (y * float64(country.NRows) * GrumpSpacing)
	lng = float64(country.XllCorner) + (x * float64(country.NCols) * GrumpSpacing)

	return lat, lng
}
