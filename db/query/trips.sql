
-- Trips
-- name: GetTrip :one
SELECT * FROM trips WHERE id = $1 LIMIT 1;

-- name: GetTripByBookingID :one
SELECT * FROM trips WHERE booking_id = $1 LIMIT 1;

-- name: ListTrips :many
SELECT * FROM trips ORDER BY id DESC;

-- name: CreateTrip :one
INSERT INTO trips (
  booking_id, trip_status, pickup_location, pickup_lat, pickup_long, dropoff_location, dropoff_lat, dropoff_long, driver_id, driver_name, driver_mobile, car_id, car_type, car_image, fare
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: UpdateTripStatus :exec
UPDATE trips
SET trip_status = $2
WHERE id = $1;

-- name: DeleteTrip :exec
DELETE FROM trips WHERE id = $1;
