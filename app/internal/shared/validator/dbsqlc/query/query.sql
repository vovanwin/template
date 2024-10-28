-- Проверить что у пользователя есть разрешения на все устройства
-- name: CountUserDevicesWithPermissions :one
SELECT COUNT(devices.id)
FROM devices
         JOIN users_has_devices ON users_has_devices.device_id = devices.id
WHERE devices.uuid = ANY(@uuidDevices::uuid[])
  AND devices.tenant_id = @tenant_id
  AND users_has_devices.user_id = @user_id
  AND users_has_devices.read = true;
