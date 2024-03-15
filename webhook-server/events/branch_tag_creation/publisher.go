package branchtagcreationevt

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RawBranchTagCreation struct {
	Id int
	RepositoryName string
	Date time.Time
	Author string
	BranchTagName string
	FormattedDate string
	IsPublished bool
}

type RawBranchTagCreationSlice []RawBranchTagCreation

const EVENT_TYPE = "BRANCH_TAG_CREATION"

func GetUnpublishedBranchTagCreation(p *pgxpool.Pool, ctx context.Context) (RawBranchTagCreationSlice, error) {
	rows, err := p.Query(ctx, `
	SELECT branch_tag_creation_id, repository_name, date, author, branch_tag_name, formatted_date, is_published
	FROM github.branch_tag_creation WHERE is_published = $1`, false)

	if err != nil {
		log.Println("Error fetching unpublished branch tag creation:", err.Error())
		return nil, err
	}
	defer rows.Close()
	
	entries, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (RawBranchTagCreation, error) {
		var r RawBranchTagCreation
		err := row.Scan(&r.Id, &r.RepositoryName, &r.Date, &r.Author, &r.BranchTagName, &r.FormattedDate, &r.IsPublished)
		return r, err
	})

	if err != nil {
		log.Println("Error parsing row", err.Error());
		return nil, err
	}

	var branchTagCreations RawBranchTagCreationSlice
	for _, entry := range entries {
		if !entry.IsPublished {
			branchTagCreations = append(branchTagCreations, entry)
		}
	}

	return branchTagCreations, nil
 }

 func (r RawBranchTagCreation) ParseString() string {
	return fmt.Sprintf(`
		A BRANCH/TAG with the name of %s was made in repository:%s. This was pushed by %s on %s (PHILIPPINE TIME)`,
		r.BranchTagName,
		r.RepositoryName,
		r.Author,
		r.FormattedDate,
	)
}


func (rs RawBranchTagCreationSlice) ParseString() string {
	parsedString := ""

	for _, r := range rs {
		parsedString += r.ParseString()
 	}

	return parsedString
}


func (rs RawBranchTagCreationSlice) GetEventType() string {
	return EVENT_TYPE
}

func (rbs RawBranchTagCreationSlice) ParseStringByBatch() []string {
	var parsedEntries []string
	for _, entry := range rbs {
		parsedEntries = append(parsedEntries, entry.ParseString())
	} 

	return parsedEntries
}


func (rbs RawBranchTagCreationSlice) MarkEventsAsPublished(p *pgxpool.Pool, ctx context.Context){
	b := &pgx.Batch{}

	for _, entry := range rbs {
		b.Queue("UPDATE github.branch_tag_creation SET is_published = true WHERE branch_tag_creation_id = $1", entry.Id)
 	}

	err := p.SendBatch(ctx, b).Close()
	if err != nil {
		log.Println("FAILED TO MARK BRANCH/TAG CREATION EVENTS AS PUBLISHED:", err.Error())
		return
	}

	log.Printf("MARKED %d BRANCH/TAG CREATION EVENTS AS PUBLISHED.", len(rbs))
}
