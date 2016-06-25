# ghwatch3
github recommendation version 3

idea: download data from bigquery github table, analyze the data
 and visualize the data.

## Analyze and Download Data From BigQuery

1. Use github.com/ferryzhou/gcutil

1. Download pemission json file and save it to
   ~/.config/gcloud/application_default_credentials.json, see
   https://godoc.org/golang.org/x/oauth2/google#DefaultTokenSource

1. Enable storage JSON API

1. Setup bigquery project, dataset and storage bucket

1. Download data
  ```
  ./run.sh get_bq_data
  ```

1. Result Data
  1. repos.csv
	  repo_url, name, owner, created_at, watchers, language, description, ...
  1. recs.csv
    repo1_url, repo2_url, count

## Process Data

raw data is large and we don't need them all. here we sequencing the url and
truncate recommendations data.

map shortPath to int and vice-versa
recs[i] is a slice of

```
./run.sh process_data
```

## Serve the Processed Data With Restful API

1. Prerequisite

```
npm install http-server -g
//install postgrest
```

1. Serve repos data: /repos?
  1. Load data to postgres
    ```
   ./run.sh csv2db
    ```
  1. Run postgrest
    ```
    ./run.sh serve_repos
    ```

1. Serves recommendation data
```
./run.sh serve_recs
```

1. Test api

```
./run.sh test_api
```

## Frontend
```
./run.sh serve_frontend
```

## References

1. http://postgrest.com/api/reading/
