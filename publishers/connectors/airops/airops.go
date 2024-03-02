package airopsconnect

import (
	"context"
	eventspublisher "github-webhook/publishers/events"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Publish(s []eventspublisher.UnpublishedEventSlice, p *pgxpool.Pool, ctx context.Context) {
	var collatedString string
	for _, eventSlice := range s {
		collatedString += eventSlice.ParseString() + "\n"
	}

	// CALL AIROPS API
	// MARK EVENTS AS PUBLISHED
	for _, eventSlice := range s {
		eventSlice.MarkEventsAsPublished(p, ctx)
	}
}