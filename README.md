10 0000
=======

[![Go Report Card](https://goreportcard.com/badge/github.com/thomaspeugeot/tkv)](https://goreportcard.com/report/github.com/thomaspeugeot/tkv)
[![Godoc](https://img.shields.io/badge/godoc-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/thomaspeugeot/tkv)


This code is the implementation of the "10 000" concept (see https://10ktblog.wordpress.com/a-propos/ for a description of the concept). 

You can check the end product at https://tenktorg.appspot.com/. Zoom in the top area (france). Click, get your territory and check the twin territory in Haïti. 

This reopository programs related
* the "extractor" program that turns an open source density file from a country into a country body file at initial configuration
* the "simulation" programm that simulates the spreading of bodies (it takes a country body file and output a new one with updated positions)  
* the "runtime" 10 000 web server for the end user who need to find his territory in france among the 10 000 territories and the sister territory in Haiti. This program take a 2 country body files (one at init and one at the end of the simulation)



The 10000 runtime server
-------------------------

**Running the web server with go command tool**

Have go (golang.org) latest (>= v1.11) installed
blank GOPATH & GOROOT env

```
cd
go get github.com/gyuho/goraph
go get github.com/thomaspeugeot/tkv
go get github.com/thomaspeugeot/pq
go get github.com/ajstarks/svgo
cd go/src/github.com/thomaspeugeot/tkv/runtime_server
go run runtime_server.go -sourceCountry fra -sourceCountryNbBodiesPtr 697529 -sourceCountryStep 4723 -targetCountry hti -targetCountryNbBodiesPtr 927787 -targetCountryStep 8564
```

a vscode configuration is available to run and debug the server.

**Running the web client**


launch your browser at http://localhost:8002/tkv-client.html

On the top panel, zoom to your place of interest (it is currently limited to france). Left click. You terrritory appears as well as the matching territory in Haiti.


The extractor program
-------------------------
You can run with default parameters
```
cd grump-reader
go run grump-reader.go -tkvdata="C:\Users\peugeot\tkv-data"
```
to see the execution flags
```
cd grump-reader
go run grump-reader.go -help
```

depending on the input country the program exectutes in less that a minute

The simulation server
-------------------------

```
cd sim_server
go run sim_server.go -sourceCountry=hti -sourceCountryNbBodies=82990
```

you can monitor sim_server progress running by opening the file  tkv-client/tkv-monitor.html in your favorite browser


**The "movie" program**

This program generates a movie from the simulation steps
```
cd sim-movie
go run sim-movie.go <some flags to be completed>
```