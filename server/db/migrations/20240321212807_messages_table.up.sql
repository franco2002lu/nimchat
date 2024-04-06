CREATE TABLE "messages" (
    "message_id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "chat_id" bigserial NOT NULL,
    "content" varchar NOT NULL,
    "votes" bigserial NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp NOT NULL DEFAULT now()
)