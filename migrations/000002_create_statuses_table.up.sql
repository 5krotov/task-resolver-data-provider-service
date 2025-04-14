CREATE TABLE statuses (
    id BIGSERIAL PRIMARY KEY,
    status INTEGER NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE INDEX idx_statuses_task_id ON statuses(task_id);