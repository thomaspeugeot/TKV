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
// go run grump-reader.go -tkvdata="C:\Users\peugeot\tkv-data" -step=
func main() {

	// flag "country"
	countryPtr := flag.String("country","fra","iso 3166 country code")

	// flag "step"
	stepPtr := flag.String("step","0000","simulation step for the displaced bodies")

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
	
	mux.HandleFunc("/area", area)
	
	
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
