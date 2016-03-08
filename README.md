# ghwatch3
github recomendation version 3

idea: download data from bigquery github table, analyze the data 
 and visualize the data.

## Setup
1. Download pemission json file and save it to 
   ~/.config/gcloud/application_default_credentials.json, see 
   https://godoc.org/golang.org/x/oauth2/google#DefaultTokenSource
1. Enable storage JSON API

## Test Setup
1. Create project, dataset and storage bucket, set in bqclient_test.go.

## Operations
1. Download data
```
go run pipeline.go --jobs_path=jobs/*.yml  --project=950350008903 --dst_dir=results
```
1. Start Server
```
go run main.go -repos_path="../pipeline/results/repos.gzip" -recs_path_pattern="../pipeline/results/repo_repo_*.gzip"
```

## implementation
1. download gzip csv from bigquery. 
1. load data into memory
1. serve the data

## required data
1. repos.csv
	name, owner, created_at, watchers, language, discription, ...
1. repo_recs.csv
    shortPath, shortPath, score

# Developer Guide

examples:

https://github.com/google/google-api-go-client/tree/master/examples
https://cloud.google.com/storage/docs/json_api/v1/json-api-go-samples

