/*
Package main contains code for running the 10000 web server as a standalone server (no need for the cloud)
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/thomaspeugeot/tkv/handler"
	"github.com/thomaspeugeot/tkv/server"
	"github.com/thomaspeugeot/tkv/translation"
)

// go run runtime_server.go -targetCountryStep=43439
func main() {

	// flags  for source country
	sourceCountryPtr := flag.String("sourceCountry", "fra", "iso 3166 sourceCountry code")
	sourceCountryNbBodiesPtr := flag.String("sourceCountryNbBodiesPtr", "697529", "nb of bodies")
	sourceCountryStepPtr := flag.String("sourceCountryStep", "4723", "simulation step for the spread bodies for source country")

	// flags  for target country
	targetCountryPtr := flag.String("targetCountry", "hti", "iso 3166 targetCountry code")
	targetCountryNbBodiesPtr := flag.String("targetCountryNbBodiesPtr", "927787", "nb of bodies for target country")
	targetCountryStepPtr := flag.String("targetCountryStep", "8564", "simulation step for the spread bodies for target country")

	flag.Parse()

	// init sourceCountry from flags
	var sourceCountry translation.Country
	sourceCountry.Name = *sourceCountryPtr
	{
		_, errScan := fmt.Sscanf(*sourceCountryNbBodiesPtr, "%d", &sourceCountry.NbBodies)
		if errScan != nil {
			log.Fatal(errScan)
			return
		}
	}
	{
		_, errScan := fmt.Sscanf(*sourceCountryStepPtr, "%d", &sourceCountry.Step)
		if errScan != nil {
			log.Fatal(errScan)
			return
		}
	}

	// init targetCountry from flags
	var targetCountry translation.Country
	targetCountry.Name = *targetCountryPtr
	{
		_, errScan := fmt.Sscanf(*targetCountryNbBodiesPtr, "%d", &targetCountry.NbBodies)
		if errScan != nil {
			log.Fatal(errScan)
			return
		}
	}
	{
		_, errScan := fmt.Sscanf(*targetCountryStepPtr, "%d", &targetCountry.Step)
		if errScan != nil {
			log.Fatal(errScan)
			return
		}
	}

	server.Info.Printf("sourceCountry to parse %s", sourceCountry.Name)
	server.Info.Printf("sourceCountry nbBodies in file to parse %d", sourceCountry.NbBodies)
	server.Info.Printf("sourceCountry at step %d", sourceCountry.Step)

	server.Info.Printf("targetCountry to parse %s", targetCountry.Name)
	server.Info.Printf("targetCountry nbBodies to parse %d", targetCountry.NbBodies)
	server.Info.Printf("targetCountry at step to parse %d", targetCountry.Step)

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
