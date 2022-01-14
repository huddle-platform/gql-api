
-- name: CreateUser :one
INSERT INTO users (id,username, profile_image) VALUES ($1,$2,$3) RETURNING id;

-- name: GetUserByID :one
SELECT * FROM users WHERE id=$1;

-- name: GetUserParticipations :many
SELECT projects.*
FROM participations INNER JOIN projects
ON participations.project_id = projects.id
WHERE participations.user_id = $1;

-- name: GetSavedProjectsForUser :many
select projects.*
FROM projects INNER JOIN project_saves
ON projects.id = project_saves.project_id
WHERE project_saves.user_id = $1;

-- name: SaveProjectForUser :exec
INSERT INTO project_saves (user_id,project_id) VALUES ($1,$2);

-- name: UnsaveProjectForUser :exec
DELETE FROM project_saves WHERE user_id = $1 AND project_id = $2;

-- name: SetUserName :exec
UPDATE users SET username = $1 WHERE id = $2;
-- name: SetDescription :exec
UPDATE users SET description = $1 WHERE id = $2;
-- name: SetProfileImage :exec
UPDATE users SET profile_image = $1 WHERE id = $2;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username=$1;