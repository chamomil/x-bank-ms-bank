CREATE TYPE status_account AS ENUM ('ACTIVE', 'BLOCKED');
CREATE TYPE status_transaction AS ENUM ('BLOCKED', 'CONFIRMED', 'CANCELLED');

CREATE TABLE "accounts"
(
    "id"           BIGSERIAL PRIMARY KEY,
    "balanceCents" BIGINT         NOT NULL CHECK ( "balanceCents" >= 0 ),
    "ownerId"      BIGINT         NOT NULL REFERENCES "accountOwners" ("id"),
    "status"       status_account NOT NULL DEFAULT 'ACTIVE'
);

CREATE TABLE "atms"
(
    "id"        BIGSERIAL PRIMARY KEY,
    "cashCents" BIGINT             NOT NULL CHECK ( "cashCents" >= 0 ),
    "login"     VARCHAR(32) UNIQUE NOT NULL CHECK ( login ~ '^[a-z0-9_-]+$'),
    "password"  BYTEA              NOT NULL CHECK ( length(password) <= 60 )
);

CREATE TABLE "accountOwners"
(
    "id"     BIGSERIAL PRIMARY KEY,
    "userId" BIGINT,
    "atmId"  BIGINT REFERENCES "atms" ("id")
);

CREATE TABLE "transactions"
(
    "id"          BIGINT             NOT NULL PRIMARY KEY,
    "senderId"    BIGINT             NOT NULL REFERENCES "accounts" ("id"),
    "receiverId"  BIGINT             NOT NULL REFERENCES "accounts" ("id"),
    "status"      status_transaction NOT NULL DEFAULT 'BLOCKED',
    "createdAt"   TIMESTAMP          NOT NULL DEFAULT current_timestamp,
    "amountCents" BIGINT             NOT NULL,
    "description" TEXT
);

CREATE INDEX "transactions_senderId_index" ON "transactions" ("senderId");
CREATE INDEX "transactions_receiverId_index" ON "transactions" ("receiverId");
CREATE INDEX "accounts_ownerId_index" ON "accounts" ("ownerId")