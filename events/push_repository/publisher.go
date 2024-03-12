package pushrepositoryevt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
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

	return entries, nil
 }

 func (u UnpublishedRepositoryPushSlice) ParseString() string {
	parsedString := ""
	loc, _ := time.LoadLocation("Asia/Manila")

	commitMessageMaxLen, err := strconv.ParseInt(os.Getenv("PUSH_REPOSITORY_COMMIT_MESSAGE_MAX_LENGTH"), 10, 0)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	modifiedFilesMaxEntries, err := strconv.ParseInt(os.Getenv("PUSH_REPOSITORY_MODIFIED_FILES_MAX_ENTRIES"), 10, 0)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	prRegex, err := regexp.Compile(`^Merge pull request.*`)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	mergeRegex, err := regexp.Compile(`Merge.*branch.*into.*`)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	
	for _, entry := range u {
		modifiedFilesMap := make(map[string]bool)

		commitString := ""
		// filter out commits not made by pusher
		var ownCommits []PushCommit
		for _, commit := range entry.Commits {
			if commit.Committer.Username == entry.Pusher {
				ownCommits = append(ownCommits, commit)
			}
		}

		if len(ownCommits) == 0 {
			continue
		}

		for commitIndex, commit := range ownCommits {
			commitTimeInUTC, err := time.Parse(time.RFC3339, commit.Timestamp)
			if err != nil {
				log.Println("Error parsing date.", err.Error())
				break;
			}
		
			commitTimeInPH := commitTimeInUTC.In(loc)

			// do not include modified files when push is a merge branch because it's unnecessary
			if !(prRegex.MatchString(commit.Message) || mergeRegex.MatchString(commit.Message)){
				for _, file := range commit.ModifiedFiles {
					modifiedFilesMap[file] = true
				}
			}

			message := commit.Message
			messageRune := []rune(message)
			if len(messageRune) > int(commitMessageMaxLen) {
				message = string(messageRune[:commitMessageMaxLen]) + "..."
			}

			commitString += fmt.Sprintf("%d.%s commited with message:%s on %s\n",
				commitIndex+1,
				commit.Committer.Username,
				message,
				commitTimeInPH,
			)
		}

		var modifiedFiles []string
		for key := range modifiedFilesMap {
			modifiedFiles = append(modifiedFiles, key)
			if len(modifiedFiles) > int(modifiedFilesMaxEntries) {
				break
			}
		}

		parsedString += fmt.Sprintf("PUSH REPOSITORY EVENT\nReference:%s\nRepository:%s\nAuthor:%s\nDate/Time:%s\nCommits:%sModified files:%s\n\n",
			entry.Reference,
			entry.RepositoryName,
			entry.Pusher,
			entry.Date.Format(time.RFC850),
			commitString,
			strings.Join(modifiedFiles, ","),
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
