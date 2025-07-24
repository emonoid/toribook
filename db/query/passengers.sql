-- Passengers
-- name: GetPassenger :one
SELECT * FROM passengers WHERE id = $1 LIMIT 1;

-- name: GetPassengerByEmail :one
SELECT * FROM passengers WHERE email = $1 LIMIT 1;

-- name: ListPassengers :many
SELECT * FROM passengers ORDER BY full_name;

-- name: CreatePassenger :one
INSERT INTO passengers (
  hashed_password, full_name, email, rating
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdatePassenger :exec
UPDATE passengers
SET hashed_password = $2,
    full_name = $3,
    email = $4,
    rating = $5
WHERE id = $1;

-- name: DeletePassenger :exec
DELETE FROM passengers WHERE id = $1;