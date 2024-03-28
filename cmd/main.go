package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
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
	templates.ExecuteTemplate(w, "account.tmpl", struct {DisplayName string}{DisplayName: u.displayname})
}

func chatHandler (w http.ResponseWriter, r *http.Request, u User) {
	templates.ExecuteTemplate(w, "chat.tmpl", struct {DisplayName string}{DisplayName: u.displayname})
}
func saveAccountDetails (w http.ResponseWriter, r *http.Request, u User) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form params", http.StatusBadRequest)
			return
		}
		displayName := r.Form.Get("display_name")
		if !strings.Contains(displayName, "fraser") {
			u.displayname = displayName
			fmt.Printf("%s", displayName)
			users[u.username] = u
			w.Write([]byte("<div><p>New Display Name " + displayName + "</p><script>alert('hello!')</script></div>"))
		} else {
			w.Write([]byte("<p>You cannot include that word!</p>"))
		}
	}
}

func createUser(username, password string) User {
	return User{username: username, password: password, displayname: username}
}

func init() {
	users["oli1"] = createUser("oli1", "124")
	users["oli2"] = createUser("oli2", "356")
	users["oli3"] = createUser("oli3", "789")
	users["oli4"] = createUser("oli4", "918")
	users["oli5"] = createUser("oli5", "184")
}

var wsConnections []*websocket.Conn

type Message struct {
	ChatMessage string `json:"chat_message"`
}

func EchoServer(ws *websocket.Conn) {
	defer ws.Close()
	cookie, err := ws.Request().Cookie("session")
	if err == http.ErrNoCookie {
		ws.Write([]byte("No valid session Cookie"))
		return
	}
	username := users[cookie.Value].displayname
	fmt.Println("username: " + username)
	wsConnections = append(wsConnections, ws)
	var m Message
	dec := json.NewDecoder(ws)
	for {
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else {
			for _, wsi := range wsConnections {
				wsi.Write([]byte("<div id=\"chat_room\" hx-swap-oob=\"beforeend\">" + username + ": " + m.ChatMessage + "<br /></div>"))
			}
		}
	}
}

func main() {
	fmt.Println("Server Started on port 8080")

	mux := http.NewServeMux()
	mux.Handle("/", withAuth(indexHandler))
	mux.Handle("/login", http.HandlerFunc(loginHandler))
	mux.Handle("/account", withAuth(accountHandler))
	mux.Handle("/save-account-details", withAuth(saveAccountDetails)) // do an api call not change url on address also this doesnt work right now

	mux.Handle("/ws", websocket.Handler(EchoServer))
	mux.Handle("/chat", withAuth(chatHandler))

	server := http.Server{Addr: "localhost:8080", Handler: mux}
	log.Fatal(server.ListenAndServe())
}
