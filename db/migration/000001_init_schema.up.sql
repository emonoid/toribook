CREATE TABLE "passengers" (
  "id" bigserial PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "rating" int NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "drivers" (
  "id" bigserial PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "driving_license" varchar UNIQUE NOT NULL,
  "mobile" varchar UNIQUE NOT NULL,
  "car_id" bigserial NOT NULL,
  "car_type" varchar NOT NULL,
  "car_image" text NOT NULL,
  "online_status" bool NOT NULL,
  "rating" int NOT NULL,
  "profile_status" int NOT NULL,
  "subscription_status" bool NOT NULL,
  "subscription_package" varchar NOT NULL,
  "subscription_amount" varchar NOT NULL,
  "subscription_validity" int NOT NULL,
  "subscription_expire_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "subscriptions" (
  "id" bigserial PRIMARY KEY,
  "subscription_package" varchar NOT NULL,
  "subscription_amount" varchar NOT NULL,
  "subscription_validity" int NOT NULL,
  "status" bool NOT NULL
);

CREATE TABLE "trips" (
  "id" bigserial PRIMARY KEY,
  "booking_id" varchar UNIQUE NOT NULL,
  "trip_status" varchar NOT NULL,
  "pickup_location" varchar NOT NULL,
  "pickup_lat" varchar NOT NULL,
  "pickup_long" varchar NOT NULL,
  "dropoff_location" varchar NOT NULL,
  "dropoff_lat" varchar NOT NULL,
  "dropoff_long" varchar NOT NULL,
  "driver_id" bigserial,
  "driver_name" varchar,
  "driver_mobile" varchar,
  "car_id" bigserial,
  "car_type" varchar,
  "car_image" text,
  "fare" int,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "cars" (
  "id" bigserial PRIMARY KEY,
  "car_type" varchar NOT NULL,
  "car_model" varchar NOT NULL,
  "car_image" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

CREATE TABLE "fares" (
  "id" bigserial PRIMARY KEY,
  "base" int NOT NULL,
  "per_km" int NOT NULL,
  "per_min" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);