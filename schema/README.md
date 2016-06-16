createdb bqtest
psql bqtest < schema.sql
psql bqtest < import_csv.psql
psql bqtest
