package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("internal/pages/views/*.tmpl"))

func withAuth (h func(http.ResponseWriter, *http.Request, User)) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user:= users[cookie.Value]
		h(w, r , user)
	}) 

}

func getUser (w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request, u User) {
	templates.ExecuteTemplate(w, "index.tmpl", struct {Username string}{Username: u.username})
}

type User struct {
	username string
	password string // pls hack me
	displayname string 
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
			Value:    username,
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

func accountHandler (w http.ResponseWriter, r *http.Request, u User) {
	var dname string
	if (u.displayname == "") {
		dname = u.username

	}else {
	  dname = u.displayname
	}

	templates.ExecuteTemplate(w, "account.tmpl", struct {DisplayName string}{DisplayName: dname})
}

func saveAccountDetails (w http.ResponseWriter, r *http.Request, u User){
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form params", http.StatusBadRequest)
			return
		}
		displayName := r.Form.Get("display_name")
		u.displayname = displayName
		fmt.Printf("%s", displayName)


}}

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
	mux.Handle("/", withAuth(indexHandler))
	mux.Handle("/login", http.HandlerFunc(loginHandler))
	mux.Handle("/account", withAuth(accountHandler))
	mux.Handle("/save-account-details", withAuth(accountHandler)) // do an api call not change url on address also this doesnt work right now

	server := http.Server{Addr: "localhost:8080", Handler: mux}
	log.Fatal(server.ListenAndServe())
}
