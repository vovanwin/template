// Code generated by ent, DO NOT EDIT.

package gen

import (
	"context"
	"fmt"
	"log"

	"entgo.io/ent/dialect/sql"
)

// Database is the client that holds all ent builders.
type Database struct {
	client *Client
}

// NewDatabase creates a new database based on Client.
func NewDatabase(client *Client) *Database {
	return &Database{client: client}
}

// RunInTx runs the given function f within a transaction.
// Inspired by https://entgo.io/docs/transactions/#best-practices.
// If there is already a transaction in the context, then the method uses it.
func (db *Database) RunInTx(ctx context.Context, f func(context.Context) error) error {
	tx := TxFromContext(ctx) // берем существующую транзакцию, инче создаем новую
	var err error
	if nil == tx {
		tx, err = db.client.Tx(ctx)
		if err != nil {
			return fmt.Errorf("create transaxtion: %v", err)
		}
	}

	var rollbackErr error
	defer func() {
		if r := recover(); r != nil {
			// Откатываем транзакию в случае паники
			rollbackErr = tx.Rollback()
			if rollbackErr != nil {
				log.Printf("Rollback error: %v", rollbackErr)
			}
			panic(r)
		} else if err != nil {
			// Откатываем транзакцию в случае ошибки
			rollbackErr = tx.Rollback()
			if rollbackErr != nil {
				log.Printf("Rollback error: %v", rollbackErr)
			}
		} else {
			// Коммитим транзакцию если всё хорошо
			err = tx.Commit()
			if err != nil {
				log.Printf("Commit error: %v", err)
			}
		}
	}()

	if err := f(NewTxContext(ctx, tx)); err != nil {
		rollbackErr = tx.Rollback()
		return fmt.Errorf("NewTxContext transaction: %w", err)
	}

	return nil
}

func (db *Database) loadClient(ctx context.Context) *Client {
	tx := TxFromContext(ctx)
	if tx != nil {
		return tx.Client()
	}
	return db.client
}

// Exec executes a query that doesn't return rows. For example, in SQL, INSERT or UPDATE.
func (db *Database) Exec(ctx context.Context, query string, args ...any) (*sql.Result, error) {
	var res sql.Result
	err := db.loadClient(ctx).driver.Exec(ctx, query, args, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Query executes a query that returns rows, typically a SELECT in SQL.
func (db *Database) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	var rows sql.Rows
	err := db.loadClient(ctx).driver.Query(ctx, query, args, &rows)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}

// Post is the client for interacting with the Post builders.
func (db *Database) Post(ctx context.Context) *PostClient {
	return db.loadClient(ctx).Post
}

// User is the client for interacting with the User builders.
func (db *Database) User(ctx context.Context) *UserClient {
	return db.loadClient(ctx).User
}
