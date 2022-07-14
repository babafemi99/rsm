CREATE TABLE "User" (
  "id" uuid PRIMARY KEY,
  "firstname" varchar NOT NULL,
  "lastname" varchar NOT NULL,
  "email" varchar NOT NULL,
  "created_at" timestamptz NOT NULL
);

CREATE TABLE "Restaurants" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "location" varchar,
  "description" varchar,
  "created_at" timestamptz NOT NULL,
  "status" boolean NOT NULL
);

CREATE TABLE "Menu" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "restaurant_id" uuid NOT NULL,
  "item" varchar NOT NULL,
  "price" varchar NOT NULL,
  "item_type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL,
  "servings" int
);

ALTER TABLE "Menu" ADD FOREIGN KEY ("restaurant_id") REFERENCES "Restaurants" ("id");
