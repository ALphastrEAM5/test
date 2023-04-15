package main

import (
    "fmt"
    "net/http"

)

type Config struct {
	GmailUser     string `json:"gmail-user"`
	GmailPassword string `json:"gmail-password"`
}

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

func dbMiddleware() gin.HandlerFunc {
	db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/covid?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{}, &HospitalUser{}, &HospitalData{}, &BookingPatient{})

	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
		db.Close()
	}
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })



    http.ListenAndServe(":8080", nil)
}
