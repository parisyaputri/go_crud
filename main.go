package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type Post struct {
	ID     int
	Title  string
	Body   string
	Author string
}

var posts = []Post{
	{ID: 1, Title: "First Post", Body: "This is the content of the first post.", Author: "Author 1"},
	{ID: 2, Title: "Second Post", Body: "This is the content of the second post.", Author: "Author 2"},
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/add_user", addUser)
	http.HandleFunc("/posts", postsHandler)
	http.HandleFunc("/posts/", postHandler)
	http.HandleFunc("/contact", contactHandler)
	http.ListenAndServe(":8080", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/layout.html", "templates/"+tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = t.ExecuteTemplate(w, "layout", data)
}
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "welcome.html", nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact.html", nil)
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login.html", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "register.html", nil)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	renderTemplate(w, "add_user.html", nil)
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "posts.html", posts)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/posts/"):]
	id, err := strconv.Atoi(idStr) // Convert string to int
	if err != nil {
		http.NotFound(w, r)
		return
	}
	for _, post := range posts {
		if post.ID == id { // Compare integers
			renderTemplate(w, "post.html", post)
			return
		}
	}
	http.NotFound(w, r)
}
