query

repos

SELECT
  repository_name,
  repository_owner,
  repository_organization,
  repository_language,
  repository_url,
  repository_created_at,
  repository_description,
  MAX(repository_watchers) AS watchers,
  MAX(repository_pushed_at) pushed_at,
FROM
  [githubarchive:github.timeline]
WHERE
  repository_watchers > 1
GROUP EACH BY 1, 2, 3, 4, 5, 6, 7
ORDER BY
  watchers DESC;