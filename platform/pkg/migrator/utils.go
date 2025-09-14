package migrator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GenerateMigrationName генерирует имя файла миграции
func GenerateMigrationName(name string) string {
	// Удаляем специальные символы и заменяем пробелы на подчеркивания
	re := regexp.MustCompile(`[^a-zA-Z0-9_\s]`)
	cleanName := re.ReplaceAllString(name, "")
	cleanName = strings.ReplaceAll(cleanName, " ", "_")
	cleanName = strings.ToLower(cleanName)

	// Генерируем timestamp
	timestamp := time.Now().Format("20060102150405")

	return fmt.Sprintf("%s_%s", timestamp, cleanName)
}

// ParseMigrationVersion извлекает версию из имени файла миграции
func ParseMigrationVersion(filename string) (int64, error) {
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) < 1 {
		return 0, fmt.Errorf("invalid migration filename format: %s", filename)
	}

	version, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid version number in filename %s: %w", filename, err)
	}

	return version, nil
}

// ValidateMigrationName проверяет корректность имени миграции
func ValidateMigrationName(name string) error {
	if name == "" {
		return fmt.Errorf("migration name cannot be empty")
	}

	// Проверяем, что имя содержит только разрешенные символы
	re := regexp.MustCompile(`^[a-zA-Z0-9_\s]+$`)
	if !re.MatchString(name) {
		return fmt.Errorf("migration name can only contain letters, numbers, spaces and underscores")
	}

	// Проверяем длину
	if len(name) > 100 {
		return fmt.Errorf("migration name is too long (max 100 characters)")
	}

	return nil
}

// FormatDuration форматирует продолжительность для логов
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1e6)
	}
	if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	return d.String()
}

// IsTestDatabase проверяет, является ли база данных тестовой
func IsTestDatabase(dbName, testKeyword string) bool {
	return strings.Contains(strings.ToLower(dbName), strings.ToLower(testKeyword))
}

// SanitizeDatabaseName очищает имя базы данных от опасных символов
func SanitizeDatabaseName(name string) string {
	// Удаляем потенциально опасные символы
	re := regexp.MustCompile(`[^\w\-_]`)
	return re.ReplaceAllString(name, "")
}
