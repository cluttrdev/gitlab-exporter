package sqlite

import (
	"context"
	"database/sql"
	"strings"
	"testing"
)

// mockPreparer implements the preparer interface for testing
type mockPreparer struct {
	lastQuery string
}

func (m *mockPreparer) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	m.lastQuery = query
	return nil, nil
}

func TestPrepareBatchInsert(t *testing.T) {
	tests := []struct {
		name      string
		table     string
		rows      int
		cols      int
		wantQuery string
		wantErr   bool
	}{
		{
			name:      "single row, single column",
			table:     "test_table",
			rows:      1,
			cols:      1,
			wantQuery: "INSERT OR REPLACE INTO test_table VALUES (?)",
			wantErr:   false,
		},
		{
			name:      "single row, multiple columns",
			table:     "test_table",
			rows:      1,
			cols:      3,
			wantQuery: "INSERT OR REPLACE INTO test_table VALUES (?,?,?)",
			wantErr:   false,
		},
		{
			name:      "multiple rows, single column",
			table:     "test_table",
			rows:      3,
			cols:      1,
			wantQuery: "INSERT OR REPLACE INTO test_table VALUES (?),(?),(?)",
			wantErr:   false,
		},
		{
			name:      "multiple rows, multiple columns",
			table:     "test_table",
			rows:      2,
			cols:      3,
			wantQuery: "INSERT OR REPLACE INTO test_table VALUES (?,?,?),(?,?,?)",
			wantErr:   false,
		},
		{
			name:      "zero rows",
			table:     "test_table",
			rows:      0,
			cols:      3,
			wantQuery: "INSERT OR REPLACE INTO test_table VALUES ",
			wantErr:   false,
		},
		{
			name:      "zero columns",
			table:     "test_table",
			rows:      3,
			cols:      0,
			wantQuery: "INSERT OR REPLACE INTO test_table VALUES (),(),()",
			wantErr:   false,
		},
		{
			name:    "negative rows",
			table:   "test_table",
			rows:    -1,
			cols:    3,
			wantErr: true,
		},
		{
			name:    "negative columns",
			table:   "test_table",
			rows:    3,
			cols:    -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockPreparer{}
			ctx := context.Background()

			_, err := prepareBatchInsert(ctx, mock, tt.table, tt.rows, tt.cols)

			if (err != nil) != tt.wantErr {
				t.Errorf("prepareBatchInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && !tt.wantErr {
				// Normalize whitespace for comparison
				gotQuery := strings.Join(strings.Fields(mock.lastQuery), " ")
				wantQuery := strings.Join(strings.Fields(tt.wantQuery), " ")

				if gotQuery != wantQuery {
					t.Errorf("prepareBatchInsert() query = %q, want %q", mock.lastQuery, tt.wantQuery)
				}
			}
		})
	}
}

func TestPrepareBatchInsertRealDatabase(t *testing.T) {
	// Create an in-memory SQLite database for testing
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer func() { _ = db.Close() }()

	// Create a test table
	ctx := context.Background()
	_, err = db.ExecContext(ctx, "CREATE TABLE test_table (id INTEGER, name TEXT, value INTEGER)")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Test batch insert preparation
	stmt, err := prepareBatchInsert(ctx, db, "test_table", 2, 3)
	if err != nil {
		t.Fatalf("prepareBatchInsert() error = %v", err)
	}
	defer func() { _ = stmt.Close() }()

	// Execute the statement with actual values
	_, err = stmt.ExecContext(ctx, 1, "first", 100, 2, "second", 200)
	if err != nil {
		t.Fatalf("Failed to execute batch insert: %v", err)
	}

	// Verify the data was inserted
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_table").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query count: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 rows, got %d", count)
	}
}

func TestPrepareBatchInsertWithTransaction(t *testing.T) {
	// Create an in-memory SQLite database for testing
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer func() { _ = db.Close() }()

	// Create a test table
	ctx := context.Background()
	_, err = db.ExecContext(ctx, "CREATE TABLE test_table (id INTEGER, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Test withTransaction
	err = withTransaction(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		stmt, err := prepareBatchInsert(ctx, tx, "test_table", 2, 2)
		if err != nil {
			return err
		}
		defer func() { _ = stmt.Close() }()

		_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, 1, "first", 2, "second")
		return err
	})

	if err != nil {
		t.Fatalf("withTransaction() error = %v", err)
	}

	// Verify the data was committed
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_table").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query count: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 rows, got %d", count)
	}
}

func TestWithTransactionRollback(t *testing.T) {
	// Create an in-memory SQLite database for testing
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer func() { _ = db.Close() }()

	// Create a test table
	ctx := context.Background()
	_, err = db.ExecContext(ctx, "CREATE TABLE test_table (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Insert initial data
	_, err = db.ExecContext(ctx, "INSERT INTO test_table VALUES (1, 'initial')")
	if err != nil {
		t.Fatalf("Failed to insert initial data: %v", err)
	}

	// Test withTransaction with an error (should rollback)
	err = withTransaction(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "INSERT INTO test_table VALUES (2, 'second')")
		if err != nil {
			return err
		}

		// Cause an error (duplicate primary key)
		_, err = tx.ExecContext(ctx, "INSERT INTO test_table VALUES (1, 'duplicate')")
		return err
	})

	if err == nil {
		t.Error("Expected error from duplicate primary key, got nil")
	}

	// Verify the second insert was rolled back
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_table").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query count: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 row (rollback successful), got %d", count)
	}
}
