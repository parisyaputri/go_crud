package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Author struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"author_name"`
}

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}

	createAuthorsTable := `CREATE TABLE IF NOT EXISTS authors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL
	);`

	createPostsTable := `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		author_name TEXT NOT NULL
	);`

	_, err = db.Exec(createAuthorsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createPostsTable)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func createAuthor(w http.ResponseWriter, r *http.Request) {
	var author Author
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := initDB()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO authors (email) VALUES (?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(author.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	author.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(author)
}

func getAuthors(w http.ResponseWriter, r *http.Request) {
	db := initDB()
	defer db.Close()

	rows, err := db.Query("SELECT id, email FROM authors")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var author Author
		if err := rows.Scan(&author.ID, &author.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		authors = append(authors, author)
	}

	json.NewEncoder(w).Encode(authors)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := initDB()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO posts (title, description, author_name) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(post.Title, post.Description, post.AuthorName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	post.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	db := initDB()
	defer db.Close()

	rows , err := db.Query("SELECT id, title, description, author_name FROM posts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.AuthorName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	var updatedPost Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := initDB()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE posts SET title = ?, description = ?, author_name = ? WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updatedPost.Title, updatedPost.Description, updatedPost.AuthorName, updatedPost.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedPost)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	db := initDB()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	http.HandleFunc("/authors", createAuthor)
	http.HandleFunc("/authors/get", getAuthors)
	http.HandleFunc("/posts", createPost)
	http.HandleFunc("/posts/get", getPosts)
	http.HandleFunc("/posts/update", updatePost)
	http.HandleFunc("/posts/delete", deletePost)
	log.Fatal(http.ListenAndServe(":8080", nil))
}