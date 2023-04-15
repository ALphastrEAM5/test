package main

import (
    "fmt"
    "net/http"
    "encoding/json"
    "net/smtp"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/go-gomail/gomail"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    })



    http.ListenAndServe(":8080", nil)
}
