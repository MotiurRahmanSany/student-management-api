-- +goose Up
-- +goose StatementBegin

CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    full_name TEXT NOT NULL,
    date_of_birth DATE,
    department TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('active', 'suspended', 'graduated')) DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_students_department ON students(department);
CREATE INDEX idx_students_status ON students(status);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS students;

-- +goose StatementEnd