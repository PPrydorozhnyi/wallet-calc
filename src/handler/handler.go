package handler

import (
	"encoding/json"
	"fmt"
	"github.com/PPrydorozhnyi/wallet/db"
	"github.com/PPrydorozhnyi/wallet/model"
	"io"
	"net/http"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostsHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleGetPosts(w, r)
	case "POST":
		handleCreatePost(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func AccountHandle(w http.ResponseWriter, r *http.Request) {
	idString := strings.Split(r.URL.Path, "/")[2]

	switch r.Method {
	case "GET":
		handleGetAccount(w, r, idString)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func handleGetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(model.Posts())
	if err != nil {
		http.Error(w, "Cannot serialize posts", http.StatusInternalServerError)
	}
}

func handleGetAccount(w http.ResponseWriter, r *http.Request, id string) {
	account, err := db.GetWallet(id)

	if err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	serializeAccount(w, account)
}

func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}

	var p model.Post
	err = json.Unmarshal(body, &p)

	if err != nil {
		http.Error(w, "Cannot parse body", http.StatusBadRequest)
		return
	}

	post, cErr := model.CreatePost(p)

	if cErr != nil {
		http.Error(w, cErr.Error(), http.StatusInternalServerError)
		return
	}

	serializePost(w, post)
}

func serializePost(w http.ResponseWriter, post *model.Post) {
	err := json.NewEncoder(w).Encode(post)
	if err != nil {
		http.Error(w, "Cannot serialize post", http.StatusInternalServerError)
	}
}

func serializeAccount(w http.ResponseWriter, post *model.Account) {
	err := json.NewEncoder(w).Encode(post)
	if err != nil {
		http.Error(w, "Cannot serialize account", http.StatusInternalServerError)
	}
}
