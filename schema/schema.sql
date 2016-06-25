BEGIN;

create schema "ghwatch3_";

CREATE TABLE "ghwatch3_".repos
(
  id integer NOT NULL,
  short_path varchar NOT NULL,
  url varchar NOT NULL,
  name varchar NOT NULL,
  description text,
  lang varchar,
  owner varchar,
  org varchar,
  homepage varchar,
  watchers integer,
  created_at date,
  pushed_at date
);

COMMIT;
