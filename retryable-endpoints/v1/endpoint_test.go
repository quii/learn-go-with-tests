package retryableendpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"testing/synctest"
	"uuid"
)

type Accounts interface {
	AddCredit(accountID string, amountPence int)
	Balance(accountID string) int
}

type InMemoryAccounts struct {
	balances map[string]int
}

func NewInMemoryAccounts() *InMemoryAccounts {
	return &InMemoryAccounts{balances: make(map[string]int)}
}

func (a *InMemoryAccounts) AddCredit(accountID string, amountPence int) {
	a.balances[accountID] += amountPence
}

func (a *InMemoryAccounts) Balance(accountID string) int {
	return a.balances[accountID]
}

type PausingAccounts struct {
	Accounts
	pauseNext bool
	resume    <-chan struct{}
}

func (a *PausingAccounts) AddCredit(accountID string, amountPence int) {
	a.Accounts.AddCredit(accountID, amountPence)

	if a.pauseNext {
		a.pauseNext = false
		<-a.resume
	}
}

type TopUpRequest struct {
	AccountID   string `json:"account_id"`
	AmountPence int    `json:"amount_pence"`
}

type claimResult int

const (
	claimed claimResult = iota
	alreadyInProgress
	alreadyCompleted
)

type IdempotencyStore struct {
	mu       sync.Mutex
	requests map[string]claimResult
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{requests: make(map[string]claimResult)}
}

func (s *IdempotencyStore) Claim(key string) claimResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	if result, exists := s.requests[key]; exists {
		return result
	}

	s.requests[key] = alreadyInProgress
	return claimed
}

func (s *IdempotencyStore) Complete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requests[key] = alreadyCompleted
}

func RetryableEndpoint(accounts Accounts) http.Handler {
	store := NewIdempotencyStore()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var topUp TopUpRequest
		if err := json.NewDecoder(r.Body).Decode(&topUp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idempotencyKey := r.Header.Get("Idempotency-Key")
		if idempotencyKey == "" {
			http.Error(w, "missing idempotency key", http.StatusBadRequest)
			return
		}

		switch store.Claim(idempotencyKey) {
		case alreadyInProgress:
			http.Error(w, "top-up is still in progress", http.StatusConflict)
		case alreadyCompleted:
			w.WriteHeader(http.StatusOK)
		case claimed:
			accounts.AddCredit(topUp.AccountID, topUp.AmountPence)
			store.Complete(idempotencyKey)
			w.WriteHeader(http.StatusOK)
		}
	})
}

func TestCreditAccount(t *testing.T) {
	t.Run("adds credit to an account", func(t *testing.T) {
		accounts := NewInMemoryAccounts()
		handler := RetryableEndpoint(accounts)

		topUp := TopUpRequest{
			AccountID:   "user-123",
			AmountPence: 1000,
		}

		response := postTopUp(t, handler, topUp, uuid.New().String())

		assertStatus(t, response, http.StatusOK)
		assertBalance(t, accounts, "user-123", 1000)
	})

	t.Run("credits the account only once when a request is retried", func(t *testing.T) {
		accounts := NewInMemoryAccounts()
		handler := RetryableEndpoint(accounts)

		topUp := TopUpRequest{
			AccountID:   "user-123",
			AmountPence: 1000,
		}

		idempotencyKey := uuid.New().String()

		res1 := postTopUp(t, handler, topUp, idempotencyKey)
		assertStatus(t, res1, http.StatusOK)
		assertBalance(t, accounts, "user-123", 1000)

		res2 := postTopUp(t, handler, topUp, idempotencyKey)
		assertStatus(t, res2, http.StatusOK)
		assertBalance(t, accounts, "user-123", 1000)
	})

	t.Run("bad request when idempotency key is missing", func(t *testing.T) {
		accounts := NewInMemoryAccounts()
		handler := RetryableEndpoint(accounts)

		topUp := TopUpRequest{
			AccountID:   "user-123",
			AmountPence: 1000,
		}

		res := postTopUp(t, handler, topUp, "")
		assertStatus(t, res, http.StatusBadRequest)
		assertBalance(t, accounts, "user-123", 0)
	})

	t.Run("does not credit twice when a retry arrives before the first request completes", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			accounts := NewInMemoryAccounts()
			resume := make(chan struct{})
			handler := RetryableEndpoint(&PausingAccounts{
				Accounts:  accounts,
				pauseNext: true,
				resume:    resume,
			})
			topUp := TopUpRequest{
				AccountID:   "user-123",
				AmountPence: 1000,
			}
			idempotencyKey := uuid.New().String()
			request := newTopUpRequest(t, topUp, idempotencyKey)
			firstResponse := httptest.NewRecorder()

			go handler.ServeHTTP(firstResponse, request)
			synctest.Wait()

			// The credit has been added, but the first request hasn't finished.
			assertBalance(t, accounts, "user-123", 1000)

			// The client hasn't received confirmation, so it retries the same top-up.
			retryResponse := postTopUp(t, handler, topUp, idempotencyKey)

			close(resume)
			synctest.Wait()

			assertStatus(t, firstResponse, http.StatusOK)
			assertStatus(t, retryResponse, http.StatusConflict)
			assertBalance(t, accounts, "user-123", 1000)
		})
	})
}

func postTopUp(t testing.TB, handler http.Handler, topUp TopUpRequest, idempotencyKey string) *httptest.ResponseRecorder {
	t.Helper()

	req := newTopUpRequest(t, topUp, idempotencyKey)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
}

func newTopUpRequest(t testing.TB, topUp TopUpRequest, idempotencyKey string) *http.Request {
	t.Helper()

	payload, err := json.Marshal(topUp)
	if err != nil {
		t.Fatalf("could not marshal top-up request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/top-up", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", idempotencyKey)
	return req
}

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	if response.Code != want {
		t.Errorf("got status %d, want %d; response body: %s", response.Code, want, response.Body.String())
	}
}

func assertBalance(t testing.TB, accounts Accounts, accountID string, want int) {
	t.Helper()
	if got := accounts.Balance(accountID); got != want {
		t.Errorf("got balance %d pence for account %q, want %d", got, accountID, want)
	}
}
