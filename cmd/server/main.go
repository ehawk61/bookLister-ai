package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ehawk61/bookLister-ai/internal/handler"
	"github.com/ehawk61/bookLister-ai/internal/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthcheck", handler.Healthcheck)

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
