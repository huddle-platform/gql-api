-- name: WriteMessage :exec
INSERT INTO messages (sender_id,receiver_id,content) VALUES ($1,$2,$3);


-- name: WriteProjectMessageToUser :exec
INSERT INTO projectMessages (project_id,user_id,content) VALUES ($1,$2,$3);

-- name: WriteUserMessageToProject :exec
INSERT INTO projectMessages (user_id,project_id,content) VALUES ($1,$2,$3);

-- name: GetMessagesBetweenUsers :many
SELECT * FROM messages
    WHERE ((sender_id=$1 AND receiver_id=$2) OR (sender_id=$2 AND receiver_id=$1))
        AND time <= $3 
    ORDER BY time DESC LIMIT $4;

-- name: GetChatsWithUser :many
SELECT DISTINCT ON(sender_id) sender_id AS id FROM messages
WHERE messages.receiver_id=$1
UNION
SELECT DISTINCT ON(receiver_id) receiver_id AS id FROM messages
WHERE messages.sender_id=$1;

-- name: GetProjectChatsWithUser :many
SELECT DISTINCT ON(project_id) project_id AS id FROM projectMessages
WHERE user_id=$1;


-- name: GetChatsWithProject :many
SELECT DISTINCT ON(user_id) user_id AS id FROM projectMessages
WHERE project_id=$1;

-- name: GetChatsWithCreatedProjects :many
SELECT DISTINCT ON(projectMessages.user_id) projectMessages.user_id
    FROM projectMessages INNER JOIN projects
    ON projectMessages.project_id = projects.id
    WHERE projects.creator=$1;

-- name: GetMessagesBetweenUserAndProject :many
SELECT * FROM projectMessages
    WHERE project_id = $1 AND user_id = $2 AND time <= $3
    ORDER BY time DESC LIMIT $4;