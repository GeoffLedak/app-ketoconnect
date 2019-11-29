package main

import (
	"database/sql"
	_ "database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your-password"
	dbname   = "ketoconnect"
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

		if err := r.ParseForm(); err != nil {
			panic("error parsing signup form. " + err.Error())
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		fmt.Fprintf(w, "Yay you submitted the form\n\n")

		fmt.Fprintf(w, "Username = %s\n", username)
		fmt.Fprintf(w, "Password = %s\n", password)

		sqlStatement := `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
		RETURNING id`
		id := 0
		err := db.QueryRow(sqlStatement, username, password).Scan(&id)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "New record ID is: %d\n", id)

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

var db *sql.DB

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", makeHandler(homeHandler))
	http.HandleFunc("/signup/", makeHandler(signupHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
