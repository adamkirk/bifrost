-- name: InsertEnvironment :exec
INSERT INTO environments (id, name)
VALUES ($1, $2);

-- name: GetEnvironmentByName :one
SELECT
    *
FROM environments
WHERE
    name = $1;

-- name: ListEnvironments :many
SELECT *
FROM environments
ORDER BY name ASC
LIMIT $1
OFFSET $2;

-- name: CountEnvironments :one
SELECT COUNT(*) FROM environments;

-- name: UpsertEnvironment :one
INSERT INTO environments (id, name)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name
RETURNING *;