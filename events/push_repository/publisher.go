package pushrepositoryevt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryPush struct {
	Id int
	Reference string
	RepositoryName string
	Date time.Time
	Pusher string
	Commits []PushCommit
	IsPublished bool
}

type UnpublishedRepositoryPushSlice []RepositoryPush

func GetUnpublishedRepositoryPush(p *pgxpool.Pool, ctx context.Context) (UnpublishedRepositoryPushSlice, error) {
	rows, err := p.Query(ctx, `
	SELECT repository_push_id, reference, pusher_name, repository_name, commits, date, is_published
	FROM github.repository_push WHERE is_published = $1`, false)

	if err != nil {
		log.Println("Error fetching unpublished branch tag creation:", err.Error())
		return nil, err
	}
	defer rows.Close()
	
	entries, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (RepositoryPush, error) {
		var r RepositoryPush
		var commits []byte
		if err := row.Scan(&r.Id, &r.Reference, &r.Pusher, &r.RepositoryName, &commits, &r.Date, &r.IsPublished); err != nil {
			return r, err
		}

		if err := json.Unmarshal(commits, &r.Commits); err != nil {
			return r, err
		}

		return r, nil
	})

	if err != nil {
		log.Println("Error parsing row", err.Error());
		return nil, err
	}

	var repositoryPush UnpublishedRepositoryPushSlice
	for _, entry := range entries {
		if !entry.IsPublished {
			repositoryPush = append(repositoryPush, entry)
		}
	}

	return repositoryPush, nil
 }

 func (u UnpublishedRepositoryPushSlice) ParseString() string {
	parsedString := ""
	loc, _ := time.LoadLocation("Asia/Manila")
	
	for _, entry := range u {
		commitString := ""
		for _, commit := range entry.Commits {
			commitTimeInUTC, err := time.Parse(time.RFC3339, commit.Timestamp)
			if err != nil {
				log.Println("Error parsing date.", err.Error())
				break;
			}
		
			commitTimeInPH := commitTimeInUTC.In(loc)

			commitString += fmt.Sprintf(`%s/%s/%s commited with the message:%s on %s (PHILIPPINE TIME). This modified the following files: %s\n\n`,
				commit.Committer.Username,
				commit.Committer.Email,
				commit.Committer.Name,
				commit.Message,
				commitTimeInPH,
				strings.Join(commit.ModifiedFiles, ","),
			)
		}

		parsedString += fmt.Sprintf(`
			A Repository Push Event was made with a reference of %s to repository:%s. This was pushed by %s on %s (PHILIPPINE TIME).
			It contained the following commits:%s`,
			entry.Reference,
			entry.RepositoryName,
			entry.Pusher,
			entry.Date.Format(time.RFC850),
			commitString,
		)
 	}

	return parsedString
}

func (u UnpublishedRepositoryPushSlice) MarkEventsAsPublished(p *pgxpool.Pool, ctx context.Context){
	b := &pgx.Batch{}

	for _, entry := range u {
		b.Queue("UPDATE github.repository_push SET is_published = true WHERE repository_push_id = $1", entry.Id)
 	}

	err := p.SendBatch(ctx, b).Close()
	if err != nil {
		log.Println("FAILED TO MARK REPOSITORY PUSH EVENTS AS PUBLISHED:", err.Error())
		return
	}

	log.Printf("MARKED %d REPOSITORY PUSH EVENTS AS PUBLISHED.", len(u))
}
