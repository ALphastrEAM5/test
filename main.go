package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "net/smtp"
    "os"

)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })



    http.ListenAndServe(":8080", nil)
}
