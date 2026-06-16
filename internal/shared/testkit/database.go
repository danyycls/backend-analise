package testkit

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type mockRows struct{}

func (m *mockRows) Close()                                       {}
func (m *mockRows) Err() error                                   { return nil }
func (m *mockRows) CommandTag() pgconn.CommandTag                { return pgconn.CommandTag{} }
func (m *mockRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (m *mockRows) Next() bool                                   { return false }
func (m *mockRows) Scan(dest ...any) error                       { return nil }
func (m *mockRows) Values() ([]any, error)                       { return nil, nil }
func (m *mockRows) RawValues() [][]byte                          { return nil }
func (m *mockRows) Conn() *pgx.Conn                              { return nil }

type mockRow struct{}

func (m *mockRow) Scan(dest ...any) error { return pgx.ErrNoRows }

type MockDB struct{}

func (m *MockDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return &mockRows{}, nil
}
func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return &mockRow{}
}
func (m *MockDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (m *MockDB) Begin(ctx context.Context) (pgx.Tx, error) { return nil, nil }
func (m *MockDB) Ping(ctx context.Context) error            { return nil }
