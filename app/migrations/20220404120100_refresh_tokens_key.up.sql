ALTER TABLE "refresh_tokens" ADD CONSTRAINT "fk_refresh_tokens_users"
    FOREIGN KEY("user_id") REFERENCES "users"("id")
    ON DELETE CASCADE
    ON UPDATE CASCADE;