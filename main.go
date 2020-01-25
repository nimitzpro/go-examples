package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Post struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type DB = map[int]Post

var db DB

func (post *Post) ToJSON() []byte {
	js, _ := json.Marshal(post)
	return js
}

func init() {
	db = make(map[int]Post)
	db[1] = Post{
		Name: "init",
		Body: "body",
	}
	db[2] = Post{
		Name: "Init 2",
		Body: "body 2",
	}
}

func main() {
	fmt.Println("Hello World")
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.RealIP,
		middleware.Timeout(60*time.Second),
	)

	router.Get("/post/{id}", GetPost)
	http.ListenAndServe(":3000", router)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	id_str := chi.URLParam(r, "id")
	var (
		id  int
		err error
	)
	if id, err = strconv.Atoi(id_str); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if post, ok := db[id]; ok {
		w.Header().Add("Content-Type", "application/json")
		js := post.ToJSON()
		code, _ := w.Write(js)
		fmt.Println(code)
		return
	} else {
		http.Error(w, http.StatusText(404), 404)
	}
}
