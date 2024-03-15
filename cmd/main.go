package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("internal/pages/views/*.tmpl"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	templates.ExecuteTemplate(w, "index.tmpl", nil)
}

type User struct {
	username string
	password string // pls hack me
}

var users = make(map[string]User)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "login.tmpl", nil)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form params", http.StatusBadRequest)
			return
		}
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		user, userExists := users[username]

		if !userExists || user.password != password {
			templates.ExecuteTemplate(w, "login.tmpl", struct{ Error string }{Error: "Incorrect username or password."})
			return
		}

		cookie := http.Cookie{
			Name:     "session",
			Value:    "random thingy",
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		http.NotFound(w, r)
	}
}

func init() {
	users["oli"] = User{username: "oli", password: "123"}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(indexHandler))
	mux.Handle("/login", http.HandlerFunc(loginHandler))

	server := http.Server{Addr: "localhost:8080", Handler: mux}
	log.Fatal(server.ListenAndServe())
}
