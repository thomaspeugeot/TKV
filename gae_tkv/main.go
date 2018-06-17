// app engine for 10kt runtime server

package main

import (
	"fmt"
	"net/http"

	"github.com/thomaspeugeot/tkv/handler"

	"github.com/thomaspeugeot/tkv/translation"

	"google.golang.org/appengine" // Required external App Engine library
	"google.golang.org/appengine/log"
)

// singloton pattern to init translation
func getTranslateCurrent() *translation.Translation {
	if handler.TranslateCurrent.GetSourceCountryName() != "fra" {
		var sourceCountry translation.Country
		var targetCountry translation.Country

		sourceCountry.Name = "fra"
		sourceCountry.NbBodies = 697529
		sourceCountry.Step = 4723

		targetCountry.Name = "hti"
		targetCountry.NbBodies = 927787
		targetCountry.Step = 8564

		handler.TranslateCurrent.Init(sourceCountry, targetCountry)
	}

	return &handler.TranslateCurrent
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main() // Starts the server to receive requests
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Index handler")

	t := getTranslateCurrent()
	log.Infof(ctx, "Translation inited %s", t.GetSourceCountryName())

	fmt.Fprintln(w, "Translation inited!")
}
