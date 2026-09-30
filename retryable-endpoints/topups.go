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

// Idempotent wraps a plain operation so repeated calls with the same key
// run it once and replay the stored result thereafter.
type Idempotent[Req, Res any] struct {
	mu      sync.Mutex
	results map[string]Res
	work    func(Req) (Res, error)
}

func NewIdempotent[Req, Res any](work func(Req) (Res, error)) *Idempotent[Req, Res] {
	return &Idempotent[Req, Res]{results: make(map[string]Res), work: work}
}

func (s *Idempotent[Req, Res]) Apply(ctx context.Context, key string, request Req) (Res, error) {
	var zero Res
	if key == "" {
		return zero, ErrMissingKey
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return zero, err
	}
	if result, exists := s.results[key]; exists {
		return result, nil
	}

	result, err := s.work(request)
	if err != nil {
		return zero, err
	}
	s.results[key] = result
	return result, nil
}

// Query runs read under the same lock Apply uses, so a caller can never
// observe work's side effects before Apply has finished recording its result.
func (s *Idempotent[Req, Res]) Query(read func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	read()
}

type InMemoryTopUps struct {
	*Idempotent[TopUpRequest, TopUpResult]
	balances map[string]int
}

func NewInMemoryTopUps() *InMemoryTopUps {
	s := &InMemoryTopUps{balances: make(map[string]int)}
	s.Idempotent = NewIdempotent(func(request TopUpRequest) (TopUpResult, error) {
		s.balances[request.AccountID] += request.AmountPence
		return TopUpResult{
			AccountID:    request.AccountID,
			BalancePence: s.balances[request.AccountID],
		}, nil
	})
	return s
}

func (s *InMemoryTopUps) Balance(ctx context.Context, accountID string) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	var balance int
	s.Query(func() { balance = s.balances[accountID] })
	return balance, nil
}
