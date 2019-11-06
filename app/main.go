package main

import (
    "fmt"
    "log"
    "net/http"
)

const HTTPListenPort = 3434

func rootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Rendering %s", r.URL.Path)
}

func main() {
    http.HandleFunc("/", rootHandler);
    log.Printf("Listening on http://localhost:%d/\n", HTTPListenPort)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", HTTPListenPort), nil))
}
