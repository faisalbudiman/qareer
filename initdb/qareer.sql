CREATE TABLE "locations" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar UNIQUE NOT NULL,
  "active" bool,
  "created_at" timestamp DEFAULT now(),
  "updated_at" timestamp DEFAULT now()
);