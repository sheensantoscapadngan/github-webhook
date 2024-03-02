package airopsconnect

import (
	eventspublisher "github-webhook/publishers/events"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Publish(p *pgxpool.Pool, s []eventspublisher.UnpublishedEventSlice) {
	var collatedString string
	for _, eventSlice := range s {
		collatedString += eventSlice.ParseString() + "\n"
	}

	// CALL AIROPS API
	// MARK EVENTS AS PUBLISHED
	for _, eventSlice := range s {
		eventSlice.MarkEventsAsPublished(p)
	}
}