package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, get a fucking PS4 %s!", r.URL.Path[1:])
}

func main() {
    port := "8080"
    fmt.Printf("Server started on port %s", port);
    http.HandleFunc("/", handler)
    http.ListenAndServe(":"+port, nil)
}
