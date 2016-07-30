package main

import (
	"github.com/thomaspeugeot/tkv/server"
	"github.com/thomaspeugeot/tkv/translation"
	"net/http"
	"log"
	"flag"
	"fmt"
	"encoding/json"
)


// load the data
var t translation.Translation

// for decoding the rendering window
type test_struct struct {
	X1, X2, Y1, Y2 float64
}
//
// go run grump-reader.go -tkvdata="C:\Users\peugeot\tkv-data" -nbBodies=222317 -step=8542
func main() {

	// flag "country"
	countryPtr := flag.String("country","fra","iso 3166 country code")

	// flag "nbBodies"
	nbBodiesPtr := flag.String("nbBodies","222317","nb of bodies")

	// flag "step"
	stepPtr := flag.String("step","8542","simulation step for the spread bodies")

	flag.Parse()

	// init country from flags
	var country translation.Country
	country.Name = *countryPtr
	{
		_, errScan := fmt.Sscanf(*stepPtr, "%d", & country.Step)
		if( errScan != nil) {
			log.Fatal(errScan)
			return			
		}
	}
	{
		_, errScan := fmt.Sscanf(*nbBodiesPtr, "%d", & country.NbBodies)
		if( errScan != nil) {
			log.Fatal(errScan)
			return			
		}
	}


	server.Info.Printf("country to parse %s", country.Name)
	server.Info.Printf("nbBodies to parse %d", country.NbBodies)
	server.Info.Printf("step to parse %d", country.Step)

	t.Init(country)

	port := "localhost:8001"

	server.Info.Printf("begin listen on port %s", port)
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("../tkv-client/")) )
	
	mux.HandleFunc("/area", area)
	mux.HandleFunc("/villageCoordinates", villageCoordinates)
	
	
	log.Fatal(http.ListenAndServe(port, mux))
	server.Info.Printf("end")
}

func area(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder( req.Body)
	var renderingWindow test_struct
	err := decoder.Decode( &renderingWindow)
	if err != nil {
		log.Println("error decoding ", err)
	}
	t.SetRenderingWindow( renderingWindow.X1, renderingWindow.X2, renderingWindow.Y1, renderingWindow.Y2)
}

type latLng struct {
	lat, lng float64
}

// get village coordinates from lat/long
func villageCoordinates(w http.ResponseWriter, req *http.Request) {
	
	// parse lat long from client
	decoder := json.NewDecoder( req.Body)
	var t latLng
	err := decoder.Decode( &t)
	if err != nil {
		log.Println("error decoding ", err)
	}
	server.Info.Printf("villageCoordinates for lat %f, lng %d", t.lat, t.lng)

	// dircontent, _ := json.MarshalIndent( r.DirConfig(), "", "	")
	// fmt.Fprintf(w, "%s", dircontent)
}
