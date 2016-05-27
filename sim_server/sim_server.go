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
var r * barnes_hut.Run

func main() {
	
	r = barnes_hut.NewRun()
	
	var bodies []quadtree.Body
	quadtree.InitBodiesUniform( &bodies, 200000)

	barnes_hut.SpreadOnCircle( & bodies)
	
	r.Init( & bodies)

	output, _ := os.Create("essai200Kbody_6Ksteps.gif")
	go r.OutputGif( output, 100000)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/status", status)

	mux.HandleFunc("/toggleManualAuto", toggleManualAuto)

	mux.HandleFunc("/play", play)
	mux.HandleFunc("/pause", pause)
	mux.HandleFunc("/oneStep", oneStep)
	mux.HandleFunc("/captureConfig", captureConfig)

	mux.HandleFunc("/render", render)
	mux.HandleFunc("/renderSVG", renderSVG)

	mux.HandleFunc("/stats", stats)
	mux.HandleFunc("/area", area)
	mux.HandleFunc("/dt", dt)
	mux.HandleFunc("/theta", theta)
	mux.HandleFunc("/dirConfig", dirConfig)
	mux.HandleFunc("/loadConfig", loadConfig)
	mux.HandleFunc("/loadConfigOrig", loadConfigOrig)
	mux.HandleFunc("/getDensityTenciles", getDensityTenciles)
	mux.HandleFunc("/nbVillagesPerAxe", nbVillagesPerAxe)
	mux.HandleFunc("/nbRoutines", nbRoutines)
	mux.HandleFunc("/updateRatioBorderBodies", updateRatioBorderBodies)
	mux.HandleFunc("/toggleRenderChoice", toggleRenderChoice)

	mux.Handle("/", http.FileServer(http.Dir("../tkv-client/")) )
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
//!-main

func status(w http.ResponseWriter, req *http.Request) {
	
	fmt.Fprintf(w, "%s Dt Adjust %s\n%s", 
				r.State(), 
				barnes_hut.DtAdjustMode,
				r.Status())
}

func play(w http.ResponseWriter, req *http.Request) {
		
	r.SetState( barnes_hut.RUNNING)
	fmt.Fprintf(w, "Run status %s\n", r.State())
}

func toggleRenderChoice(w http.ResponseWriter, req *http.Request) { r.ToggleRenderChoice() }
func toggleManualAuto(w http.ResponseWriter, req *http.Request) { r.ToggleManualAuto() }

func pause(w http.ResponseWriter, req *http.Request) {
	
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

func render(w http.ResponseWriter, req *http.Request) { r.RenderGif( w) }
func renderSVG(w http.ResponseWriter, req *http.Request) { r.RenderSVG( w) }

func stats(w http.ResponseWriter, req *http.Request) {
	
	stats, _ := json.MarshalIndent( r.BodyCountGini(), "", "	")
	// stats, _ := json.MarshalIndent( r.GiniOverTimeTransposed(), "", "	")
	// fmt.Println( string( stats))
	fmt.Fprintf(w, "%s", stats)
}

func getDensityTenciles(w http.ResponseWriter, req *http.Request) {
	
	
	tenciles, _ := json.MarshalIndent( r.ComputeDensityTencilePerVillageString(), "", "	")
	fmt.Fprintf(w, "%s", tenciles)
}

type test_struct struct {
	X1, X2, Y1, Y2 float64
}

func area(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder( req.Body)
	var t test_struct
	err := decoder.Decode( &t)
	if err != nil {
		log.Println("error decoding ", err)
	}
	r.SetRenderingWindow( t.X1, t.X2, t.Y1, t.Y2)
}

func dt(w http.ResponseWriter, req *http.Request) {
	
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

	decoder := json.NewDecoder( req.Body)
	var thetaRequest float64
	err := decoder.Decode( &thetaRequest)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		barnes_hut.ThetaRequest = thetaRequest
	}
}

func nbVillagesPerAxe(w http.ResponseWriter, req *http.Request) {
	
	decoder := json.NewDecoder( req.Body)
	var nbVillagesPerAxe int
	err := decoder.Decode( &nbVillagesPerAxe)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		barnes_hut.SetNbVillagePerAxe( nbVillagesPerAxe)
	}
}

func nbRoutines(w http.ResponseWriter, req *http.Request) {
	
	decoder := json.NewDecoder( req.Body)
	var nbRoutines int
	err := decoder.Decode( &nbRoutines)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		barnes_hut.SetNbRoutines( nbRoutines)
	}
}

func updateRatioBorderBodies(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder( req.Body)
	var ratioBorderBodies float64
	err := decoder.Decode( &ratioBorderBodies)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		barnes_hut.SetRatioBorderBodies( ratioBorderBodies)
	}
}

// list the content of the available config files
func dirConfig(w http.ResponseWriter, req *http.Request) {
	
	dircontent, _ := json.MarshalIndent( r.DirConfig(), "", "	")
	fmt.Fprintf(w, "%s", dircontent)
}

// load config files
func loadConfig(w http.ResponseWriter, req *http.Request) {
	
	// get the file
	fileSlice := req.URL.Query()["file"]

	fmt.Println(fileSlice[0])
	// get the file name

	loadResult := r.LoadConfig( fileSlice[0])
	fmt.Println( "load result ", loadResult )
}

// list config files in orig
func loadConfigOrig(w http.ResponseWriter, req *http.Request) {
	
	// get the file
	fileSlice := req.URL.Query()["file"]

	fmt.Println(fileSlice[0])
	// get the file name

	loadResult := r.LoadConfigOrig( fileSlice[0])
	fmt.Println( "load result ", loadResult )
}

