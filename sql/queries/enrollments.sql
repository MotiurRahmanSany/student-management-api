-- name: EnrollStudent :one
INSERT INTO enrollments (
    student_id,
    course_id
) VALUES (
    $1, $2
)
RETURNING id, student_id, course_id, enrolled_at;


-- name: UnenrollStudent :exec
DELETE FROM enrollments
WHERE student_id = $1 AND course_id = $2;


-- name: ListCoursesByStudent :many
SELECT c.id, c.title, c.code, c.credit, c.capacity, c.created_at, c.updated_at
FROM courses c
JOIN enrollments e ON e.course_id = c.id
WHERE e.student_id = $1;


-- name: ListStudentsByCourse :many
SELECT s.id, s.full_name, s.department, s.status, s.created_at
FROM students s
JOIN enrollments e ON e.student_id = s.id
WHERE e.course_id = $1;