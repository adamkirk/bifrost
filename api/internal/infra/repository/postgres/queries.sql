-- name: InsertEnvironment :exec
INSERT INTO environments (id, name)
VALUES ($1, $2);

-- name: GetEnvironmentByName :one
SELECT 
    *
FROM environments
WHERE
    name = $1;