
-- Drivers
-- name: GetDriver :one
SELECT * FROM drivers WHERE id = $1 LIMIT 1;

-- name: GetDriverByMobile :one
SELECT * FROM drivers WHERE mobile = $1 LIMIT 1;

-- name: ListDrivers :many
SELECT * FROM drivers ORDER BY full_name;

-- name: CreateDriver :one
INSERT INTO drivers (
  hashed_password, full_name, driving_license, mobile, car_id, car_type, car_image, online_status, rating, profile_status, subscription_status, subscription_package, subscription_amount, subscription_validity
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING *;

-- name: UpdateDriver :exec
UPDATE drivers
SET hashed_password = $2,
    full_name = $3,
    driving_license = $4,
    mobile = $5,
    car_id = $6,
    car_type = $7,
    car_image = $8,
    online_status = $9,
    rating = $10,
    profile_status = $11,
    subscription_status = $12,
    subscription_package = $13,
    subscription_amount = $14,
    subscription_validity = $15
WHERE id = $1;

-- name: DeleteDriver :exec
DELETE FROM drivers WHERE id = $1;
