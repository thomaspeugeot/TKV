package grump

import (
	"fmt"
	"os"
	"log"
	"encoding/json"
	)



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

	file.Close()
}