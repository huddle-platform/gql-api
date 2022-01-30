-- name: WriteMessage :exec
INSERT INTO messages (sender_id,receiver_id,content) VALUES ($1,$2,$3);


-- name: WriteProjectUserMessage :exec
INSERT INTO projectMessages (project_id,user_id,userIsSender,content) VALUES ($1,$2,$3,$4);


-- name: GetMessagesBetweenUsers :many
SELECT * FROM messages
    WHERE ((sender_id=$1 AND receiver_id=$2) OR (sender_id=$2 AND receiver_id=$1))
        AND time <= $3 
    ORDER BY time DESC LIMIT $4;

-- name: GetChatsWithUser :many
SELECT DISTINCT ON(messages.sender_id) users.*
FROM messages INNER JOIN users
ON messages.sender_id = users.id
WHERE messages.receiver_id=$1
UNION
SELECT DISTINCT ON(messages.receiver_id) users.*
FROM messages INNER JOIN users
ON messages.receiver_id = users.id
WHERE messages.sender_id=$1;

-- name: GetProjectChatsWithUser :many
SELECT DISTINCT ON(projects.id) projects.*
FROM projectMessages INNER JOIN projects
ON projectMessages.project_id = projects.id
WHERE projectMessages.user_id=$1;


-- name: GetChatsWithProject :many
SELECT DISTINCT ON(users.user_id) users.*
    FROM projectMessages
    INNER JOIN users ON projectMessages.user_id = users.id
    WHERE projectMessages.project_id=$1;

-- name: GetChatsWithCreatedProjects :many
SELECT DISTINCT ON(projectMessages.user_id) projectMessages.user_id,projects.id
    FROM projectMessages INNER JOIN projects
    ON projectMessages.project_id = projects.id
    WHERE projects.creator=$1;

-- name: GetMessagesBetweenUserAndProject :many
SELECT * FROM projectMessages
    WHERE project_id = $1 AND user_id = $2 AND time <= $3
    ORDER BY time DESC LIMIT $4;