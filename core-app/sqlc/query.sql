-- name: GetTitle :one
SELECT title FROM assets
WHERE id = ?;

-- name: GetAboutMe :one
SELECT aboutMe FROM assets
WHERE id = ?;