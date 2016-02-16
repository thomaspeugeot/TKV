package main

import (
	"github.com/thomaspeugeot/tkv/barnes-hut"
	"github.com/thomaspeugeot/tkv/quadtree"
	// "testing"
	"fmt"
	"log"
	"net/http"
	"os"
	// "math/rand"
)

//!+main
var r barnes_hut.Run

func main() {
	
	var bodies []quadtree.Body
	quadtree.InitBodiesUniform( &bodies, 200000)

	barnes_hut.SpreadOnCircle( & bodies)
	
	r.Init( & bodies)

	output, _ := os.Create("essai200Kbody_6Ksteps.gif")

	go r.OutputGif( output, 6000)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/status", status)
	mux.HandleFunc("/play", play)
	mux.HandleFunc("/pause", pause)
	mux.HandleFunc("/render", render)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
//!-main

func status(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Run status %s\n", r.GetState())
}

func play(w http.ResponseWriter, req *http.Request) {
	r.SetState( barnes_hut.RUNNING)
	fmt.Fprintf(w, "Run status %s\n", r.GetState())
}

func pause(w http.ResponseWriter, req *http.Request) {
	r.SetState( barnes_hut.STOPPED)
	fmt.Fprintf(w, "Run status %s\n", r.GetState())
}

func render(w http.ResponseWriter, req *http.Request) {

	r.RenderGif( w)
	fmt.Fprintf(w, "Run status %s\n", r.GetState())
}


