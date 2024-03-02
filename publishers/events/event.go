package eventspublisher

import "github.com/jackc/pgx/v5/pgxpool"

type UnpublishedEventSlice interface {
	ParseString() string
	MarkEventsAsPublished(*pgxpool.Pool) error
}