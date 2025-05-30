package store

import (
	"encoding/json"
	"os"
	"time"

	"github.com/g3co/4hbty/pkg/models"
)

type SeedData struct {
	Posts []struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	} `json:"posts"`
}

// LoadSeedData loads initial blog posts from blog_data.json into the store
func (s *PostStore) LoadSeedData(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	var seedData SeedData
	if err := json.Unmarshal(data, &seedData); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range seedData.Posts {
		post := &models.Post{
			ID:        p.ID,
			Title:     p.Title,
			Content:   p.Content,
			Author:    p.Author,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		s.posts[post.ID] = post
	}

	return nil
}
