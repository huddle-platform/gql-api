-- name: GetProjects :many
SELECT * FROM projects;

-- name: CreateProject :one
INSERT into projects (name, description,creator) VALUES ($1,$2,$3) RETURNING id;

-- name: GetProjectByID :one
SELECT * FROM projects WHERE id=$1;

-- name: GetProjectsByUserID :many
SELECT *
FROM projects
WHERE creator=$1;

-- name: GetParticipantsOfProject :many
SELECT DISTINCT users.*
FROM users INNER JOIN participations
ON users.id = participations.user_id
WHERE participations.project_id = $1;

-- name: AddParticipantToProject :exec
INSERT INTO participations (user_id,project_id) VALUES($1,$2);


-- name: RemoveParticipantFromProject :exec
DELETE FROM participations WHERE user_id = $1 AND project_id = $2;