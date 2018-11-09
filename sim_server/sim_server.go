/*
Contains the main package for the simulation server. Also contains all handlers for the simulation server.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/thomaspeugeot/tkv/barnes-hut"
	"github.com/thomaspeugeot/tkv/server"
	"github.com/thomaspeugeot/tkv/translation"
)

//!+main
var r *barneshut.Run

//
// to start with haiti
// go run sim_server.go -sourceCountry=hti -sourceCountryNbBodies=82990
//
func main() {

	// flags  for source country
	sourceCountryPtr := flag.String("sourceCountry", "fra", "iso 3166 sourceCountry code")
	sourceCountryNbBodiesPtr := flag.String("sourceCountryNbBodies", "934136", "nb of bodies")
	sourceCountryStepPtr := flag.String("sourceCountryStep", "0", "simulation step for the spread bodies for source country")

	cutoffPtr := flag.String("cutoff", "2", "cutoff code distance")

	shutdownCriteriaPtr := flag.String("shutdownCriteria", "0.00001", "If energy decreases ratio is below this threshold during a simulation step, simulation shutdowns")

	portPtr := flag.String("port", "8000", "listening port")

	startPtr := flag.Bool("start", false, "if true, start simulation run immediatly") 
	
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
	{
		_, errScan := fmt.Sscanf(*cutoffPtr, "%f", &barneshut.CutoffDistance)
		if errScan != nil {
			log.Fatal(errScan)
			return
		}
	}
	server.Info.Printf("CutoffDistance %f", barneshut.CutoffDistance)

	{
		_, errScan := fmt.Sscanf(*shutdownCriteriaPtr, "%f", &barneshut.ShutdownCriteria)
		if errScan != nil {
			log.Fatal(errScan)
			return
		}
	}
	server.Info.Printf("Studown Criteria %f", barneshut.ShutdownCriteria)
	port := 8000
	{
		_, errScan := fmt.Sscanf(*portPtr, "%d", &port)
		if errScan != nil {
			log.Fatal(errScan)
			return
		}
	}
	server.Info.Printf("will listen on port %d", port)
	r = barneshut.NewRun()

	// load configuration files.
	filename := fmt.Sprintf(barneshut.CountryBodiesNamePattern, sourceCountry.Name, sourceCountry.NbBodies, sourceCountry.Step)
	server.Info.Printf("filename for init %s", filename)
	r.LoadConfig(filename)

	if !*startPtr	{
		r.SetState(barneshut.STOPPED)
	} else {
		r.SetState(barneshut.RUNNING)
	}
	
	go r.RunSimulation()

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
	mux.HandleFunc("/minDistanceCoord", minDistanceCoord)
	mux.HandleFunc("/nbVillagesPerAxe", nbVillagesPerAxe)
	mux.HandleFunc("/nbRoutines", nbRoutines)
	mux.HandleFunc("/fieldGridNb", fieldGridNb)
	mux.HandleFunc("/updateRatioBorderBodies", updateRatioBorderBodies)
	mux.HandleFunc("/toggleRenderChoice", toggleRenderChoice)
	mux.HandleFunc("/toggleFieldRendering", toggleFieldRendering)

	mux.Handle("/", http.FileServer(http.Dir("../tkv-client/")))
	adressToListen := fmt.Sprintf("localhost:%d", port)
	server.Info.Printf("adressToListen %s", adressToListen)
	log.Fatal(http.ListenAndServe(adressToListen, mux))
}

//!-main

func status(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "%s Dt Adjust %s\n%s",
		r.State(),
		barneshut.DtAdjustMode,
		r.Status())
}

func play(w http.ResponseWriter, req *http.Request) {

	r.SetState(barneshut.RUNNING)
	fmt.Fprintf(w, "Run status %s\n", r.State())
}

func toggleRenderChoice(w http.ResponseWriter, req *http.Request)   { r.ToggleRenderChoice() }
func toggleFieldRendering(w http.ResponseWriter, req *http.Request) { r.ToggleFieldRendering() }
func toggleManualAuto(w http.ResponseWriter, req *http.Request)     { r.ToggleManualAuto() }

func pause(w http.ResponseWriter, req *http.Request) {

	r.SetState(barneshut.STOPPED)
	fmt.Fprintf(w, "Run status %s\n", r.State())
}

func oneStep(w http.ResponseWriter, req *http.Request) {
	if r.State() == barneshut.STOPPED {
		r.OneStep()
	}
	fmt.Fprintf(w, "Run status %s\n", r.State())
}

func captureConfig(w http.ResponseWriter, req *http.Request) {
	r.CaptureConfig()
}

func render(w http.ResponseWriter, req *http.Request)    { r.RenderGif(w, true) }
func renderSVG(w http.ResponseWriter, req *http.Request) { r.RenderSVG(w) }

func stats(w http.ResponseWriter, req *http.Request) {

	stats, _ := json.MarshalIndent(r.BodyCountGini(), "", "	")
	// stats, _ := json.MarshalIndent( r.GiniOverTimeTransposed(), "", "	")
	// fmt.Println( string( stats))
	fmt.Fprintf(w, "%s", stats)
}

func getDensityTenciles(w http.ResponseWriter, req *http.Request) {

	tenciles, _ := json.MarshalIndent(r.ComputeDensityTencilePerTerritoryString(), "", "	")
	fmt.Fprintf(w, "%s", tenciles)
}

type testStruct struct {
	X1, X2, Y1, Y2 float64
}

func area(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t testStruct
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("error decoding ", err)
	}
	r.SetRenderingWindow(t.X1, t.X2, t.Y1, t.Y2)
}

func dt(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var dtRequest float64
	err := decoder.Decode(&dtRequest)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		barneshut.DtRequest = dtRequest
	}
}

func theta(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var thetaRequest float64
	err := decoder.Decode(&thetaRequest)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		barneshut.BN_THETA_Request = thetaRequest
	}
}

func nbVillagesPerAxe(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var nbVillagesPerAxe int
	err := decoder.Decode(&nbVillagesPerAxe)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		barneshut.SetNbVillagePerAxe(nbVillagesPerAxe)
	}
}

func nbRoutines(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var nbRoutines int
	err := decoder.Decode(&nbRoutines)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		barneshut.SetNbRoutines(nbRoutines)
	}
}

func fieldGridNb(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var gridNb float64
	err := decoder.Decode(&gridNb)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		r.SetGridFieldNb(int(math.Floor(gridNb)))
	}
}

func updateRatioBorderBodies(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var ratioBorderBodies float64
	err := decoder.Decode(&ratioBorderBodies)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		barneshut.SetRatioBorderBodies(ratioBorderBodies)
	}
}

// list the content of the available config files
func dirConfig(w http.ResponseWriter, req *http.Request) {

	dircontent, _ := json.MarshalIndent(r.DirConfig(), "", "	")
	fmt.Fprintf(w, "%s", dircontent)
}

// send coordinates of minimal distance
func minDistanceCoord(w http.ResponseWriter, req *http.Request) {

	minDistanceCoordResp, _ := json.MarshalIndent(r.GetMaxRepulsiveForce(), "", "	")
	fmt.Fprintf(w, "%s", minDistanceCoordResp)
}

// load config files
func loadConfig(w http.ResponseWriter, req *http.Request) {

	// get the file
	fileSlice := req.URL.Query()["file"]

	server.Info.Println(fileSlice[0])
	// get the file name

	loadResult := r.LoadConfig(fileSlice[0])
	server.Info.Println("load result ", loadResult)
}

// list config files in orig
func loadConfigOrig(w http.ResponseWriter, req *http.Request) {

	// get the file
	fileSlice := req.URL.Query()["file"]

	server.Info.Println(fileSlice[0])
	// get the file name

	loadResult := r.LoadConfigOrig(fileSlice[0])
	server.Info.Println("load result ", loadResult)
}
