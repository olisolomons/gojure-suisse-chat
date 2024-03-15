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
	users["oli1"] = User{username: "oli1", password: "124"}
	users["oli2"] = User{username: "oli2", password: "125"}
	users["oli3"] = User{username: "oli3", password: "126"}
	users["oli4"] = User{username: "oli4", password: "127"}
	users["oli5"] = User{username: "oli5", password: "128"}
	users["oli6"] = User{username: "oli6", password: "129"}
	users["oli7"] = User{username: "oli7", password: "130"}
	users["oli8"] = User{username: "oli8", password: "131"}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(indexHandler))
	mux.Handle("/login", http.HandlerFunc(loginHandler))

	server := http.Server{Addr: "localhost:8080", Handler: mux}
	log.Fatal(server.ListenAndServe())
}
