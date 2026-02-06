package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// TxManager управляет транзакциями.
type TxManager struct {
	conn *pgx.Conn
}

// NewTxManager создает новый TxManager.
func NewTxManager(conn *pgx.Conn) *TxManager {
	return &TxManager{conn: conn}
}

// RunInTx выполняет функцию f внутри транзакции.
// Если транзакция уже существует в контексте, используется она.
func (tm *TxManager) RunInTx(ctx context.Context, f func(context.Context) error) (err error) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if !ok {
		// Создаем новую транзакцию, если она еще не существует
		tx, err = tm.conn.Begin(ctx)
		if err != nil {
			return fmt.Errorf("создание транзакции: %v", err)
		}

		// Обрабатываем завершение транзакции
		defer func() {
			if p := recover(); p != nil {
				// Откатываем транзакцию в случае паники
				if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
					log.Printf("Ошибка отката транзакции: %v", rollbackErr)
				}
				panic(p)
			} else if err != nil {
				// Откатываем транзакцию в случае ошибки
				if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
					log.Printf("Ошибка отката транзакции: %v", rollbackErr)
				}
			} else {
				// Коммитим транзакцию, если все прошло успешно
				if commitErr := tx.Commit(ctx); commitErr != nil {
					log.Printf("Ошибка коммита транзакции: %v", commitErr)
				}
			}
		}()

		// Добавляем транзакцию в контекст
		ctx = context.WithValue(ctx, txKey{}, tx)
	}

	// Выполняем функцию внутри транзакции
	if err := f(ctx); err != nil {
		return fmt.Errorf("выполнение функции внутри транзакции: %w", err)
	}

	return nil
}

// txKey используется как ключ для хранения транзакции в контексте.
type txKey struct{}

// TxFromContext получает транзакцию из контекста.
func TxFromContext(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(txKey{}).(pgx.Tx)
	return tx
}
