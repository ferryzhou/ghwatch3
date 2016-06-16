BEGIN;

create schema "ghwatch3_";

CREATE TABLE "ghwatch3_".repos
(
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

CREATE TABLE "ghwatch3_".recs
(
  repo1 varchar,
  repo2 varchar,
  score real,
  PRIMARY KEY(repo1, repo2)
);

COMMIT;
