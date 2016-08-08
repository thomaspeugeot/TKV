package grump

import (
	"fmt"
	"os"
	"log"
	"encoding/json"
	)

const GrumpSpacing float64 = 0.0083333333333

// store country code
type Country struct {
	Name string
	NCols, NRows, XllCorner, YllCorner int
}	

func (country * Country) Serialize() {

	filename := fmt.Sprintf("conf-%s.coord", country.Name)
	file, err := os.Create(filename)
	if( err != nil) {
		log.Fatal(err)
		return
	}
	jsonCountry, _ := json.MarshalIndent( country, "","\t")
	file.Write( jsonCountry)
	file.Close()
}

func (country * Country) Unserialize() {

	filename := fmt.Sprintf("conf-%s.coord", country.Name)
	file, err := os.Open(filename)
	if( err != nil) {
		log.Fatal(err)
		return
	}

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(country); err != nil {
		log.Fatal( fmt.Sprintf( "parsing config file %s", err.Error()))
	}

	Info.Printf( "(Grump) Unserialize country %s", country.Name)

	Info.Printf("Init Country size lng %f lat %f", 
		float64(country.NCols) * GrumpSpacing,
		float64(country.NRows) * GrumpSpacing)

	file.Close()
}

// given a lat/lng, provides the relative coordinate within the country
func (country * Country) LatLng2XY( lat, lng float64) (x, y float64) {

	// compute relative coordinates within the square
	x = (lng - float64( country.XllCorner)) / (float64(country.NCols) * GrumpSpacing)
	y = (lat - float64( country.YllCorner)) / (float64(country.NRows) * GrumpSpacing) // y is 0 at northest point and 1.0 at southest point

	return x, y
}

// given a lat/lng, provides the relative coordinate within the country
func (country * Country) XY2LatLng(x, y float64) ( lat, lng float64) {

	lat = float64( country.YllCorner) + (y * float64(country.NCols) * GrumpSpacing)
	lng = float64( country.XllCorner) + (x * float64(country.NRows) * GrumpSpacing)

	return lat, lng
}



