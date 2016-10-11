package main

import (
	"github.com/thomaspeugeot/tkv/server"
	"github.com/thomaspeugeot/tkv/translation"
	"net/http"
	"log"
	"flag"
	"fmt"
	"encoding/json"
	"github.com/thomaspeugeot/pq"
	"time"

)


// load the data
var t translation.Translation

// for decoding the rendering window
type test_struct struct {
	X1, X2, Y1, Y2 float64
}
//
//
// on pc
// go run runtime_server.go -targetCountryStep=43439
// 
func main() {

	// flags  for source country
	sourceCountryPtr := flag.String("sourceCountry","fra","iso 3166 sourceCountry code")
	sourceCountryNbBodiesPtr := flag.String("sourceCountryNbBodiesPtr","34413","nb of bodies")
	sourceCountryStepPtr := flag.String("sourceCountryStep","3563","simulation step for the spread bodies for source country")

	// flags  for target country
	targetCountryPtr := flag.String("targetCountry","hti","iso 3166 targetCountry code")
	targetCountryNbBodiesPtr := flag.String("targetCountryNbBodiesPtr","82990","nb of bodies for target country")
	targetCountryStepPtr := flag.String("targetCountryStep","36719","simulation step for the spread bodies for target country")

	flag.Parse()

	// init sourceCountry from flags
	var sourceCountry translation.Country
	sourceCountry.Name = *sourceCountryPtr
	{
		_, errScan := fmt.Sscanf(*sourceCountryNbBodiesPtr, "%d", & sourceCountry.NbBodies)
		if( errScan != nil) {
			log.Fatal(errScan)
			return			
		}
	}
	{
		_, errScan := fmt.Sscanf(*sourceCountryStepPtr, "%d", & sourceCountry.Step)
		if( errScan != nil) {
			log.Fatal(errScan)
			return			
		}
	}

	// init targetCountry from flags
	var targetCountry translation.Country
	targetCountry.Name = *targetCountryPtr
	{
		_, errScan := fmt.Sscanf(*targetCountryNbBodiesPtr, "%d", & targetCountry.NbBodies)
		if( errScan != nil) {
			log.Fatal(errScan)
			return			
		}
	}
	{
		_, errScan := fmt.Sscanf(*targetCountryStepPtr, "%d", & targetCountry.Step)
		if( errScan != nil) {
			log.Fatal(errScan)
			return			
		}
	}


	server.Info.Printf("sourceCountry to parse %s", sourceCountry.Name)
	server.Info.Printf("nbBodies to parse %d", sourceCountry.NbBodies)
	server.Info.Printf("step to parse %d", sourceCountry.Step)
	
	server.Info.Printf("targetCountry to parse %s", targetCountry.Name)
	server.Info.Printf("nbBodies to parse %d", targetCountry.NbBodies)
	server.Info.Printf("step to parse %d", targetCountry.Step)
	
	t.Init(sourceCountry, targetCountry)

	port := "localhost:8001"

	server.Info.Printf("begin listen on port %s", port)
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("../tkv-client/")) )
	
	mux.HandleFunc("/villageCoordinates", villageCoordinates)
	mux.HandleFunc("/villageTargetBorder", villageTargetBorder)
	mux.HandleFunc("/villageSourceBorder", villageSourceBorder)
		
	log.Fatal(http.ListenAndServe(port, mux))
	server.Info.Printf("end")
}

type LatLng struct {
	Lat, Lng float64
}

type VillageCoordResponse struct {
	X, Y float64
	Distance float64
	LatClosest, LngClosest float64
	LatTarget, LngTarget float64	
}

// get village coordinates from lat/long
func villageCoordinates(w http.ResponseWriter, req *http.Request) {
	
	server.Info.Printf("villageCoordinates begin")
	
	// parse lat long from client
	decoder := json.NewDecoder( req.Body)
	var ll LatLng
	err := decoder.Decode( &ll)
	if err != nil {
		log.Println("error decoding ", err)
	}
	server.Info.Printf("villageCoordinates for lat %f, lng %f", ll.Lat, ll.Lng)

	x, y, distance, latClosest, lngClosest, xSpread, ySpread, _ := t.VillageCoordinates( ll.Lat, ll.Lng)
	server.Info.Printf("villageCoordinates is %f %f, distance %f", x, y, distance)

	var xy VillageCoordResponse
	xy.X = x
	xy.Y = y
	xy.Distance = distance
	xy.LatClosest = latClosest
	xy.LngClosest = lngClosest

	latTarget, lngTarget := t.TargetVillage( xSpread, ySpread)
	xy.LatTarget = latTarget
	xy.LngTarget = lngTarget

	VillageCoordResponsejson, _ := json.MarshalIndent( xy, "", "	")
	fmt.Fprintf(w, "%s", VillageCoordResponsejson)
	
	server.Info.Printf("villageCoordinates end")
}

