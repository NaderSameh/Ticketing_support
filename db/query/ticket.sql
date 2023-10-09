-- name: CreateTicket :one
INSERT INTO tickets (
  title, description, status, user_assigned, category_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;


-- name: GetTicket :one
SELECT * FROM tickets
WHERE ticket_id = $1 LIMIT 1;


-- name: GetTicketForUpdate :one
SELECT * FROM tickets
WHERE ticket_id = $1 LIMIT 1
FOR NO KEY UPDATE;


-- name: ListTickets :many
SELECT * FROM tickets
WHERE user_assigned = $1
ORDER BY ticket_id
LIMIT $2
OFFSET $3;


-- name: UpdateTicket :one
UPDATE tickets
SET updated_at = $2
WHERE ticket_id = $1
RETURNING *;

-- name: DeleteTicket :exec
DELETE FROM tickets
WHERE ticket_id = $1;
