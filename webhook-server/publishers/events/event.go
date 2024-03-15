package eventspublisher

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UnpublishedEvent interface {
	ParseString() string
	MarkEventAsPublished(*pgxpool.Pool, context.Context)
}

type UnpublishedEventSlice interface {
	ParseString() string
	ParseStringByBatch() []string
	MarkEventsAsPublished(*pgxpool.Pool, context.Context)
	GetEventType() string
}
