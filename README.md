# ghwatch3
github recomendation version 3

idea: download data from bigquery github table, analyze the data 
 and visualize the data.

## Test Setup
1. Download pemission json file and save to ~/.gcp/default.pem, which is 
   used by test code.
2. Create project, dataset and storage bucket, set in bqclient_test.go.

## implementation
1. download gzip csv from bigquery. 
1. load data into memory
1. serve the data

## required data
1. repos.csv
	name, owner, created_at, watchers, language, discription, ...
1. repo_recs.csv
    shortPath, shortPath, score
