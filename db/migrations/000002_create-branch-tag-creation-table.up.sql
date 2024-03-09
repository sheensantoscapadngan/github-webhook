CREATE TABLE github.branch_tag_creation (
    branch_tag_creation_id serial,
    repository_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
    date timestamp without time zone NOT NULL,
    author character varying(100) COLLATE pg_catalog."default" NOT NULL,
    branch_tag_name text COLLATE pg_catalog."default" NOT NULL,
    formatted_date character varying(50) COLLATE pg_catalog."default" NOT NULL,
    is_published boolean NOT NULL DEFAULT false,
    CONSTRAINT branch_tag_creation_pkey PRIMARY KEY (branch_tag_creation_id)
);