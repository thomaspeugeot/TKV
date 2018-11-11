/*
Package handler contains all the function that are handled by the 10000 runtime server
*/
package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thomaspeugeot/pq"
	"github.com/thomaspeugeot/tkv/server"
	"github.com/thomaspeugeot/tkv/translation"
)

// for decoding the rendering window
type testStruct struct {
	X1, X2, Y1, Y2 float64
}

type LatLngCountry struct {
	Lat, Lng float64
	Country  string
}

var lastReqest LatLngCountry // store last request.

type VillageCoordResponse struct {
	Source, Target         string
	Distance               float64
	LatClosest, LngClosest float64
	LatTarget, LngTarget   float64
	X, Y                   float64
	SourceBorderPoints     GeoJSONBorderCoordinates
	TargetBorderPoints     GeoJSONBorderCoordinates
}

// get village coordinates from lat/long
func GetTranslationResult(w http.ResponseWriter, req *http.Request) {

	// parse lat long from client
	decoder := json.NewDecoder(req.Body)
	var llc LatLngCountry
	err := decoder.Decode(&llc)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		lastReqest = llc
	}

	// check wether country is the source country
	// if not, switch source & target
	sourceCountry := translation.GetTranslateCurrent().GetSourceCountryName()
	if llc.Country != sourceCountry {
		translation.GetTranslateCurrent().Swap()
	}

	distance, latClosest, lngClosest, xSpread, ySpread, _ :=
		translation.GetTranslateCurrent().BodyCoordsInSourceCountry(llc.Lat, llc.Lng)

	var response VillageCoordResponse
	response.Source = translation.GetTranslateCurrent().GetSourceCountryName()
	response.Target = translation.GetTranslateCurrent().GetTargetCountryName()
	response.Distance = distance
	response.LatClosest = latClosest
	response.LngClosest = lngClosest
	response.X = xSpread
	response.Y = ySpread

	latTarget, lngTarget := translation.GetTranslateCurrent().LatLngToXYInTargetCountry(xSpread, ySpread)
	response.LatTarget = latTarget
	response.LngTarget = lngTarget

	// add source border
	sourceBorderPoints := translation.GetTranslateCurrent().SourceBorder(llc.Lat, llc.Lng)
	response.SourceBorderPoints = make(GeoJSONBorderCoordinates, 1)
	response.SourceBorderPoints[0] = make([][]float64, len(sourceBorderPoints))
	for idx := range sourceBorderPoints {
		response.SourceBorderPoints[0][idx] = make([]float64, 2)
		response.SourceBorderPoints[0][idx][0] = sourceBorderPoints[idx].Y // Y is longitude
		response.SourceBorderPoints[0][idx][1] = sourceBorderPoints[idx].X // X is latitude
	}

	// add target border
	targetBorderPoints := translation.GetTranslateCurrent().TargetBorder(xSpread, ySpread)
	response.TargetBorderPoints = make(GeoJSONBorderCoordinates, 1)
	response.TargetBorderPoints[0] = make([][]float64, len(targetBorderPoints))
	for idx := range targetBorderPoints {
		response.TargetBorderPoints[0][idx] = make([]float64, 2)
		response.TargetBorderPoints[0][idx][0] = targetBorderPoints[idx].Y // Y is longitude
		response.TargetBorderPoints[0][idx][1] = targetBorderPoints[idx].X // X is latitude
	}

	VillageCoordResponsejson, _ := json.MarshalIndent(response, "", "	")
	fmt.Fprintf(w, "%s", VillageCoordResponsejson)
}

