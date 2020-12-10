-- +migrate Up
CREATE TABLE IF NOT EXISTS "Payment"
(
    id              UUID PRIMARY KEY NOT NULL,
    organisation_id UUID             NOT NULL,
    version         INT              NOT NULL,
    is_deleted      BOOLEAN          NOT NULL,
    is_locked       BOOLEAN          NOT NULL,
    actioned_by     UUID             NOT NULL,
    created_on      TIMESTAMP,
    modified_on     TIMESTAMP,
    record          JSONB,
    pagination_id   SERIAL
);

ALTER TABLE "Payment" REPLICA IDENTITY FULL;

CREATE UNIQUE INDEX Payment_paginationid ON "Payment" (pagination_id);

-- +migrate Down
DROP TABLE IF EXISTS "Payment";
