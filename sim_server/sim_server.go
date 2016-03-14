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
	mux.HandleFunc("/oneStep", oneStep)
	mux.HandleFunc("/captureConfig", captureConfig)
	mux.HandleFunc("/render", render)
	mux.HandleFunc("/stats", stats)
	mux.HandleFunc("/area", area)
	mux.HandleFunc("/dt", dt)
	mux.HandleFunc("/theta", theta)
	mux.Handle("/", http.FileServer(http.Dir("../tkv-client/")) )
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
//!-main

func status(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Run status %s step %d\n", r.State(), r.GetStep())
}

func play(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.SetState( barnes_hut.RUNNING)
	fmt.Fprintf(w, "Run status %s\n", r.State())
}

func pause(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.SetState( barnes_hut.STOPPED)
	fmt.Fprintf(w, "Run status %s\n", r.State())
}

func oneStep(w http.ResponseWriter, req *http.Request) {
	if (r.State() == barnes_hut.STOPPED) {
		r.OneStep()
	}
	fmt.Fprintf(w, "Run status %s\n", r.State())
}

func captureConfig(w http.ResponseWriter, req *http.Request) {
	r.CaptureConfig()
}

func render(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.RenderGif( w)
	// fmt.Fprintf(w, "Run status %s\n", r.State())
}

func stats(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// stats, _ := json.MarshalIndent( r.BodyCountGini(), "", "	")
	// stats, _ := json.MarshalIndent( r.GiniOverTimeTransposed(), "", "	")
	// fmt.Println( string( stats))
	fmt.Fprintf(w, "%s", stats)
}

type test_struct struct {
	X1, X2, Y1, Y2 float64
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
		log.Println("error decoding ", err)
	}
	r.SetRenderingWindow( t.X1, t.X2, t.Y1, t.Y2)
	log.Println(t.X1)
	log.Println(t.X2)
	log.Println(t.Y1)
	log.Println(t.Y2)

}

func dt(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	decoder := json.NewDecoder( req.Body)
	var dtRequest float64
	err := decoder.Decode( &dtRequest)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		barnes_hut.DtRequest = dtRequest
	}
}


func theta(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	decoder := json.NewDecoder( req.Body)
	var thetaRequest float64
	err := decoder.Decode( &thetaRequest)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		barnes_hut.ThetaRequest = thetaRequest
	}
}

