package retryableendpoints

import (
	"context"
	"errors"
	"sync"
)

type TopUpRequest struct {
	AccountID   string `json:"account_id"`
	AmountPence int    `json:"amount_pence"`
}

type TopUpResult struct {
	AccountID    string `json:"account_id"`
	BalancePence int    `json:"balance_pence"`
}

var ErrMissingKey = errors.New("missing idempotency key")

type TopUps interface {
	// Apply credits the account once per key and replays the original result.
	// Retries must use the same request. Concurrent calls with the same key
	// wait for the result rather than returning an in-progress response.
	Apply(ctx context.Context, key string, request TopUpRequest) (TopUpResult, error)
	Balance(ctx context.Context, accountID string) (int, error)
}

type InMemoryTopUps struct {
	mu       sync.Mutex
	balances map[string]int
	results  map[string]TopUpResult
}

func NewInMemoryTopUps() *InMemoryTopUps {
	return &InMemoryTopUps{
		balances: make(map[string]int),
		results:  make(map[string]TopUpResult),
	}
}

func (s *InMemoryTopUps) Apply(ctx context.Context, key string, request TopUpRequest) (TopUpResult, error) {
	if key == "" {
		return TopUpResult{}, ErrMissingKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return TopUpResult{}, err
	}
	if result, exists := s.results[key]; exists {
		return result, nil
	}

	s.balances[request.AccountID] += request.AmountPence
	result := TopUpResult{
		AccountID:    request.AccountID,
		BalancePence: s.balances[request.AccountID],
	}
	s.results[key] = result
	return result, nil
}

func (s *InMemoryTopUps) Balance(ctx context.Context, accountID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return 0, err
	}
	return s.balances[accountID], nil
}
