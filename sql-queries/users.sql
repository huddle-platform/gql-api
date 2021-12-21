
-- name: CreateUser :one
INSERT INTO users (id,username, email, profile_image) VALUES ($1,$2,$3,$4) RETURNING id;

-- name: GetUserByID :one
SELECT * FROM users WHERE id=$1;