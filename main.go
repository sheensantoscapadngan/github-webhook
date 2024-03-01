package main

import (
	"log"
	"net/http"
	"os"

	"github-webhook/app"
	branchhandler "github-webhook/handlers"

	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		log.Println("Successfully loaded envs")
	}
	
	app := app.NewApp()

	app.Router.Post("/branch-tag-creation", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling branch/tag creation event")
		branchhandler.HandleBranchTagCreation(app, w, r)
	})

	app.Serve()
}