// return all points within source borders
// get village coordinates from lat/long
func AllSourceBorderPointsCoordinates(w http.ResponseWriter, req *http.Request) {

	server.Info.Printf("allSourcPointsCoordinates begin")

	// parse lat long from client
	decoder := json.NewDecoder(req.Body)
	var ll LatLngCountry
	err := decoder.Decode(&ll)
	if err != nil {
		log.Println("error decoding ", err)
		ll = lastReqest
	}
	server.Info.Printf("allSourcPointsCoordinates for lat %f, lng %f", ll.Lat, ll.Lng)

	points := translation.GetTranslateCurrent().SourceBorder(ll.Lat, ll.Lng)

	coord := make(GeoJSONBorderCoordinates, 1)
	coord[0] = make([][]float64, len(points))
	for idx := range points {
		coord[0][idx] = make([]float64, 2)
		coord[0][idx][0] = points[idx].Y // Y is longitude
		coord[0][idx][1] = points[idx].X // X is latitude
	}

	allSourcPointsCoordinatesResponsejson, _ := json.MarshalIndent(coord, "", "	")
	fmt.Fprintf(w, "%s", allSourcPointsCoordinatesResponsejson)

	server.Info.Printf("allSourcPointsCoordinates end")
}

// get target village border from lat/long
func VillageTargetBorder(w http.ResponseWriter, req *http.Request) {

	// parse lat long from client
	decoder := json.NewDecoder(req.Body)
	var ll LatLngCountry
	err := decoder.Decode(&ll)
	if err != nil {
		log.Println("error decoding ", err)
	}
	server.Info.Printf("villageTargetBorder for lat %f, lng %f", ll.Lat, ll.Lng)

	_, _, _, xSpread, ySpread, _ := translation.GetTranslateCurrent().BodyCoordsInSourceCountry(ll.Lat, ll.Lng)

	points := translation.GetTranslateCurrent().TargetBorder(xSpread, ySpread)

	// available convex hull code (in perfect precision but robust)
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
	server.Info.Printf("villageTargetBorder time to compute convex hull %s", TT.String())

	server.Info.Printf("Lower# %d", len(lower))
	server.Info.Printf("Upper# %d", len(upper))

	PQtoGeoJSONBorderCoordinates(lower, upper)

	VillageBorderResponsejson, _ := json.MarshalIndent(PQtoGeoJSONBorderCoordinates(lower, upper), "", "	")
	fmt.Fprintf(w, "%s", VillageBorderResponsejson)
}

// get target village border from lat/long
func VillageSourceBorder(w http.ResponseWriter, req *http.Request) {

	server.Info.Printf("villageSourceBorder begin")

	// parse lat long from client
	decoder := json.NewDecoder(req.Body)
	var ll LatLngCountry
	err := decoder.Decode(&ll)
	if err != nil {
		log.Println("error decoding ", err)
	}
	server.Info.Printf("villageSourceBorder for lat %f, lng %f", ll.Lat, ll.Lng)

	points := translation.GetTranslateCurrent().SourceBorder(ll.Lat, ll.Lng)

	// available convex hull code (in perfect precision but robust)
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
	server.Info.Printf("villageTargetBorder time to compute convex hull %s", TT.String())

	server.Info.Printf("Lower# %d", len(lower))
	server.Info.Printf("Upper# %d", len(upper))

	VillageBorderResponsejson, _ := json.MarshalIndent(PQtoGeoJSONBorderCoordinates(lower, upper), "", "	")
	fmt.Fprintf(w, "%s", VillageBorderResponsejson)

	server.Info.Printf("villageSourceBorder end")
}

// Type GeoJSONBorderCoordinates is an array of an array of an array of int
// convert pointList to array of array of array of float
// this is necessary since the client only understand a border expressed as [][][]float
type GeoJSONBorderCoordinates [][][]float64

func PQtoGeoJSONBorderCoordinates(lower, upper []pq.Point2q) GeoJSONBorderCoordinates {

	coord := make(GeoJSONBorderCoordinates, 1)
	coord[0] = make([][]float64, len(lower)+len(upper))
	for idx := range lower {
		coord[0][idx] = make([]float64, 2)
		coord[0][idx][0] = lower[idx].Y().Float64() // Y is longitude
		coord[0][idx][1] = lower[idx].X().Float64() // X is latitude
	}
	for idx := range upper {
		coord[0][len(lower)+idx] = make([]float64, 2)
		coord[0][len(lower)+idx][0] = upper[idx].Y().Float64() // Y is longitude
		coord[0][len(lower)+idx][1] = upper[idx].X().Float64() // X is latitude
	}

	return coord
}
