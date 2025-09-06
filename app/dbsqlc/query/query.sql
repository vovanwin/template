-- Получить пользователя для me запроса
-- name: FindMeForId :one
SELECT * FROM users WHERE id = $1;

-- Получить пользователя по email
-- name: GetUserByEmail :one
SELECT id, email, password_hash, first_name, last_name, role, tenant_id, 
       is_active, email_verified, settings, components, created_at, updated_at
FROM users 
WHERE email = $1 AND is_active = true;

-- Получить пользователя по ID
-- name: GetUserByID :one
SELECT id, email, password_hash, first_name, last_name, role, tenant_id, 
       is_active, email_verified, settings, components, created_at, updated_at
FROM users 
WHERE id = $1 AND is_active = true;

-- Создать нового пользователя
-- name: CreateUser :one
INSERT INTO users (id, email, password_hash, first_name, last_name, role, tenant_id, 
                  is_active, email_verified, settings, components)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, email, password_hash, first_name, last_name, role, tenant_id, 
          is_active, email_verified, settings, components, created_at, updated_at;

-- Обновить пользователя
-- name: UpdateUser :one
UPDATE users 
SET email = $2, password_hash = $3, first_name = $4, last_name = $5, 
    role = $6, tenant_id = $7, is_active = $8, email_verified = $9,
    settings = $10, components = $11, updated_at = NOW()
WHERE id = $1
RETURNING id, email, password_hash, first_name, last_name, role, tenant_id, 
          is_active, email_verified, settings, components, created_at, updated_at;