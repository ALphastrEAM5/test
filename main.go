package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })

    // Set GIN_MODE environment variable to "debug"
    os.Setenv("GIN_MODE", "debug")

    http.ListenAndServe(":8080", nil)
}
