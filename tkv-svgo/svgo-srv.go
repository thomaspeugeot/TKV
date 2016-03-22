package main

import (
    "log"
    "github.com/ajstarks/svgo/float"
    "net/http"
)

func main() {
    http.Handle("/circle", http.HandlerFunc(circle))
    err := http.ListenAndServe(":2003", nil)
    if err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}

func circle(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "image/svg+xml")
  s := svg.New(w)
  s.Start(504.3, 500)
  s.Circle(250, 250, 125, "fill:none;stroke:black")
  s.End()
}