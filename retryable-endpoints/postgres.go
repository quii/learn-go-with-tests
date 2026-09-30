package retryableendpoints

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PostgresTopUps struct {
	db *sql.DB
}

func NewPostgresTopUps(db *sql.DB) *PostgresTopUps {
	return &PostgresTopUps{db: db}
}

func (s *PostgresTopUps) Apply(ctx context.Context, key string, request TopUpRequest) (TopUpResult, error) {
	if key == "" {
		return TopUpResult{}, ErrMissingKey
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return TopUpResult{}, fmt.Errorf("begin top-up: %w", err)
	}
	defer tx.Rollback()

	claim, err := tx.ExecContext(ctx, `
		INSERT INTO top_up_results (idempotency_key) VALUES ($1)
		ON CONFLICT (idempotency_key) DO NOTHING`, key)
	if err != nil {
		return TopUpResult{}, fmt.Errorf("claim top-up: %w", err)
	}
	inserted, err := claim.RowsAffected()
	if err != nil {
		return TopUpResult{}, fmt.Errorf("read claim outcome: %w", err)
	}

	var result TopUpResult
	if inserted == 0 {
		// At READ COMMITTED, this new statement sees the competing transaction's
		// committed result, even if our INSERT had to wait for it.
		err = tx.QueryRowContext(ctx, `
			SELECT account_id, balance_pence FROM top_up_results
			WHERE idempotency_key = $1`, key).Scan(&result.AccountID, &result.BalancePence)
		if err != nil {
			return TopUpResult{}, fmt.Errorf("read top-up result: %w", err)
		}
		return result, nil
	}

	result.AccountID = request.AccountID
	err = tx.QueryRowContext(ctx, `
		INSERT INTO accounts (account_id, balance_pence) VALUES ($1, $2)
		ON CONFLICT (account_id) DO UPDATE
		SET balance_pence = accounts.balance_pence + EXCLUDED.balance_pence
		RETURNING balance_pence`, request.AccountID, request.AmountPence).Scan(&result.BalancePence)
	if err != nil {
		return TopUpResult{}, fmt.Errorf("add credit: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE top_up_results SET account_id = $2, balance_pence = $3
		WHERE idempotency_key = $1`, key, result.AccountID, result.BalancePence)
	if err != nil {
		return TopUpResult{}, fmt.Errorf("record top-up result: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return TopUpResult{}, fmt.Errorf("commit top-up: %w", err)
	}
	return result, nil
}

func (s *PostgresTopUps) Balance(ctx context.Context, accountID string) (int, error) {
	var balance int
	err := s.db.QueryRowContext(ctx, `
		SELECT balance_pence FROM accounts WHERE account_id = $1`, accountID).Scan(&balance)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("read balance: %w", err)
	}
	return balance, nil
}
