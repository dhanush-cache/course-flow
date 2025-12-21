-- name: GetCourse :one
SELECT id, slug, timestamp_id, platform_id
FROM courses
WHERE id = ?
LIMIT 1;

-- name: ListCourses :many
SELECT id, slug, timestamp_id, platform_id
FROM courses
ORDER BY id;

-- name: CreateCourse :one
INSERT INTO courses (id, slug, timestamp_id, platform_id)
VALUES (?, ?, ?, ?)
RETURNING id, slug, timestamp_id, platform_id;

-- name: DeleteCourse :exec
DELETE
FROM courses
WHERE id = ?;

-- name: CreateTimestamp :one
INSERT INTO timestamps (first, rest)
VALUES (?, ?)
RETURNING id, first, rest;

-- name: GetTimestamp :one
SELECT id, first, rest
FROM timestamps
WHERE id = ?
LIMIT 1;

-- name: CreateCourseURL :one
INSERT INTO course_urls (course_id, url, category, position)
VALUES (?, ?, ?, ?)
RETURNING id, course_id, url, category, position;

-- name: GetCourseURLs :many
SELECT id, course_id, url, category, position
FROM course_urls
WHERE course_id = ?
ORDER BY position;

-- name: UpdateCourseURL :exec
UPDATE course_urls
SET url      = ?,
    category = ?,
    position = ?
WHERE id = ?;

-- name: DeleteCourseURL :exec
DELETE
FROM course_urls
WHERE id = ?;

-- name: GetCourseBySlug :one
SELECT id, slug, timestamp_id, platform_id
FROM courses
WHERE slug = ?
LIMIT 1;

-- name: UpdateCourse :one
UPDATE courses
SET slug         = ?,
    timestamp_id = ?,
    platform_id  = ?
WHERE id = ?
RETURNING id, slug, timestamp_id, platform_id;

-- name: CreatePlatform :one
INSERT INTO platforms (id, title, url)
VALUES (?, ?, ?)
RETURNING id, title, url;

-- name: GetPlatform :one
SELECT id, title, url
FROM platforms
WHERE id = ?
LIMIT 1;

-- name: ListPlatforms :many
SELECT id, title, url
FROM platforms
ORDER BY title;

-- name: UpdatePlatform :one
UPDATE platforms
SET title = ?,
    url   = ?
WHERE id = ?
RETURNING id, title, url;

-- name: DeletePlatform :exec
DELETE
FROM platforms
WHERE id = ?;
