package branchhandler

import (
	"encoding/json"
	"github-webhook/app"
	"log"
	"net/http"
	"time"
)

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

func HandleBranchTagCreation(a *app.App, w http.ResponseWriter, r *http.Request) {
	var payload BranchTagCreationPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	defer r.Body.Close()
	if err != nil {
		log.Println("Oops an error occured.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	timeInUTC, err := time.Parse(time.RFC3339, payload.Repository.PushedAt)
	if err != nil {
		log.Println("Error parsing date.")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	loc, _ := time.LoadLocation("Asia/Manila")
	timeInPH := timeInUTC.In(loc)

	tag, err := a.Pool.Exec(r.Context(), "INSERT INTO github.branch_tag_creation(repository_name,date,formatted_date,author,branch_tag_name) VALUES($1, $2, $3, $4, $5)",
		payload.Repository.Name,
		timeInPH.Format(time.DateTime),
		timeInPH.Format(time.RFC850),
		payload.Sender.Name,
		payload.BranchOrTagName,
	)
	
	log.Println("Inserted", tag.RowsAffected(), "BRANCH/TAG CREATION event")
	if err != nil {
		log.Println("Oops, an error occured when inserting event...")
		http.Error(w, err.Error(), http.StatusInternalServerError);
		return
	}
	
	w.Write([]byte("HANDLED BRANCH/TAG CREATION EVENT"))
}