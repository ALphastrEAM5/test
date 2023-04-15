package main

import (
    "fmt"
    "net/http"
)

type User struct {
	ID    uint   `gorm:"primary_key"`
	SRFID string `gorm:"unique"`
	Email string `gorm:"unique"`
	DOB   string
}

type HospitalUser struct {
	ID       uint   `gorm:"primary_key"`
	HCode    string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}

type HospitalData struct {
	ID        uint   `gorm:"primary_key"`
	HCode     string `gorm:"unique"`
	HName     string
	NormalBed int
	HICUBed   int
	ICUBed    int
	VBed      int
}

type BookingPatient struct {
	ID       uint   `gorm:"primary_key"`
	SRFID    string `gorm:"unique"`
	BedType  string
	HCode    int
	SpO2     int
	PName    string
	PPhone   string
	PAddress string
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })
    http.ListenAndServe(":8080", nil)
}
