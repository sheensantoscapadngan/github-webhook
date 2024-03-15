package branchtagcreationevt

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

func HandleBranchTagCreation(p *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	var payload BranchTagCreationPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	defer r.Body.Close()
	if err != nil {
		log.Println("Oops an error occured.", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	timeInUTC, err := time.Parse(time.RFC3339, payload.Repository.PushedAt)
	if err != nil {
		log.Println("Error parsing date.", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	loc, _ := time.LoadLocation("Asia/Manila")
	timeInPH := timeInUTC.In(loc)

	tag, err := p.Exec(r.Context(), "INSERT INTO github.branch_tag_creation(repository_name,date,formatted_date,author,branch_tag_name) VALUES($1, $2, $3, $4, $5)",
		payload.Repository.Name,
		timeInPH.Format(time.DateTime),
		timeInPH.Format(time.RFC850),
		payload.Sender.Name,
		payload.BranchOrTagName,
	)
	
	log.Println("Inserted", tag.RowsAffected(), "BRANCH/TAG CREATION event")
	if err != nil {
		log.Println("Oops, an error occured when inserting event...", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError);
		return
	}
	
	w.Write([]byte("HANDLED BRANCH/TAG CREATION EVENT"))
}