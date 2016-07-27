package main

import (
	"github.com/thomaspeugeot/tkv/server"
	"github.com/thomaspeugeot/tkv/translation"
	"net/http"
	"log"
	"flag"
	"fmt"
)

//
// go run grump-reader.go -tkvdata="C:\Users\peugeot\tkv-data" -step=

func main() {

	// load the data
	var t translation.Translation

	// flag "country"
	countryPtr := flag.String("country","fra","iso 3166 country code")

	// flag "step"
	stepPtr := flag.String("step","0000","simulation step for the displaced bodies")

	// get the directory containing tkv data through the flag "tkvdata"
	// dirTKVDataPtr := flag.String("tkvdata","/Users/thomaspeugeot/the-mapping-data/","directory containing input tkv data")
	
	flag.Parse()

	var country translation.Country
	country.Name = *countryPtr
	_, errScan := fmt.Sscanf(*stepPtr, "%d", & country.Step)
	if( errScan != nil) {
		log.Fatal(errScan)
		return			
	}

	server.Info.Printf("country to parse %s", country.Name)
	server.Info.Printf("step to parse %d", country.Step)

	t.Init(country)

	port := "localhost:8001"

	server.Info.Printf("begin listen on port %s", port)
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("../leaflets/")) )
	
	log.Fatal(http.ListenAndServe(port, mux))
	server.Info.Printf("end")
}

