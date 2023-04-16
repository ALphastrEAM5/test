package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", fs)
	log.Print("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./bootstrap.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	str := "This is a Golang example\nPlease, select a valid repository changing the LAUNCHER_TARGET_REPO value"
	fmt.Fprintf(w, str)
}