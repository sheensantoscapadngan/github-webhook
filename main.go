package main

import (
	"log"
	"net/http"

	"github-webhook/app"
	branchhandler "github-webhook/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	app := app.NewApp()

	app.Router.Post("/branch-tag-creation", func(w http.ResponseWriter, r *http.Request) {
		branchhandler.HandleBranchTagCreation(app, w, r)
	})

	app.Serve()
}