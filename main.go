package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type GithubEvent interface {
	marshalEvent()([]byte, error)
}

// should this really be camel-case? 
type BranchTagCreationEvent struct {
	Action string `json:"action"`
	Repository string `json:"repository"`
	BranchOrTagName string `json:"branchOrTagName"`
	Author string `json:"author"`
	Date string `json:"date"`
}

type BranchTagCreationPayload struct {
	Repository struct {
		Name string `json:"full_name"`
		PushedAt string `json:"pushed_at"`
	} `json:"repository"`
	Sender struct {
		Name string `json:"login"`
	} `json:"sender"`
	BranchOrTagName string `json:"ref"`
}

func (b BranchTagCreationPayload) marshalEvent() ([]byte, error){
	event := BranchTagCreationEvent{
		Action: "BRANCH/TAG CREATION",
		Repository: b.Repository.Name,
		BranchOrTagName: b.BranchOrTagName,
		Author: b.Sender.Name,
		Date: b.Repository.PushedAt,
	}
	data, err := json.Marshal(event)

	if err != nil {
		log.Print("Oops something went wrong")
		return nil, err
	}

	return data, nil

}

func WriteToMemory(ge GithubEvent) (error) {
	data, err := ge.marshalEvent()
	if err != nil {
		return err
	}	

	resp, err := http.Post(os.Getenv("MEMORY_URL"), "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	return nil
}

func handleBranchTagCreation(w http.ResponseWriter, r *http.Request) {
	var payload BranchTagCreationPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	defer r.Body.Close()
	if err != nil {
		log.Println("Oops an error occured.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := WriteToMemory(payload); err != nil {
		log.Println("Write to memory failed!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("HANDLED BRANCH/TAG CREATION EVENT"))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := chi.NewRouter()
	r.Get("/", func (w http.ResponseWriter, r *http.Request)  {
		w.Write([]byte("Hello world!"))
	})

	r.Post("/branch-tag-creation", handleBranchTagCreation)

	http.ListenAndServe(":" + os.Getenv("PORT"), r)
}