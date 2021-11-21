-- name: GetProjects :many
SELECT * FROM projects;

-- name: CreateProject :one
INSERT into projects (name, description) VALUES ($1,$2) RETURNING id;

-- name: GetProjectByID :one
SELECT *
FROM projects
WHERE id=$1;

-- name: AddRole :one
INSERT INTO roles (type,project_id) VALUES ($1,$2) RETURNING id;
-- name: GrantRoleToUser :exec
INSERT INTO has_role (user_id,role_id) VALUES ($1,$2);