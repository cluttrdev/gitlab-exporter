package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"
)

type preparer interface {
	PrepareContext(context.Context, string) (*sql.Stmt, error)
}

func prepareBatchInsert(ctx context.Context, prep preparer, table string, rows int, cols int) (*sql.Stmt, error) {
	if rows < 0 {
		return nil, errors.New("number of rows cannot be negative")
	}
	if cols < 0 {
		return nil, errors.New("number of columns cannot be negative")
	}

	colVals := "(" + strings.Join(slices.Repeat([]string{"?"}, cols), ",") + ")" // (?,?,?)
	vals := strings.Join(slices.Repeat([]string{colVals}, rows), ",")            // (?,?,?),(?,?,?),(?,?,?)

	query := fmt.Sprintf("INSERT OR REPLACE INTO %s VALUES %s", table, vals)
	return prep.PrepareContext(ctx, query)
}

func withTransaction(ctx context.Context, db *sql.DB, query func(context.Context, *sql.Tx) error) error {
	// cancel after 5 seconds
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// get connection
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("get connection: %w", err)
	}
	defer func() { _ = conn.Close() }()

	// begin transaction
	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// run transaction
	if err := query(ctx, tx); err != nil {
		if rerr := tx.Rollback(); err != nil {
			err = errors.Join(err, fmt.Errorf("rollback transaction: %w", rerr))
		}
		return fmt.Errorf("run transaction: %w", err)
	}

	// end transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
