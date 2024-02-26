package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	branchhandler "github-webhook/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()
	r.Get("/", func (w http.ResponseWriter, r *http.Request)  {
		w.Write([]byte("Hello world!"))
	})

	r.Post("/branch-tag-creation", branchhandler.HandleBranchTagCreation)

	http.ListenAndServe(":" + os.Getenv("PORT"), r)
}