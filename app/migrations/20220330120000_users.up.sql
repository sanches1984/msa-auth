CREATE TABLE "users" (
    "id" SERIAL NOT NULL PRIMARY KEY,
    "login" VARCHAR(100) NOT NULL,
    "password_hash" VARCHAR(255) NOT NULL,
    "created" TIMESTAMPTZ NOT NULL,
    "updated" TIMESTAMPTZ NOT NULL,
    "deleted" TIMESTAMPTZ
);