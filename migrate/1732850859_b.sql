-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS "b" (
	"id" bigint NOT NULL,
	"value" json NOT NULL,
	"ptr_value" json,
	"n" varchar(255) NOT NULL
);
ALTER TABLE "b"
	ADD COLUMN "id" bigint NOT NULL,
	ADD COLUMN "value" json NOT NULL,
	ADD COLUMN "ptr_value" json,
	ADD COLUMN "n" varchar(255) NOT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE "b"
	DROP COLUMN "id",
	DROP COLUMN "value",
	DROP COLUMN "ptr_value",
	DROP COLUMN "n";