CREATE TABLE github.repository_push (
    repository_push_id serial NOT NULL,
    reference text NOT NULL,
    pusher_name text NOT NULL,
    repository_name text NOT NULL,
    commits jsonb [] NOT NULL,
    date timestamp without time zone NOT NULL,
    PRIMARY KEY (repository_push_id)
);