package main

import (
	"encoding/json"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

var (
	localServer = true
	configPath  = "config.json"
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

func main() {
	app := gin.Default()
	app.Use(dbMiddleware())
	app.Use(gin.Logger())

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	app.Run(":8080")
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

func sendEmail(to, subject, body string) error {
	cfg := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(cfg); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.GmailUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 465, cfg.GmailUser, cfg.GmailPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func beforeRequest() {
    g.hospitalUser = nil
    if session["hospital_user_id"] != nil {
        user := hospitaluser.query.get(session["hospital_user_id"])
        g.hospitalUser = user
    }
}

func hospitalLoginRequired(f http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if g.hospitalUser == nil {
            fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++Redirecting to hospital login page")
            http.Redirect(w, r, "/hospitallogin?next="+r.URL.String(), http.StatusFound)
            return
        }
        f(w, r)
    }
}

func home(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "index.html", nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        srfid := r.FormValue("srf")
        email := r.FormValue("email")
        dob := r.FormValue("dob")

        encpassword := generatePasswordHash(dob)
        user := User.query.filter_by(srfid=srfid).first()
        emailUser := User.query.filter_by(email=email).first()
        if user != nil || emailUser != nil {
            flash("Email or srfid is already taken", "warning", w, r)
            renderTemplate(w, "usersignup.html", nil)
            return
        }

        newUser := User{
            Srfid: srfid,
            Email: email,
            Dob:   encpassword,
        }
        db.session.add(&newUser)
        db.session.commit()

        flash("SignUp Success Please Login Success", "success", w, r)
        renderTemplate(w, "userlogin.html", nil)
        return
    }
    renderTemplate(w, "/usersignup.html", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        srfid := r.FormValue("srf")
        dob := r.FormValue("dob")
        user := User.query.filter_by(srfid=srfid).first()

        if user != nil && checkPasswordHash(user.Dob, dob) {
            login_user(user)
            flash("Login Success", "info", w, r)
            renderTemplate(w, "index.html", nil)
            return
        } else {
            flash("Invalid Credentials", "danger", w, r)
            renderTemplate(w, "userlogin.html", nil)
            return
        }
    }
    renderTemplate(w, "/userlogin.html", nil)
}

func hospitalLogin(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        email := r.FormValue("email")
        password := r.FormValue("password")
        user := hospitaluser.query.filter_by(email=email).first()

        if user != nil && checkPasswordHash(user.Password, password) {
            session["hospital_user_id"] = user.ID
            flash("Login Success", "info", w, r)
            http.Redirect(w, r, "/addhospitalinfo", http.StatusFound)
            return
        } else {
            flash("Invalid Credentials", "danger", w, r)
            renderTemplate(w, "hospitallogin.html", nil)
            return
        }
    }
    renderTemplate(w, "/hospitallogin.html", nil)
}
func admin(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        username := r.FormValue("username")
        password := r.FormValue("password")
        if username == params["user"] && password == params["password"] {
            session, _ := store.Get(r, "session-name")
            session.Values["user"] = username
            session.Save(r, w)
            flash("Login Success", "info", w)
            http.Redirect(w, r, "/addHosUser.html", http.StatusSeeOther)
            return
        } else {
            flash("Invalid Credentials", "danger", w)
        }
    }
    renderTemplate("admin.html", w)
}

func logout(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")
    session.Values["user"] = ""
    session.Options.MaxAge = -1
    session.Save(r, w)
    flash("Logout Successful", "warning", w)
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func main() {
    // Define your routes and handlers here using the "http" package.
    // ...

    // Set GIN_MODE environment variable to "debug"
    os.Setenv("GIN_MODE", "debug")

    // Start the server
    http.ListenAndServe(":8080", nil) // Replace 8080 with the port number you want to use
}