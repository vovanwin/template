-- Получить пользователя для me запроса
-- name: FindMeForId :one
 SELECT * FROM users WHERE id = $1;