-- name: GetAllPermissionsForUser :one
SELECT permissions.code
FROM permissions
         INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
         INNER JOIN users ON users_permissions.user_id = users.id
WHERE users.id = $1;

-- name: AddPermissionsForUser :many
INSERT INTO users_permissions (user_id, permission_id)
SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)
RETURNING user_id, permission_id;