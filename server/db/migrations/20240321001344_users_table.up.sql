CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "password" varchar NOT NULL,
    "first_name" varchar NOT NULL,
    "last_name" varchar NOT NULL
)
