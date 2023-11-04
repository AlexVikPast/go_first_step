package main

import (
	"fmt"
	"net/http"
	"html/template"
	"database/sql"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type Article struct {
	Id uint16
	Title, Anons, FullText string
}

var articles = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	psqlInfo := "dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
			panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM articles")

	if err != nil {
		panic(err)
	}

	articles = []Article{}
	for res.Next() {
		var article Article 
		err = res.Scan(&article.Id, &article.Title, &article.Anons, &article.FullText)
		if err != nil { panic(err) }

		fmt.Println(fmt.Sprintf("Post: %s with id: %d", article.Title, article.Id))

		articles = append(articles, article)
	}

	t.ExecuteTemplate(w, "index", articles)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func showPost(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	vars := mux.Vars(r)

	psqlInfo := "dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
			panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM articles WHERE id = %s", vars["id"]))

	if err != nil {
		panic(err)
	}

	showPost = Article{}
	for res.Next() {
		var article Article 
		err = res.Scan(&article.Id, &article.Title, &article.Anons, &article.FullText)

		if err != nil { 
			panic(err) 
		}

		fmt.Println(fmt.Sprintf("Post: %s with id: %d", article.Title, article.Id))
		showPost = article
	}

	t.ExecuteTemplate(w, "show", showPost)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	psqlInfo := "dbname=golang sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
			panic(err)
	}
	defer db.Close()
	
	insert, err := db.Query(fmt.Sprintf("INSERT INTO articles (title, anons, full_text) VALUES ('%s', '%s', '%s')", title, anons, full_text))
	if err != nil {
		panic(err)
	}
	defer insert.Close()

	http.Redirect(w, r, "/", 301)
}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	// http.HandleFunc("/save_article", save_article)
	rtr.HandleFunc("/post/{id:[0-9]+}", showPost).Methods("GET")

	http.Handle("/", rtr)
	http.ListenAndServe(":8888", nil)
}

func main() {
	handleFunc()
}