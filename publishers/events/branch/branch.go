package branchpublisher

import (
	"context"
	"fmt"
	"github-webhook/app"
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

type UnpublishedBranchTagCreationSlice []RawBranchTagCreation

func GetUnpublishedBranchTagCreation(a *app.App, ctx context.Context) (UnpublishedBranchTagCreationSlice, error) {
	rows, err := a.Pool.Query(ctx, `
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

	var branchTagCreations UnpublishedBranchTagCreationSlice
	var ids []int
	for _, entry := range entries {
		if !entry.IsPublished {
			branchTagCreations = append(branchTagCreations, entry)
			ids = append(ids, entry.Id)
		}
	}

	return branchTagCreations, nil
 }

 func (u UnpublishedBranchTagCreationSlice) ParseString() string {
	var parsedString string

	for _, entry := range u {
		parsedString += fmt.Sprintf(`
			A BRANCH/TAG with the name of %s was made in repository:%s. This was pushed by %s on %s (PHILIPPINE TIME)`,
			entry.BranchTagName,
			entry.RepositoryName,
			entry.Author,
			entry.FormattedDate,
		)
 	}

	return parsedString
}

func (u UnpublishedBranchTagCreationSlice) MarkEventsAsPublished(p *pgxpool.Pool, ctx context.Context) error {
	b := &pgx.Batch{}

	for _, entry := range u {
		b.Queue("UPDATE")
 	}

	p.SendBatch(ctx, b)

	return nil
}
