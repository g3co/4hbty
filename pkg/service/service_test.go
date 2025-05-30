package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/valeriikabisov/rakia/pkg/models"
	"github.com/valeriikabisov/rakia/pkg/store"
)

func TestCreatePost(t *testing.T) {
	store := store.NewPostStore()
	service := NewService(store)

	tests := []struct {
		name       string
		payload    models.Post
		wantStatus int
	}{
		{
			name: "valid post",
			payload: models.Post{
				Title:   "Test Post",
				Content: "Test Content",
				Author:  "Test Author",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid post - missing title",
			payload: models.Post{
				Content: "Test Content",
				Author:  "Test Author",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			service.CreatePost(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreatePost() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetPost(t *testing.T) {
	store := store.NewPostStore()
	service := NewService(store)

	// Create a test post
	post := &models.Post{
		Title:   "Test Post",
		Content: "Test Content",
		Author:  "Test Author",
	}
	store.Create(post)

	tests := []struct {
		name       string
		id         int
		wantStatus int
	}{
		{
			name:       "existing post",
			id:         post.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing post",
			id:         999,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/posts/"+strconv.Itoa(tt.id), nil)
			w := httptest.NewRecorder()

			// Set the URL parameters
			vars := map[string]string{
				"id": strconv.Itoa(tt.id),
			}
			req = mux.SetURLVars(req, vars)

			service.GetPost(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetPost() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	store := store.NewPostStore()
	service := NewService(store)

	// Create a test post
	post := &models.Post{
		Title:   "Test Post",
		Content: "Test Content",
		Author:  "Test Author",
	}
	store.Create(post)

	tests := []struct {
		name       string
		id         int
		wantStatus int
	}{
		{
			name:       "existing post",
			id:         post.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "non-existing post",
			id:         999,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/posts/"+strconv.Itoa(tt.id), nil)
			w := httptest.NewRecorder()

			// Set the URL parameters
			vars := map[string]string{
				"id": strconv.Itoa(tt.id),
			}
			req = mux.SetURLVars(req, vars)

			service.DeletePost(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("DeletePost() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				// Verify post was actually deleted
				_, err := store.Get(tt.id)
				if err == nil {
					t.Error("DeletePost() post still exists after deletion")
				}
			}
		})
	}
}

func TestUpdatePost(t *testing.T) {
	store := store.NewPostStore()
	service := NewService(store)

	// Create initial post
	post := &models.Post{
		Title:   "Initial Title",
		Content: "Initial Content",
		Author:  "Initial Author",
	}
	store.Create(post)

	tests := []struct {
		name       string
		id         int
		updateJSON string
		wantStatus int
		wantTitle  string
	}{
		{
			name:       "valid update",
			id:         post.ID,
			updateJSON: `{"title": "Updated Title", "content": "Updated Content", "author": "Updated Author"}`,
			wantStatus: http.StatusOK,
			wantTitle:  "Updated Title",
		},
		{
			name:       "non-existent post",
			id:         999,
			updateJSON: `{"title": "Updated Title"}`,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid JSON",
			id:         post.ID,
			updateJSON: `{invalid json}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "partial update",
			id:         post.ID,
			updateJSON: `{"title": "New Title"}`,
			wantStatus: http.StatusOK,
			wantTitle:  "New Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", "/posts/"+strconv.Itoa(tt.id), strings.NewReader(tt.updateJSON))
			w := httptest.NewRecorder()

			// Set the URL parameters
			vars := map[string]string{
				"id": strconv.Itoa(tt.id),
			}
			req = mux.SetURLVars(req, vars)

			service.UpdatePost(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("UpdatePost() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				// Verify post was actually updated
				updatedPost, err := store.Get(tt.id)
				if err != nil {
					t.Error("UpdatePost() failed to get updated post")
				}
				if updatedPost.Title != tt.wantTitle {
					t.Errorf("UpdatePost() title = %v, want %v", updatedPost.Title, tt.wantTitle)
				}
			}
		})
	}
}
