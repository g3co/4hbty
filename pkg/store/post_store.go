package store

import (
	"errors"
	"sync"

	"github.com/g3co/4hbty/pkg/models"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

// PostStore handles the storage of blog posts
type PostStore struct {
	posts map[int]*models.Post
	mu    sync.RWMutex
}

// NewPostStore creates a new PostStore instance
func NewPostStore() *PostStore {
	return &PostStore{
		posts: make(map[int]*models.Post),
	}
}

// Create adds a new post to the store
func (s *PostStore) Create(post *models.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	post.ID = len(s.posts) + 1
	s.posts[post.ID] = post
	return nil
}

// Get retrieves a post by ID
func (s *PostStore) Get(id int) (*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, exists := s.posts[id]
	if !exists {
		return nil, ErrPostNotFound
	}
	return post, nil
}

// GetAll retrieves all posts
func (s *PostStore) GetAll() []*models.Post {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := make([]*models.Post, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, post)
	}

	return posts
}

// Update modifies an existing post
func (s *PostStore) Update(id int, post *models.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.posts[id]; !exists {
		return ErrPostNotFound
	}

	post.ID = id
	s.posts[id] = post
	return nil
}

// Delete removes a post from the store
func (s *PostStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.posts[id]; !exists {
		return ErrPostNotFound
	}

	delete(s.posts, id)
	return nil
}
