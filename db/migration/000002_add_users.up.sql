CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_change_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00z',
  "create_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD CONSTRAINT accounts_owner_currency_key UNIQUE ("owner", "currency");

ALTER TABLE "accounts" ADD CONSTRAINT accounts_owner_fkey FOREIGN KEY ("owner") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;