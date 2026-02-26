-- +goose Up
-- +goose StatementBegin

CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    code TEXT NOT NULL UNIQUE,
    credit INT NOT NULL CHECK (credit > 0),
    capacity INT NOT NULL CHECK (capacity > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_courses_code ON courses(code);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS courses;

-- +goose StatementEnd