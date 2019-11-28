package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
	p := &Page{Title: title, Body: []byte("Welcome to Keto Cookout")}
	renderTemplate(w, "home", p)
}

func signupHandler(w http.ResponseWriter, r *http.Request, title string) {

	if r.Method == "POST" {
		fmt.Println("Posty post")
		fmt.Fprintln(w, "Yay you submitted the form")
		return
	}

	p := &Page{Title: title, Body: []byte("This is the body")}
	renderTemplate(w, "signup", p)
}

var templates = template.Must(template.ParseFiles("home.html", "signup.html"))

func renderTemplate(w http.ResponseWriter, template string, p *Page) {
	err := templates.ExecuteTemplate(w, template+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(|signup)/?$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			fmt.Println("Shiet")
			http.NotFound(w, r)
			return
		}

		fn(w, r, "some title")
	}
}

func main() {
	http.HandleFunc("/", makeHandler(homeHandler))
	http.HandleFunc("/signup/", makeHandler(signupHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
