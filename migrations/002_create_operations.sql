CREATE TABLE IF NOT EXISTS operations
(
    id         UUID PRIMARY KEY,
    name       TEXT      NOT NULL,
    type       TEXT      NOT NULL CHECK (type in ('income', 'expense')),
    amount     numeric   NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    user_id    UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_operations_user_id ON operations (user_id);