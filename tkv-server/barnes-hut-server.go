package main

import (
	"github.com/thomaspeugeot/tkv/barnes-hut"
	"github.com/thomaspeugeot/tkv/quadtree"
	// "testing"
	"fmt"
	"log"
	"net/http"
	// "os"
	// "math/rand"
)

//!+main
var r barnes_hut.Run

func main() {
	
	var bodies []quadtree.Body
	quadtree.InitBodiesUniform( &bodies, 200000)

	barnes_hut.SpreadOnCircle( & bodies)
	
	r.Init( & bodies)

	mux := http.NewServeMux()
	mux.HandleFunc("/status", status)
	mux.HandleFunc("/render", render)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

func status(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Run status %s\n", r.GetState().String())
}

func render(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Run status %s\n", r.GetState().String())
}

//!-main
