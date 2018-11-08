/*
Package main contains code for running the 10000 web server as a standalone server (no need for the cloud)
*/
package main

import (
	"log"
	"net/http"

	"github.com/thomaspeugeot/tkv/handler"
	"github.com/thomaspeugeot/tkv/server"
	"github.com/thomaspeugeot/tkv/translation"
)

// go run runtime_server.go -targetCountryStep=43439
func main() {

	t := translation.GetTranslateCurrent()
	server.Info.Printf("ended init of country %s", t.GetSourceCountryName())

	port := "localhost:8002"

	server.Info.Printf("begin listen on port %s", port)
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("../gae_tkv/")))

	mux.HandleFunc("/translateLatLngInSourceCountryToLatLngInTargetCountry",
		handler.GetTranslationResult)
	mux.HandleFunc("/villageTargetBorder", handler.VillageTargetBorder)
	mux.HandleFunc("/villageSourceBorder", handler.VillageSourceBorder)
	mux.HandleFunc("/allSourcPointsCoordinates", handler.AllSourceBorderPointsCoordinates)

	log.Fatal(http.ListenAndServe(port, mux))
	server.Info.Printf("end")

}
