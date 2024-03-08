package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

//go:embed index.tmpl
var index string
var indexTemplate *template.Template = template.Must(template.New("index").Parse(index))

//go:embed login.tmpl
var login string
var loginTemplate *template.Template = template.Must(template.New("login").Parse(login))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	indexTemplate.Execute(w, nil)
}

type User struct {
	username string
	password string // pls hack me
}

var users = make(map[string]User)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loginTemplate.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form params", http.StatusBadRequest)
			return
		}
		username := r.Form.Get("username")
		password := r.Form.Get("password")
        user, ok := users[username]
        if !ok {
            http.Error(w, "No such User", http.StatusNotFound)
            return
        }
        // NOT TODO
		//users[username] = User{username: username, password: password}

        // TODO
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(indexHandler))
	mux.Handle("/login", http.HandlerFunc(loginHandler))

	server := http.Server{Addr: "localhost:8080", Handler: mux}
	log.Fatal(server.ListenAndServe())
}
