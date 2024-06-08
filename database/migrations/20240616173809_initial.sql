-- Create "audit_logs" table
CREATE TABLE "audit_logs" ("id" bigserial NOT NULL, "actor" character varying NOT NULL, "action" character varying NOT NULL, "target" character varying NULL, "created_at" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("id"));
-- Create "feeds" table
CREATE TABLE "feeds" ("id" bigserial NOT NULL, "name" character varying NOT NULL, "url" character varying NOT NULL, "created_at" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("id"), CONSTRAINT "feeds_url_key" UNIQUE ("url"));
-- Create "subscriptions" table
CREATE TABLE "subscriptions" ("id" bigserial NOT NULL, "bot_name" character varying NULL, "feed_id" bigint NOT NULL, "channel_id" character varying NOT NULL, "group_id" character varying NOT NULL, "published_at" timestamptz NOT NULL DEFAULT now(), "created_at" timestamptz NOT NULL DEFAULT now(), PRIMARY KEY ("id"), CONSTRAINT "subscriptions_feed_id_channel_id_group_id_key" UNIQUE ("feed_id", "channel_id", "group_id"));
