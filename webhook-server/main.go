package main

import (
	"log"
	"net/http"
	"os"

	"github-webhook/app"
	branchtagcreationevt "github-webhook/events/branch_tag_creation"
	pushrepositoryevt "github-webhook/events/push_repository"
	publisher "github-webhook/publishers"

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

	app.Router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		event := r.Header.Get("X-GitHub-Event")
		log.Println("HANDLING GITHUB EVENT: ", event)
		switch event {
		case "create":
			branchtagcreationevt.HandleBranchTagCreation(app.Pool, w, r)
		case "push":
			pushrepositoryevt.HandlePushRepository(app.Pool, w, r)
		}
	})

	app.Router.Post("/trigger", func(w http.ResponseWriter, r *http.Request) {
		publisher.HandlePublishEvents(app, w, r)
	})

	app.Serve()
}