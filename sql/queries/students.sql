-- name: CreateStudent :one
INSERT INTO students (
    user_id,
    full_name,
    date_of_birth,
    department,
    status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, user_id, full_name, date_of_birth, department, status, created_at, updated_at;


-- name: GetStudentByID :one
SELECT id, user_id, full_name, date_of_birth, department, status, created_at, updated_at
FROM students
WHERE id = $1
LIMIT 1;


-- name: GetStudentByUserID :one
SELECT id, user_id, full_name, date_of_birth, department, status, created_at, updated_at
FROM students
WHERE user_id = $1
LIMIT 1;


-- name: ListStudents :many
SELECT id, user_id, full_name, date_of_birth, department, status, created_at, updated_at
FROM students
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: CountStudents :one
SELECT COUNT(*) FROM students;


-- name: UpdateStudent :one
UPDATE students
SET
    full_name = $2,
    date_of_birth = $3,
    department = $4,
    status = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, full_name, date_of_birth, department, status, created_at, updated_at;


-- name: DeleteStudent :exec
DELETE FROM students
WHERE id = $1;


-- name: GetStudentsByDepartment :many
SELECT id, user_id, full_name, date_of_birth, department, status, created_at, updated_at
FROM students
WHERE department = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;