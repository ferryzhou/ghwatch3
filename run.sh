
set -e

function get_bq_data() {
  ./bqrun --jobs_path=bqjobs/c1000/*.yml  --project=950350008903 \
   --dataset=ghwatch3 --bucket=ghwatch3 --dst_dir=results
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
  echo $(psql $dbname -c "select schema_name from information_schema.schemata where schema_name like '$schema_prefix%'")
  echo $schema | psql $dbname
  echo "import csv ...."
  gzip -dc results/repos_c1000_full.csv.gz > /tmp/repos.csv
  psql $dbname -c "COPY \"$schema_name\".repos (name,owner,org,lang,url,created_at,description,homepage,watchers,pushed_at) FROM '/tmp/repos.csv' DELIMITER ',' CSV HEADER;"
  echo $(psql $dbname -c "SELECT COUNT(*) FROM \"$schema_name\".repos;")
  for f in $( ls results/repo_repo_count_c1000_*.csv.gz); do
    echo "unzip $f"
    gzip -dc $f > /tmp/recs.csv
    echo "import /tmp/recs.csv"
    psql $dbname -c "COPY \"$schema_name\".recs FROM '/tmp/recs.csv' DELIMITER ',' CSV HEADER;"
  done
  echo $(psql $dbname -c "SELECT COUNT(*) FROM \"$schema_name\".recs;")
  echo "csv2db finished ...."
}

$@
