package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

var tpl *template.Template

var Store = sessions.NewCookieStore([]byte("session"))

func init() {
	tpl = template.Must(template.ParseGlob("template/*.html"))
}

type page struct {
	Status     bool
	Header1    interface{}
	IsLoggedin bool
	Valid      bool
}

var P = page{
	Status: false,
}

var userData = map[string]string{
	"email":    "akshay@gmail.com",
	"password": "akshay123",
}

func index(w http.ResponseWriter, r *http.Request) {

	ok := Middleware(w, r)
	if ok {
		P.Status = true
	}

	filename := "index.html"
	err := tpl.ExecuteTemplate(w, filename, P)

	if err != nil {
		fmt.Println("error while parsing file", err)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	ok := Middleware(w, r)

	if ok {
		http.Redirect(w, r, "/login-submit", http.StatusSeeOther)
		return
	}
	filename := "login.html"
	err := tpl.ExecuteTemplate(w, filename, P)

	if err != nil {
		fmt.Println("an error occurd", err)
	}

}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if r.Method == "GET" {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		ok := Middleware(w, r)
		if ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "there is an error parsing %v", err)
		return
	}
	emails := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	if userData["email"] == emails && userData["password"] == password && r.Method == "POST" {
		session, _ := Store.Get(r, "started")
		session.Values["id"] = "AKSHAY"
		P.Header1 = session.Values["id"]

		session.Save(r, w)

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {

		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return

	}

}

func logouthandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if P.Status == true {
		session, _ := Store.Get(r, "started")
		session.Options.MaxAge = -1
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		P.Status = false
		P.Header1 = ""
	} else if P.Status == false {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func Middleware(w http.ResponseWriter, r *http.Request) bool {

	session, _ := Store.Get(r, "started")
	if session.Values["id"] == nil {
		return false
	}

	P.Header1 = session.Values["id"]
	return true

}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/login-submit", loginHandler)
	http.HandleFunc("/logout", logouthandler)
	fmt.Println("sever running in port 8080")
	http.ListenAndServe("localhost:8080", nil)

}
