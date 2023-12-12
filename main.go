package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Article struct {
	Id           uint16
	Title, Anons string
}

var note = []Article{}

var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT * FROM `notes`")
	if err != nil {
		panic(err)
	}

	note = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons)
		if err != nil {
			panic(err)
		}
		//fmt.Println(fmt.Sprintf("Post: %s with id: %d", post.Title, post.Id))
		note = append(note, post)
	}

	t.ExecuteTemplate(w, "index", note)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO `notes` (`title`, `anons`) VALUES ('%s','%s')", title, anons))
	if err != nil {
		panic(err)
	}
	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tmplFiles := []string{"templates/show.html", "templates/header.html", "templates/footer.html"}
	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		panic(err)
	}

	// connect to database!
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}
	// close connection!
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `notes` WHERE `id` = '%s'", vars["post_id"]))

	// res, err := db.Query("SELECT * FROM `notes` WHERE `id` = " + vars["post_id"])
	if err != nil {
		panic(err)
	}

	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons)
		if err != nil {
			panic(err)
		}
		//fmt.Println(fmt.Sprintf("Post: %s with id: %d", post.Title, post.Id))
		showPost = post
	}

	t.ExecuteTemplate(w, "show", showPost)

}

func handleFunc() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{post_id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
