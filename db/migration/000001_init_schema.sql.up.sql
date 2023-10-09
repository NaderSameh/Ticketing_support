

CREATE TABLE "tickets" (
  "ticket_id" bigserial PRIMARY KEY,
  "title" varchar NOT NULL,
  "description" varchar NOT NULL,
  "status" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "closed_at" timestamptz,
  "category_id" bigserial NOT NULL,
  "user_assigned" varchar NOT NULL,
  "assigned_to" varchar
);

CREATE TABLE "categories" (
  "category_id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL
);
CREATE TABLE "comments" (
  "comment_id" bigserial PRIMARY KEY,
  "comment_text" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "ticket_id" bigserial NOT NULL,
  "user_commented" varchar NOT NULL
);

ALTER TABLE "tickets" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("category_id");

ALTER TABLE "comments" ADD FOREIGN KEY ("ticket_id") REFERENCES "tickets" ("ticket_id");
