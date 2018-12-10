/*
Package handler contains all the function that are handled by the 10000 runtime server
*/
package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/thomaspeugeot/pq"
	"github.com/thomaspeugeot/tkv/translation"
)

// for decoding the rendering window
type testStruct struct {
	X1, X2, Y1, Y2 float64
}

type LatLngCountry struct {
	Lat, Lng float64
	SourceCountry  string
	TargetCountry string
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

	// setup translation
	translation.GetTranslateCurrent().SetSourceCountry( llc.SourceCountry)
	translation.GetTranslateCurrent().SetTargetCountry( llc.TargetCountry)

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
