package branchhandler

import (
	"encoding/json"
	"github-webhook/utils/event"
	"log"
	"net/http"
)

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

func HandleBranchTagCreation(w http.ResponseWriter, r *http.Request) {
	var payload BranchTagCreationPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	defer r.Body.Close()
	if err != nil {
		log.Println("Oops an error occured.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := event.WriteToMemory(payload); err != nil {
		log.Println("Write to memory failed!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("HANDLED BRANCH/TAG CREATION EVENT"))
}

func (b BranchTagCreationPayload) MarshalEvent() ([]byte, error){
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