CREATE TABLE "image_conversions" (
    "conversion_id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL,
    "image_name" VARCHAR(255) NOT NULL,
    "extracted_text" TEXT NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "image_conversions" ("user_id");

ALTER TABLE "image_conversions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;