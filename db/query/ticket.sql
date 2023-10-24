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




-- name: ListAllTickets :many
SELECT * FROM tickets
WHERE (user_assigned = sqlc.narg('user_assigned') OR sqlc.narg('user_assigned') IS NULL)
AND (assigned_to = sqlc.narg('assigned_to') OR sqlc.narg('assigned_to') IS NULL)
AND (category_id = sqlc.narg('category_id') OR sqlc.narg('category_id') IS NULL)
ORDER BY ticket_id
LIMIT $1
OFFSET $2;



-- name: UpdateTicket :one
UPDATE tickets
SET updated_at = $2,
status = $3,
assigned_to = $4
WHERE ticket_id = $1
RETURNING *;

-- name: DeleteTicket :exec
DELETE FROM tickets
WHERE ticket_id = $1;
