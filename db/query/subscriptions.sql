
-- Subscriptions
-- name: GetSubscription :one
SELECT * FROM subscriptions WHERE id = $1 LIMIT 1;

-- name: ListSubscriptions :many
SELECT * FROM subscriptions ORDER BY id;

-- name: CreateSubscription :one
INSERT INTO subscriptions (
  subscription_package, subscription_amount, subscription_validity, status
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateSubscription :exec
UPDATE subscriptions
SET subscription_package = $2,
    subscription_amount = $3,
    subscription_validity = $4,
    status = $5
WHERE id = $1;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions WHERE id = $1;
