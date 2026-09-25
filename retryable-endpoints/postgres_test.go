package retryableendpoints

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresTopUps(t *testing.T) {
	if testing.Short() {
		t.Skip("Postgres integration test requires Docker")
	}
	db, connectionString := newPostgresDB(t)

	// The contract runs scenarios sequentially, so truncating between them
	// provides isolated state without starting a container for every test.
	TopUpsContract{New: func(t testing.TB) TopUps {
		resetPostgres(t, db)
		return NewPostgresTopUps(db)
	}}.Test(t)

	t.Run("replays through a new adapter and connection pool", func(t *testing.T) {
		resetPostgres(t, db)
		request := TopUpRequest{AccountID: "user-123", AmountPence: 1000}
		first := applyTopUp(t, NewPostgresTopUps(db), "first", request)

		otherDB := openPostgres(t, connectionString)
		other := NewPostgresTopUps(otherDB)
		applyTopUp(t, other, "second", request)
		assertResult(t, applyTopUp(t, other, "first", request), first)
		assertBalance(t, other, "user-123", 2000)
	})

	t.Run("failure to record the result rolls back the credit and key", func(t *testing.T) {
		resetPostgres(t, db)
		// This test-only constraint fails the result UPDATE, after credit was added.
		_, err := db.ExecContext(t.Context(), `ALTER TABLE top_up_results
			ADD CONSTRAINT reject_test_result CHECK (balance_pence <> 1234)`)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if _, err := db.Exec(`ALTER TABLE top_up_results DROP CONSTRAINT IF EXISTS reject_test_result`); err != nil {
				t.Error(err)
			}
		})

		topUps := NewPostgresTopUps(db)
		request := TopUpRequest{AccountID: "user-123", AmountPence: 1234}
		_, err = topUps.Apply(t.Context(), "first", request)
		if err == nil || !strings.Contains(err.Error(), "reject_test_result") {
			t.Fatalf("expected result constraint failure, got %v", err)
		}
		assertBalance(t, topUps, "user-123", 0)

		if _, err := db.ExecContext(t.Context(), `ALTER TABLE top_up_results DROP CONSTRAINT reject_test_result`); err != nil {
			t.Fatal(err)
		}
		// Retrying the same operation now succeeds: the failed transaction didn't
		// leave a claimed key or any credit behind.
		result := applyTopUp(t, topUps, "first", request)
		assertResult(t, result, TopUpResult{AccountID: "user-123", BalancePence: 1234})
		assertBalance(t, topUps, "user-123", 1234)
	})
}

func newPostgresDB(t *testing.T) (*sql.DB, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 90*time.Second)
	defer cancel()

	container, err := postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithDatabase("topups"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, container)
	if err != nil {
		t.Fatalf("start Postgres: %v", err)
	}
	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	db := openPostgres(t, connectionString)

	// Schema belongs in test setup for this example; no migration framework needed.
	_, err = db.ExecContext(ctx, `
		CREATE TABLE accounts (
			account_id TEXT PRIMARY KEY,
			balance_pence BIGINT NOT NULL
		);
		CREATE TABLE top_up_results (
			idempotency_key TEXT PRIMARY KEY,
			account_id TEXT,
			balance_pence BIGINT
		);`)
	if err != nil {
		t.Fatalf("create tables: %v", err)
	}
	return db, connectionString
}

func openPostgres(t testing.TB, connectionString string) *sql.DB {
	t.Helper()
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	return db
}

func resetPostgres(t testing.TB, db *sql.DB) {
	t.Helper()
	if _, err := db.ExecContext(t.Context(), `TRUNCATE accounts, top_up_results`); err != nil {
		t.Fatal(err)
	}
}
