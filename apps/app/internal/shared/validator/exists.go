package validator

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"strings"
)

var (
	// Проверка допустимости таблиц и полей, для доп безопасности.
	// Теоретически вызвать sql инъекцию нельзя, так как название таблиц и полей задаем мы сами,
	// но ради доп безопасности можно задать только те значения которые есть в словаре.
	// Используется в existsValidate и uniqueValidate
	validTables = map[string]bool{"users": true, "devices": true, "reports_lists": true}
	validFields = map[string]bool{"username": true, "imei": true, "email": true}
)

// existsValidate проверяет, существует ли запись в БД
func (cv *CustomValidator) existsValidate(ctx context.Context, fl validator.FieldLevel) bool {
	fieldValue := fl.Field().Int() // Значение поля, которое мы проверяем (ID)
	tableName := fl.Param()        // Имя таблицы, в которой мы проверяем наличие записи

	if !validTables[tableName] {
		fmt.Println("Неверное имя таблицы")
		return false
	}

	// Создаем запрос с использованием Squirrel
	query, args, err := cv.pgx.Builder.Select("1").
		From(tableName).
		Where(squirrel.Eq{"id": fieldValue}).
		ToSql()
	if err != nil {
		fmt.Printf("Ошибка формирования SQL-запроса: %v", err)
		return false
	}

	// Выполняем запрос
	var exists bool
	err = cv.pgx.Pool.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		fmt.Printf("Ошибка выполнения SQL-запроса: %v", err)
		return false
	}
	return exists
}

// uniqueValidate проверяет, существует ли запись с таким значением в указанной таблице и поле
func (cv *CustomValidator) uniqueValidate(ctx context.Context, fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String() // Значение поля (например, username или imei)
	param := fl.Param()               // Получаем параметр (формат: "table,field")

	// Разбираем параметр на таблицу и поле
	params := strings.Split(param, ":")
	if len(params) != 2 {
		fmt.Printf("Неверный формат тэга unique: %s", param)
		return false
	}
	tableName := params[0]
	fieldName := params[1]

	if !validTables[tableName] || !validFields[fieldName] {
		fmt.Println("Неверное имя таблицы или поля")
		return false
	}

	// Строим запрос с использованием Squirrel
	query, args, err := cv.pgx.Builder.
		Select("1").
		From(tableName).
		Where(squirrel.Eq{fieldName: fieldValue}).
		Limit(1). // Оптимизация — запросим только одну строку
		ToSql()

	if err != nil {
		fmt.Printf("Ошибка формирования SQL-запроса: %v", err)
		return false
	}

	// Выполняем запрос
	var exists bool
	err = cv.pgx.Pool.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		// Если не нашли строку, считаем что поле уникально
		if err == pgx.ErrNoRows {
			return true
		}
		// Если другая ошибка, логируем её и возвращаем false
		fmt.Printf("Ошибка выполнения запроса: %v", err)
		return false
	}

	// Если запись найдена, то поле не уникально
	return !exists
}
