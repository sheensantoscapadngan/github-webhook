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

	app.Router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		event := r.Header.Get("X-GitHub-Event")
		log.Println("HANDLING GITHUB EVENT: ", event)
		switch event {
		case "create":
			branchhandler.HandleBranchTagCreation(app, w, r)
		}
	})

	app.Router.Post("/trigger", func(w http.ResponseWriter, r *http.Request) {
		parsedBranchTagCreation, branchTagCreationIds, err := branchhandler.ParseUnpublishedBranchTagCreation(
			app, 
			w,
			r.Context(),
		)
		if err != nil {
			log.Println(err.Error())
		}
		
		log.Println("PARSED BRANCH TAG CREATION IS", parsedBranchTagCreation, "WITH IDS", branchTagCreationIds)

	})

	app.Serve()
}