package branchhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github-webhook/app"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
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

type RawBranchTagCreation struct {
	Id int
	RepositoryName string
	Date time.Time
	Author string
	BranchTagName string
	FormattedDate string
	IsPublished bool
}

func ParseUnpublishedBranchTagCreation(a *app.App, w http.ResponseWriter, ctx context.Context) (string, []int, error) {
	rows, err := a.Pool.Query(ctx, `
	SELECT branch_tag_creation_id, repository_name, date, author, branch_tag_name, formatted_date, is_published
	FROM github.branch_tag_creation WHERE is_published = $1`, false)

	if err != nil {
		log.Println("Error fetching unpublished branch tag creation:", err.Error())
		return "", nil, err
	}
	defer rows.Close()
	entries, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (RawBranchTagCreation, error) {
		var r RawBranchTagCreation
		err := row.Scan(&r.Id, &r.RepositoryName, &r.Date, &r.Author, &r.BranchTagName, &r.FormattedDate, &r.IsPublished)
		return r, err
	})

	if err != nil {
		log.Println("Error parsing row", err.Error());
		return "", nil, err
	}

	var parsedString string
	var ids []int
	for _, entry := range entries {
		if !entry.IsPublished {
			parsedString += fmt.Sprintf(`
				A BRANCH/TAG with the name of %s was made in repository:%s. This was pushed by %s on %s (PHILIPPINE TIME)`,
				entry.BranchTagName,
				entry.RepositoryName,
				entry.Author,
				entry.FormattedDate,
			)
			ids = append(ids, entry.Id)
		}
	}

	return parsedString, ids, nil
 }

func HandleBranchTagCreation(a *app.App, w http.ResponseWriter, r *http.Request) {
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

	tag, err := a.Pool.Exec(r.Context(), "INSERT INTO github.branch_tag_creation(repository_name,date,formatted_date,author,branch_tag_name) VALUES($1, $2, $3, $4, $5)",
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