
-- Fares
-- name: GetFare :one
SELECT * FROM fares WHERE id = $1 LIMIT 1;

-- name: ListFares :many
SELECT * FROM fares ORDER BY id;

-- name: CreateFare :one
INSERT INTO fares (
  base, per_km, per_min
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateFare :exec
UPDATE fares
SET base = $2,
    per_km = $3,
    per_min = $4
WHERE id = $1;

-- name: DeleteFare :exec
DELETE FROM fares WHERE id = $1;
