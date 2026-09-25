package retryableendpoints

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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

type TopUpRequest struct {
	AccountID   string `json:"account_id"`
	AmountPence int    `json:"amount_pence"`
}

func RetryableEndpoint(accounts Accounts) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var topUp TopUpRequest
		if err := json.NewDecoder(r.Body).Decode(&topUp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		accounts.AddCredit(topUp.AccountID, topUp.AmountPence)
		w.WriteHeader(http.StatusOK)
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

		response := postTopUp(t, handler, topUp, "")

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
}

func postTopUp(t testing.TB, handler http.Handler, topUp TopUpRequest, idempotencyKey string) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(topUp)
	if err != nil {
		t.Fatalf("could not marshal top-up request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/top-up", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", idempotencyKey)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
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
