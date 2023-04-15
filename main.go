package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
)

var (
    // Replace with your own values
    cookieStore = sessions.NewCookieStore([]byte("super-secret-key"))
    username    = "user"
    password    = "password"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", indexHandler)
    r.HandleFunc("/login", loginHandler).Methods("POST")
    r.HandleFunc("/logout", logoutHandler).Methods("POST")

    http.ListenAndServe(":8080", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    session, err := cookieStore.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    fmt.Fprintln(w, "You are logged in.")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    password := r.FormValue("password")

    if username == "" || password == "" {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        return
    }

    if username == username && password == password {
        session, err := cookieStore.Get(r, "session-name")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        session.Values["authenticated"] = true
        err = session.Save(r, w)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    http.Error(w, "Invalid username or password", http.StatusBadRequest)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
    session, err := cookieStore.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    session.Values["authenticated"] = false
    err = session.Save(r, w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/login", http.StatusSeeOther)
}
