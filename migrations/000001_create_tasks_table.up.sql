CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    difficulty INTEGER NOT NULL,
    status INTEGER NOT NULL,
    last_update TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_tasks_status ON tasks(status);