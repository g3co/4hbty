package rest

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/valeriikabisov/rakia/pkg/service"
)

func NewHTTPROuter(s *service.Service) http.Handler {
	router := mux.NewRouter()
	// Define routes
	router.HandleFunc("/", s.MainHandler).Methods("GET")
	router.HandleFunc("/posts", s.GetPosts).Methods("GET")
	router.HandleFunc("/posts/{id}", s.GetPost).Methods("GET")
	router.HandleFunc("/posts", s.CreatePost).Methods("POST")
	router.HandleFunc("/posts/{id}", s.UpdatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", s.DeletePost).Methods("DELETE")

	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	headersOk := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"Accept",
		"Accept-Language",
		"Content-Type",
		"Content-Language",
		"Origin",
	})

	return handlers.CORS(originsOk, headersOk, methodsOk)(router)
}
