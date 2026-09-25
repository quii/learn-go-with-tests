package retryableendpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/synctest"
	"uuid"
)

func TestCreditAccount(t *testing.T) {
	t.Run("adds credit and returns the result", func(t *testing.T) {
		topUps := NewInMemoryTopUps()
		handler := RetryableEndpoint(topUps)
		request := TopUpRequest{AccountID: "user-123", AmountPence: 1000}

		response := postTopUp(t, handler, request, uuid.New().String())

		assertStatus(t, response, http.StatusOK)
		assertResult(t, readResult(t, response), TopUpResult{AccountID: "user-123", BalancePence: 1000})
		assertBalance(t, topUps, "user-123", 1000)
		if got := response.Header().Get("Content-Type"); got != "application/json" {
			t.Errorf("got Content-Type %q, want application/json", got)
		}
	})

	t.Run("retries can reach a different handler", func(t *testing.T) {
		topUps := NewInMemoryTopUps()
		firstHandler := RetryableEndpoint(topUps)
		secondHandler := RetryableEndpoint(topUps)
		request := TopUpRequest{AccountID: "user-123", AmountPence: 1000}
		key := uuid.New().String()

		first := postTopUp(t, firstHandler, request, key)
		retry := postTopUp(t, secondHandler, request, key)

		assertStatus(t, first, http.StatusOK)
		assertStatus(t, retry, http.StatusOK)
		if first.Body.String() != retry.Body.String() {
			t.Errorf("retry body %q differs from original %q", retry.Body.String(), first.Body.String())
		}
		assertBalance(t, topUps, "user-123", 1000)
	})

	t.Run("bad request when idempotency key is missing", func(t *testing.T) {
		topUps := NewInMemoryTopUps()
		response := postTopUp(t, RetryableEndpoint(topUps), TopUpRequest{AccountID: "user-123", AmountPence: 1000}, "")

		assertStatus(t, response, http.StatusBadRequest)
		assertBalance(t, topUps, "user-123", 0)
	})

	t.Run("retry succeeds before the first response reaches the caller", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			topUps := NewInMemoryTopUps()
			resume := make(chan struct{})
			firstHandler := RetryableEndpoint(&PausingTopUps{TopUps: topUps, resume: resume})
			secondHandler := RetryableEndpoint(topUps)
			request := TopUpRequest{AccountID: "user-123", AmountPence: 1000}
			key := uuid.New().String()
			firstRequest := newTopUpRequest(t, request, key)
			firstResponse := httptest.NewRecorder()

			go firstHandler.ServeHTTP(firstResponse, firstRequest)
			synctest.Wait()

			// Apply has completed, but the first handler hasn't received its result.
			assertBalance(t, topUps, "user-123", 1000)
			retry := postTopUp(t, secondHandler, request, key)
			close(resume)
			synctest.Wait()

			assertStatus(t, firstResponse, http.StatusOK)
			assertStatus(t, retry, http.StatusOK)
			assertResult(t, readResult(t, retry), readResult(t, firstResponse))
			assertBalance(t, topUps, "user-123", 1000)
		})
	})

	t.Run("reports an operation failure", func(t *testing.T) {
		topUps := NewInMemoryTopUps()
		handler := RetryableEndpoint(FailingTopUps{TopUps: topUps})
		response := postTopUp(t, handler, TopUpRequest{AccountID: "user-123", AmountPence: 1000}, uuid.New().String())
		assertStatus(t, response, http.StatusInternalServerError)
		assertBalance(t, topUps, "user-123", 0)
	})
}

// Pause after the operation has finished, modelling a delayed acknowledgement.
type PausingTopUps struct {
	TopUps
	resume <-chan struct{}
}

func (p *PausingTopUps) Apply(ctx context.Context, key string, request TopUpRequest) (TopUpResult, error) {
	result, err := p.TopUps.Apply(ctx, key, request)
	<-p.resume
	return result, err
}

type FailingTopUps struct {
	TopUps
}

func (f FailingTopUps) Apply(context.Context, string, TopUpRequest) (TopUpResult, error) {
	return TopUpResult{}, errors.New("operation unavailable")
}

func postTopUp(t testing.TB, handler http.Handler, topUp TopUpRequest, key string) *httptest.ResponseRecorder {
	t.Helper()
	request := newTopUpRequest(t, topUp, key)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func newTopUpRequest(t testing.TB, topUp TopUpRequest, key string) *http.Request {
	t.Helper()
	payload, err := json.Marshal(topUp)
	if err != nil {
		t.Fatalf("could not marshal top-up request: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/top-up", bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	if key != "" {
		request.Header.Set("Idempotency-Key", key)
	}
	return request
}

func readResult(t testing.TB, response *httptest.ResponseRecorder) TopUpResult {
	t.Helper()
	var result TopUpResult
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("could not decode top-up result: %v", err)
	}
	return result
}

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	if response.Code != want {
		t.Errorf("got status %d, want %d; response body: %s", response.Code, want, response.Body.String())
	}
}
