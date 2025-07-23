
-- Cars
-- name: GetCar :one
SELECT * FROM cars WHERE id = $1 LIMIT 1;

-- name: ListCars :many
SELECT * FROM cars ORDER BY car_type;

-- name: CreateCar :one
INSERT INTO cars (
  car_type, car_model, car_image
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateCar :exec
UPDATE cars
SET car_type = $2,
    car_model = $3,
    car_image = $4
WHERE id = $1;

-- name: DeleteCar :exec
DELETE FROM cars WHERE id = $1;
