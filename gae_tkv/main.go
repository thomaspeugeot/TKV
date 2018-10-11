// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
}

// // app engine for 10kt runtime server

// package main

// import (
// 	"net/http"

// 	// "github.com/thomaspeugeot/tkv/handler"
// 	// "github.com/thomaspeugeot/tkv/handler"
// 	// "github.com/thomaspeugeot/tkv/translation"

// 	"google.golang.org/appengine" // Required external App Engine library
// )

// // singloton pattern to init translation
// // func getTranslateCurrent() *translation.Translation {
// // 	if handler.TranslateCurrent.GetSourceCountryName() != "fra" {
// // 		var sourceCountry translation.Country
// // 		var targetCountry translation.Country

// // 		sourceCountry.Name = "fra"
// // 		sourceCountry.NbBodies = 697529
// // 		sourceCountry.Step = 4723

// // 		targetCountry.Name = "hti"
// // 		targetCountry.NbBodies = 927787
// // 		targetCountry.Step = 8564

// // 		handler.TranslateCurrent.Init(sourceCountry, targetCountry)
// // 	}

// // 	return &handler.TranslateCurrent
// // }

// func main() {
// 	//	the example line of code

// 	//	http.HandleFunc("/", indexHandler)
// 	// http.HandleFunc("/", indexHandler)

// 	mux := http.NewServeMux()
// 	// mux.Handle("/", http.FileServer(http.Dir(".")))
// 	http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))

// 	// mux.HandleFunc("/translateLatLngInSourceCountryToLatLngInTargetCountry",
// 	// 	handler.TranslateLatLngInSourceCountryToLatLngInTargetCountry)
// 	// mux.HandleFunc("/villageTargetBorder", handler.VillageTargetBorder)
// 	// mux.HandleFunc("/villageSourceBorder", handler.VillageSourceBorder)
// 	// mux.HandleFunc("/allSourcPointsCoordinates", handler.AllSourcPointsCoordinates)

// 	appengine.Main() // Starts the server to receive requests
// }

// // [END main_func]
// // [START indexHandler]
// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	// if statement redirects all invalid URLs to the root homepage.
// 	// Ex: if URL is http://[YOUR_PROJECT_ID].appspot.com/FOO, it will be
// 	// redirected to http://[YOUR_PROJECT_ID].appspot.com.
// 	if r.URL.Path != "/" {
// 		http.Redirect(w, r, "/", http.StatusFound)
// 		return
// 	}

// 	// fmt.Fprintln(w, "Hello 10000 !")
// }

// // func indexHandler(w http.ResponseWriter, r *http.Request) {

// // 	ctx := appengine.NewContext(r)
// // 	log.Infof(ctx, "Index handler")

// // 	t := getTranslateCurrent()
// // 	log.Infof(ctx, "Translation inited %s", t.GetSourceCountryName())

// // 	fmt.Fprintln(w, "Translation inited!")
// // }
