package pushrepositoryevt

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PushRepositoryPayload struct {
	Reference string `json:"ref"`
	IsCreated bool `json:"created"`
	IsDeleted bool `json:"deleted"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	Repository struct {
		Name string `json:"full_name"`
		PushedAt int64 `json:"pushed_at"`
	} `json:"repository"`
	Commits []struct {
		Message string `json:"message"`
		Timestamp string `json:"timestamp"`
		Committer struct {
			Name string `json:"name"`
			Email string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		ModifiedFiles []string `json:"modified"`
	} `json:"commits"`
}

func HandlePushRepository(p *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	var payload PushRepositoryPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	defer r.Body.Close()
	if err != nil {
		log.Println("Oops an error occured.", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loc, _ := time.LoadLocation("Asia/Manila")
	localPushedTime := time.Unix(payload.Repository.PushedAt, 0).In(loc)
	commitsInBytes, err := json.Marshal(payload.Commits)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tag, err := p.Exec(r.Context(), `INSERT INTO github.repository_push(
			reference,
			pusher_name,
			repository_name,
			commits,
			date
		) VALUES($1, $2, $3, $4, $5)`,
		payload.Reference,
		payload.Pusher.Name,
		payload.Repository.Name,
		commitsInBytes,
		localPushedTime.Format(time.DateTime),
	)
	
	log.Println("Inserted", tag.RowsAffected(), "PUSH REPOSITORY event")
	if err != nil {
		log.Println("Oops, an error occured when inserting event...", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError);
		return
	}
	
	w.Write([]byte("HANDLED PUSH REPOSITORY EVENT"))
}