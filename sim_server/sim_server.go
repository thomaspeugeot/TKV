package main

import (
	"github.com/thomaspeugeot/tkv/barnes-hut"
	"github.com/thomaspeugeot/tkv/quadtree"
	// "testing"
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
)

//!+main
var r barnes_hut.Run

func main() {
	
	var bodies []quadtree.Body
	quadtree.InitBodiesUniform( &bodies, 200000)

	barnes_hut.SpreadOnCircle( & bodies)
	
	r.Init( & bodies)

	output, _ := os.Create("essai200Kbody_6Ksteps.gif")

	go r.OutputGif( output, 15000)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/status", status)
	mux.HandleFunc("/play", play)
	mux.HandleFunc("/pause", pause)
	mux.HandleFunc("/render", render)
	mux.HandleFunc("/stats", stats)
	mux.HandleFunc("/area", area)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
//!-main

func status(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Run status %s step %d\n", r.GetState(), r.GetStep())
}

func play(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.SetState( barnes_hut.RUNNING)
	fmt.Fprintf(w, "Run status %s\n", r.GetState())
}

func pause(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.SetState( barnes_hut.STOPPED)
	fmt.Fprintf(w, "Run status %s\n", r.GetState())
}

func render(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.RenderGif( w)
	// fmt.Fprintf(w, "Run status %s\n", r.GetState())
}

func stats(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// stats, _ := json.MarshalIndent( r.BodyCountGini(), "", "	")
	// stats, _ := json.MarshalIndent( r.GiniOverTimeTransposed(), "", "	")
	stats, _ := json.MarshalIndent( r.GiniOverTime(), "","\t")
	// fmt.Println( string( stats))
	fmt.Fprintf(w, "%s", stats)
}

type test_struct struct {
	x1, x2, y1, y2 float64
}

func area(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println( "Path ", req.URL.Path)
	fmt.Println( "Header ", req.Header)
	fmt.Println( "Form ", req.Form)
	fmt.Println( "PostForm ", req.PostForm)
	fmt.Println( "Body ",  req.Body)

	decoder := json.NewDecoder( req.Body)
	var t test_struct
	err := decoder.Decode( &t)
	if err != nil {
		log.Println(err)
	}
	log.Println(t.x1)

}



