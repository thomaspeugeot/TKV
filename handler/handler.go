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
type test_struct struct {
	X1, X2, Y1, Y2 float64
}

type LatLng struct {
	Lat, Lng float64
}

var lastReqest LatLng // store last request.

type VillageCoordResponse struct {
	X, Y                   float64
	Distance               float64
	LatClosest, LngClosest float64
	LatTarget, LngTarget   float64
	Xspread, Yspread       float64
}

// get village coordinates from lat/long
func TranslateLatLngInSourceCountryToLatLngInTargetCountry(w http.ResponseWriter, req *http.Request) {

	server.Info.Printf("translateLatLngInSourceCountryToLatLngInTargetCountry begin")

	// parse lat long from client
	decoder := json.NewDecoder(req.Body)
	var ll LatLng
	err := decoder.Decode(&ll)
	if err != nil {
		log.Println("error decoding ", err)
	} else {
		lastReqest = ll
	}
	server.Info.Printf("translateLatLngInSourceCountryToLatLngInTargetCountry for lat %f, lng %f", ll.Lat, ll.Lng)

	x, y, distance, latClosest, lngClosest, xSpread, ySpread, _ := translation.GetTranslateCurrent().ClosestBodyInOriginalPosition(ll.Lat, ll.Lng)
	server.Info.Printf("translateLatLngInSourceCountryToLatLngInTargetCountry x, y is %f %f, distance %f", x, y, distance)

	var xy VillageCoordResponse
	xy.X = x
	xy.Y = y
	xy.Distance = distance
	xy.LatClosest = latClosest
	xy.LngClosest = lngClosest
	xy.Xspread = xSpread
	xy.Yspread = ySpread

	latTarget, lngTarget := translation.GetTranslateCurrent().XYSpreadToLatLngInTargetCountry(xSpread, ySpread)
	xy.LatTarget = latTarget
	xy.LngTarget = lngTarget

	VillageCoordResponsejson, _ := json.MarshalIndent(xy, "", "	")
	fmt.Fprintf(w, "%s", VillageCoordResponsejson)

	server.Info.Printf("translateLatLngInSourceCountryToLatLngInTargetCountry end")
}

// return all points within source borders
// get village coordinates from lat/long
func AllSourcPointsCoordinates(w http.ResponseWriter, req *http.Request) {

	server.Info.Printf("allSourcPointsCoordinates begin")

	// parse lat long from client
	decoder := json.NewDecoder(req.Body)
	var ll LatLng
	err := decoder.Decode(&ll)
	if err != nil {
		log.Println("error decoding ", err)
		ll = lastReqest
	}
	server.Info.Printf("allSourcPointsCoordinates for lat %f, lng %f", ll.Lat, ll.Lng)

	points := translation.GetTranslateCurrent().SourceBorder(ll.Lat, ll.Lng)

	coord := make(GeoJSONBorderCoordinates, 1)
	coord[0] = make([][]float64, len(points))
	for idx, _ := range points {
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
	var ll LatLng
	err := decoder.Decode(&ll)
	if err != nil {
		log.Println("error decoding ", err)
	}
	server.Info.Printf("villageTargetBorder for lat %f, lng %f", ll.Lat, ll.Lng)

	x, y, distance, _, _, xSpread, ySpread, _ := translation.GetTranslateCurrent().ClosestBodyInOriginalPosition(ll.Lat, ll.Lng)
	server.Info.Printf("villageTargetBorder is %f %f, distance %f", x, y, distance)

	points := translation.GetTranslateCurrent().TargetBorder(xSpread, ySpread)

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
	var ll LatLng
	err := decoder.Decode(&ll)
	if err != nil {
		log.Println("error decoding ", err)
	}
	server.Info.Printf("villageSourceBorder for lat %f, lng %f", ll.Lat, ll.Lng)

	points := translation.GetTranslateCurrent().SourceBorder(ll.Lat, ll.Lng)

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
	server.Info.Printf("villageTargetBorder time to compute convex hull %s", TT.String())

	server.Info.Printf("Lower# %d", len(lower))
	server.Info.Printf("Upper# %d", len(upper))

	VillageBorderResponsejson, _ := json.MarshalIndent(PQtoGeoJSONBorderCoordinates(lower, upper), "", "	")
	fmt.Fprintf(w, "%s", VillageBorderResponsejson)

	server.Info.Printf("villageSourceBorder end")
}

// convert pointList to array of array of array of float
// this is necessary since the client only understand a border expressed as [][][]float
type GeoJSONBorderCoordinates [][][]float64

func PQtoGeoJSONBorderCoordinates(lower, upper []pq.Point2q) GeoJSONBorderCoordinates {

	coord := make(GeoJSONBorderCoordinates, 1)
	coord[0] = make([][]float64, len(lower)+len(upper))
	for idx, _ := range lower {
		coord[0][idx] = make([]float64, 2)
		coord[0][idx][0] = lower[idx].Y().Float64() // Y is longitude
		coord[0][idx][1] = lower[idx].X().Float64() // X is latitude
	}
	for idx, _ := range upper {
		coord[0][len(lower)+idx] = make([]float64, 2)
		coord[0][len(lower)+idx][0] = upper[idx].Y().Float64() // Y is longitude
		coord[0][len(lower)+idx][1] = upper[idx].X().Float64() // X is latitude
	}

	return coord
}
