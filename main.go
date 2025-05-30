package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/g3co/4hbty/pkg/rest"
	"github.com/g3co/4hbty/pkg/service"
	"github.com/g3co/4hbty/pkg/store"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Initialize the store
	postStore := store.NewPostStore()

	// Load seed data
	if err := postStore.LoadSeedData("blog_data.json"); err != nil {
		log.Printf("Warning: Failed to load seed data: %v", err)
	}

	// Initialize the service
	service := service.NewService(postStore)

	eg, errCtx := errgroup.WithContext(ctx)

	// Create a new router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/posts", service.GetPosts).Methods("GET")
	router.HandleFunc("/posts/{id}", service.GetPost).Methods("GET")
	router.HandleFunc("/posts", service.CreatePost).Methods("POST")
	router.HandleFunc("/posts/{id}", service.UpdatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", service.DeletePost).Methods("DELETE")

	httpPort := ":8088"
	restSrv := &http.Server{
		Addr:    httpPort,
		Handler: rest.NewHTTPROuter(service),
	}

	eg.Go(func() error {
		log.Println("Server starting on", httpPort)
		return restSrv.ListenAndServe()
	})

	eg.Go(func() error {
		<-errCtx.Done()
		return restSrv.Shutdown(errCtx)
	})

	err := eg.Wait()
	switch {
	case errors.Is(err, context.Canceled),
		errors.Is(err, http.ErrServerClosed),
		err == nil:
		log.Println("gracefully quit server")
	default:
		log.Println("error occurred when wrapping the service", err)
	}
}
