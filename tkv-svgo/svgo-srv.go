/*
Code for a web site that generates Scalable Vector Graphics (SVG). It is work in progress.
*/
package main

import (
	"log"
	"net/http"

	"github.com/ajstarks/svgo/float"
)

func main() {
	http.Handle("/circle", http.HandlerFunc(circle))
	err := http.ListenAndServe(":2003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func circle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(504.3, 500)
	s.Circle(250, 250, 125, "fill:none;stroke:black")
	s.End()
}
