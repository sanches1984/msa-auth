CREATE TABLE "refresh_tokens"
(
    "id"            SERIAL       NOT NULL PRIMARY KEY,
    "user_id"       BIGINT       NOT NULL,
    "session_id"    UUID         NOT NULL,
    "token"         VARCHAR(64)  NOT NULL,
    "password_hash" VARCHAR(255) NOT NULL,
    "expires_in"    INT          NOT NULL,
    "created"       TIMESTAMPTZ  NOT NULL,
    "updated"       TIMESTAMPTZ  NOT NULL
);