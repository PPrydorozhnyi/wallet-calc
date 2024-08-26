package handler

import (
	"encoding/json"
	"fmt"
	"github.com/PPrydorozhnyi/wallet/model"
	"github.com/google/uuid"
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

func PostHandle(w http.ResponseWriter, r *http.Request) {
	idString := strings.Split(r.URL.Path, "/")[2]
	id, err := uuid.Parse(idString)

	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		handleGetPost(w, r, id)
	case "DELETE":
		handleDeletePost(w, r, id)
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

func handleGetPost(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	post, ok := model.FindPost(id)

	if !ok {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	serializePost(w, post)
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

func handleDeletePost(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	model.DeletePost(id)
	w.WriteHeader(http.StatusNoContent)
}
