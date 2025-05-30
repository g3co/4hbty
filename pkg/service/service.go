package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/valeriikabisov/rakia/pkg/models"
	"github.com/valeriikabisov/rakia/pkg/store"
)

// PostHandler handles HTTP requests for blog posts
type Service struct {
	store *store.PostStore
}

// NewPostHandler creates a new PostHandler instance
func NewService(store *store.PostStore) *Service {
	return &Service{store: store}
}

// MainHandler handles GET /
func (s *Service) MainHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Hello this is test exercise for RAKIA"})
}

// GetPosts handles GET /posts
func (s *Service) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts := s.store.GetAll()
	respondWithJSON(w, http.StatusOK, posts)
}

// GetPost handles GET /posts/{id}
func (s *Service) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	post, err := s.store.Get(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	respondWithJSON(w, http.StatusOK, post)
}

// CreatePost handles POST /posts
func (s *Service) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := post.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	if err := s.store.Create(&post); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating post")
		return
	}

	respondWithJSON(w, http.StatusCreated, post)
}

// UpdatePost handles PUT /posts/{id}
func (s *Service) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// Get existing post first
	existingPost, err := s.store.Get(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if post.Title != "" {
		existingPost.Title = post.Title
	}
	if post.Content != "" {
		existingPost.Content = post.Content
	}
	if post.Author != "" {
		existingPost.Author = post.Author
	}

	existingPost.UpdatedAt = time.Now()

	if err := s.store.Update(id, existingPost); err != nil {
		respondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	respondWithJSON(w, http.StatusOK, existingPost)
}

// DeletePost handles DELETE /posts/{id}
func (s *Service) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	if err := s.store.Delete(id); err != nil {
		respondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Post deleted successfully"})
}

// Helper functions for JSON responses
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
