-- Create a new table with the same structure as repo_urls
CREATE TABLE repo_urls_v2 AS TABLE repo_urls WITH NO DATA;

-- -- Copy 250,000 rows from repo_urls to repo_urls_v2 and set status to 'pending'
INSERT INTO repo_urls_v2 (url, status, created_at, updated_at)
SELECT url, 'pending'::repo_status, created_at, updated_at
FROM repo_urls
LIMIT 300000;