-- name: InsertEnvironment :exec
INSERT INTO environments (id, name)
VALUES ($1, $2);
