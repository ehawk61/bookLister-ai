package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ehawk61/bookLister-ai/internal/handler"
	"github.com/ehawk61/bookLister-ai/internal/middleware"
	"github.com/ehawk61/bookLister-ai/internal/repository"
)

func main() {
	bookRepo := repository.NewMemoryBookRepository()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthcheck", handler.Healthcheck)
	mux.HandleFunc("GET /api/books", handler.ListBooks(bookRepo))
	mux.HandleFunc("POST /api/books", handler.CreateBooks(bookRepo))
	mux.HandleFunc("GET /api/books/q", handler.SearchBooks(bookRepo))
	mux.HandleFunc("GET /api/books/{id}", handler.GetBook(bookRepo))
	mux.HandleFunc("PUT /api/books/{id}", handler.UpdateBook(bookRepo))
	mux.HandleFunc("DELETE /api/books/{id}", handler.DeleteBook(bookRepo))

	chain := middleware.CorrelationID(middleware.APIVersion(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("bookLister-ai listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, chain); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
