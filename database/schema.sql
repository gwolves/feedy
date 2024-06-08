-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";
-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';

CREATE TABLE public."feeds" (
  "id" bigserial,
  "name" varchar NOT NULL,
  "url" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  UNIQUE ("url")
);

CREATE TABLE public."subscriptions" (
  "id" bigserial,
  "bot_name" varchar,
  "feed_id" bigint NOT NULL,
  "channel_id" varchar NOT NULL,
  "group_id" varchar NOT NULL,
  "published_at" timestamptz NOT NULL DEFAULT now(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  UNIQUE ("feed_id", "channel_id", "group_id")
);

CREATE TABLE public."audit_logs" (
  "id" bigserial,
  "actor" varchar NOT NULL,
  "action" varchar NOT NULL,
  "target" varchar NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
