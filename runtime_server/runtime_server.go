package main

import (
	"github.com/thomaspeugeot/tkv/server"
	"net/http"
	"log"
)

func main() {

	port := "localhost:8001"

	server.Info.Printf("begin listen on port %s", port)
	mux := http.NewServeMux()
	// mux.Handle("/", http.FileServer(http.Dir("../end_user/")) )
	mux.Handle("/", http.FileServer(http.Dir("../leaflets/")) )
	log.Fatal(http.ListenAndServe(port, mux))
	server.Info.Printf("end")
}