// get target village border from lat/long
func villageTargetBorder(w http.ResponseWriter, req *http.Request) {

	// parse lat long from client
	decoder := json.NewDecoder( req.Body)
	var ll LatLng
	err := decoder.Decode( &ll)
	if err != nil {
		log.Println("error decoding ", err)
	}
	server.Info.Printf("villageTargetBorder for lat %f, lng %f", ll.Lat, ll.Lng)
	
	x, y, distance, _, _, xSpread, ySpread, _ := t.VillageCoordinates( ll.Lat, ll.Lng)
	server.Info.Printf("villageTargetBorder is %f %f, distance %f", x, y, distance)

	
	points := t.TargetBorder( xSpread, ySpread)
	
	// availble convex hull code (in perfect precision but robust)
	ps := make([]pq.Point2q, len(points))
	for i := 0; i < len(points); i++ {
		//
		xf, yf := points[i].X, points[i].Y
		//
		xq, yq := pq.FtoQ(xf), pq.FtoQ(yf)
		ps[i] = pq.XYtoP(xq, yq)
	}
	
	T := time.Now()
	lower, upper := pq.ParConvHull2q(4, ps)
	TT := time.Since(T)
	server.Info.Printf("villageTargetBorder time to compute convex hull %s", time.Now().Format("2006-01-02 15:04:05") + " Ready; T=" + TT.String())
	
	server.Info.Printf("Lower# %d", len(lower))
	server.Info.Printf("Upper# %d", len(upper))

	PQtoGeoJSONBorderCoordinates(lower, upper)

	VillageBorderResponsejson, _ := json.MarshalIndent( PQtoGeoJSONBorderCoordinates(lower, upper), "", "	")
	fmt.Fprintf(w, "%s", VillageBorderResponsejson)
}

// get target village border from lat/long
func villageSourceBorder(w http.ResponseWriter, req *http.Request) {

	server.Info.Printf("villageSourceBorder begin")
	
	// parse lat long from client
	decoder := json.NewDecoder( req.Body)
	var ll LatLng
	err := decoder.Decode( &ll)
	if err != nil {
		log.Println("error decoding ", err)
	}
	server.Info.Printf("villageSourceBorder for lat %f, lng %f", ll.Lat, ll.Lng)
	
	points := t.SourceBorder( ll.Lat, ll.Lng)

		// availble convex hull code (in perfect precision but robust)
	ps := make([]pq.Point2q, len(points))
	for i := 0; i < len(points); i++ {
		//
		xf, yf := points[i].X, points[i].Y
		//
		xq, yq := pq.FtoQ(xf), pq.FtoQ(yf)
		ps[i] = pq.XYtoP(xq, yq)
	}
	
	T := time.Now()
	lower, upper := pq.ParConvHull2q(4, ps)
	TT := time.Since(T)
	server.Info.Printf("villageTargetBorder time to compute convex hull %s", time.Now().Format("2006-01-02 15:04:05") + " Ready; T=" + TT.String())
	
	server.Info.Printf("Lower# %d", len(lower))
	server.Info.Printf("Upper# %d", len(upper))

	VillageBorderResponsejson, _ := json.MarshalIndent( PQtoGeoJSONBorderCoordinates(lower, upper), "", "	")
	fmt.Fprintf(w, "%s", VillageBorderResponsejson)
	
	server.Info.Printf("villageSourceBorder end")
}


// convert pointList to array of array of array of float
// this is necessary since the client only understand a border expressed as [][][]float
type GeoJSONBorderCoordinates [][][]float64

func PQtoGeoJSONBorderCoordinates(lower, upper []pq.Point2q) GeoJSONBorderCoordinates {

	coord := make( GeoJSONBorderCoordinates, 1)
	coord[0] = make( [][]float64, len(lower)+len(upper))
	for idx, _ := range lower {
		coord[0][idx] = make( []float64, 2)
		coord[0][idx][0] = lower[idx].Y().Float64() // Y is longitude
		coord[0][idx][1] = lower[idx].X().Float64() // X is latitude
	}
	for idx, _ := range upper{
		coord[0][len(lower)+idx] = make( []float64, 2)
		coord[0][len(lower)+idx][0] = upper[idx].Y().Float64() // Y is longitude
		coord[0][len(lower)+idx][1] = upper[idx].X().Float64() // X is latitude
	}
	
	return coord
}
