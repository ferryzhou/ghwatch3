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
bqrun --jobs_path=bqjobs/*.yml  --project=950350008903 \
 --dataset=ghwatch3 --bucket=ghwatch3 --dst_dir=results
```

## Process Data

map shortPath to int and vice-versa
recs[i] is a slice of pairs, []<int, int>

```
proc --in_dir=results --out_dir=processed
```


## Serve the Processed Data With Restful API

1. Install postgrest

1. Serve repos data

1. Serves recommendation data


## Frontend



## Required Data
1. repos.csv
	name, owner, created_at, watchers, language, discription, ...
1. recs.csv
    shortPath, shortPath, score
