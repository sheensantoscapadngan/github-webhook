package eventspublisher

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UnpublishedEventSlice interface {
	ParseString() string
	MarkEventsAsPublished(*pgxpool.Pool, context.Context) error
}