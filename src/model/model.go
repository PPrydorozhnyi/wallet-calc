package model

import (
	"github.com/google/uuid"
	"sync"
)

type Post struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Body string    `json:"body,omitempty"`
}

var (
	posts   = make(map[uuid.UUID]*Post)
	postsMu sync.Mutex
)

func Posts() []*Post {
	ps := make([]*Post, 0, len(posts))

	for _, p := range posts {
		ps = append(ps, p)
	}

	return ps
}

func FindPost(id uuid.UUID) (*Post, bool) {
	post, ok := posts[id]
	return post, ok
}

func CreatePost(p Post) (*Post, error) {
	postsMu.Lock()
	defer postsMu.Unlock()

	id, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	post := &Post{ID: id, Body: p.Body}

	posts[id] = post

	return post, nil
}

func DeletePost(id uuid.UUID) {
	postsMu.Lock()
	defer postsMu.Unlock()

	delete(posts, id)
}
