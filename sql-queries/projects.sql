-- name: GetProjects :many
SELECT * FROM projects;

-- name: CreateProject :one
INSERT into projects (name, description) VALUES ($1,$2) RETURNING id;

-- name: GetProjectByID :one
SELECT * FROM projects WHERE id=$1;

-- name: GetProjectsByUserID :many
SELECT *
FROM projects
WHERE creator=$1;
