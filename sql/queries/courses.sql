-- name: CreateCourse :one
INSERT INTO courses (
    title,
    code,
    credit,
    capacity
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, title, code, credit, capacity, created_at, updated_at;


-- name: GetCourseByID :one
SELECT id, title, code, credit, capacity, created_at, updated_at
FROM courses
WHERE id = $1
LIMIT 1;


-- name: ListCourses :many
SELECT id, title, code, credit, capacity, created_at, updated_at
FROM courses
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: UpdateCourse :one
UPDATE courses
SET
    title = $2,
    code = $3,
    credit = $4,
    capacity = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING id, title, code, credit, capacity, created_at, updated_at;


-- name: DeleteCourse :exec
DELETE FROM courses
WHERE id = $1;