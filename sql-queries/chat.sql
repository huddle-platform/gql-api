-- name: WriteMessage :exec
INSERT INTO messages (sender_id,receiver_id,content) VALUES ($1,$2,$3);

-- name: GetMessagesBetweenUsers :many
SELECT * FROM messages WHERE (sender_id=$1 AND receiver_id=$2) OR (sender_id=$2 AND receiver_id=$1) ORDER BY time DESC LIMIT $3 OFFSET $4;

-- name: GetChatsWithUser :many
SELECT DISTINCT ON(sender_id) sender_id AS id FROM messages
WHERE messages.receiver_id=$1
UNION
SELECT DISTINCT ON(receiver_id) receiver_id AS id FROM messages
WHERE messages.sender_id=$1;