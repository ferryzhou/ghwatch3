
set -e

PORT_REC=3001
PORT_REPO=3000

function get_bq_data() {
  ./bqrun --jobs_path=bqjobs/c1000/*.yml  --project=950350008903 \
   --dataset=ghwatch3 --bucket=ghwatch3 --dst_dir=results
}

function process_data() {
  go run cmd/procrun/procrun.go --in_dir=results --out_dir=processed
}

function serve_recs() {
  go run cmd/procserv/procserv.go --in_gob_path=processed/recs.gob --port=$PORT_REC
}

dbname=bqtest
schema_prefix="ghwatch3"
function csv2db() {
  echo "csv2db $2....."
  schema_tmpl="`cat schema/schema.sql`"
  t=$(date +"%Y%m%d%H%M%S")
  schema_name=ghwatch3_$t
  schema="${schema_tmpl//ghwatch3_/$schema_name}"
  echo $schema
  list_schemas
  echo $schema | psql $dbname
  echo "import repos csv ...."
  #gzip -dc results/repos_c1000_full.csv.gz > /tmp/repos.csv
  cp -f processed/repos.csv /tmp/repos.csv
  psql $dbname -c "COPY \"$schema_name\".repos (id,short_path,name,owner,org,lang,url,created_at,description,homepage,watchers,pushed_at) FROM '/tmp/repos.csv' DELIMITER ',' CSV HEADER;"
  echo $(psql $dbname -c "SELECT COUNT(*) FROM \"$schema_name\".repos;")
  echo "csv2db finished ...."
}

function list_schemas() {
  echo $(psql $dbname -c "select schema_name from information_schema.schemata where schema_name like '$schema_prefix%' order by schema_name desc")
}

function latest_schema() {
  echo $(psql $dbname -c "select schema_name from information_schema.schemata where schema_name like '$schema_prefix%' order by schema_name desc" | sed -n '3 p')
}

function serve_repos() {
  schema_name=$(latest_schema)
  echo $schema_name
  postgrest "host=localhost dbname=$dbname" -a $(whoami) -s $schema_name -p $PORT_REPO
}

function serve_frontend() {
  http-server web
}

function test_api() {
  echo 'test recommendation api ......'
  curl http://localhost:$PORT_REC/rec/twbs/bootstrap
  curl http://localhost:$PORT_REC/recn/twbs/bootstrap
  echo 'test repo api ........'
  # cannot use curl
  #curl http://localhost:3000/repos?short_path=eq.twitter/bootstrap
}

$@
