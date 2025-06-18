-- name: GetProjects :many
SELECT *
FROM projects
ORDER BY created_at ASC;
-- name: GetProjectById :one
SELECT *
FROM projects
WHERE id = $1;
-- name: GetProjectsByUserId :many
SELECT *
FROM projects
WHERE user_id = $1
ORDER BY created_at ASC;
-- name: CreateProject :one
INSERT INTO projects(
        id,
        created_at,
        updated_at,
        name,
        user_id
    )
VALUES(
        gen_random_uuid(),
        NOW(),
        NOW(),
        $1,
        $2
    )
RETURNING *;
-- name: DeleteProjectById :exec
DELETE from projects
WHERE id = $1;