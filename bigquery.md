# Data Processing with BigQuery

Show Options, save results to table, allow big results.

Step 1. Filter Repos

``` sql
SELECT * FROM (
  SELECT repository_url, COUNT(1) AS c
  FROM [githubarchive:github.timeline]
  WHERE type="WatchEvent"
  GROUP BY 1 )
WHERE c > 100
ORDER BY 2 DESC;
```

Step 2. Filter Repo-Actor-Createdat

``` sql
SELECT
  a.repository_url AS repository_url,
  actor,
  created_at
FROM [githubarchive:github.timeline] a
RIGHT JOIN EACH [github.repos_100] b
ON a.repository_url = b.repository_url
WHERE type = "WatchEvent"
GROUP EACH BY 1, 2, 3 ;
```

Step 3. Compute Repo-Repo-Count

``` sql
SELECT * FROM (
  SELECT
    a.repository_url AS repo1,
    b.repository_url AS repo2,
    COUNT(1) AS c
  FROM [github.repo_actor_createdat_100] a
  RIGHT JOIN EACH [github.repo_actor_createdat_100] b
  ON a.actor = b.actor
  GROUP EACH BY 1, 2 )
WHERE c > 10;
```
