package retryableendpoints

import (
	"errors"
	"testing"
)

// New must provide isolated state for each scenario. A database adapter can
// register its cleanup with t.Cleanup.
type TopUpsContract struct {
	New func(t testing.TB) TopUps
}

func (c TopUpsContract) Test(t *testing.T) {
	t.Run("credits an account and returns its balance", func(t *testing.T) {
		topUps := c.New(t)
		request := TopUpRequest{AccountID: "user-123", AmountPence: 1000}

		got := applyTopUp(t, topUps, "first", request)

		assertResult(t, got, TopUpResult{AccountID: "user-123", BalancePence: 1000})
		assertBalance(t, topUps, "user-123", 1000)
		assertBalance(t, topUps, "another-user", 0)
	})

	t.Run("replays the original result even after another top-up", func(t *testing.T) {
		topUps := c.New(t)
		request := TopUpRequest{AccountID: "user-123", AmountPence: 1000}
		first := applyTopUp(t, topUps, "first", request)

		// Identical payload, different key: a genuinely new top-up.
		second := applyTopUp(t, topUps, "second", request)
		assertResult(t, second, TopUpResult{AccountID: "user-123", BalancePence: 2000})

		replayed := applyTopUp(t, topUps, "first", request)
		assertResult(t, replayed, first)
		assertBalance(t, topUps, "user-123", 2000)
	})

	t.Run("rejects a missing key without adding credit", func(t *testing.T) {
		topUps := c.New(t)
		_, err := topUps.Apply(t.Context(), "", TopUpRequest{AccountID: "user-123", AmountPence: 1000})
		if !errors.Is(err, ErrMissingKey) {
			t.Errorf("got error %v, want %v", err, ErrMissingKey)
		}
		assertBalance(t, topUps, "user-123", 0)
	})

	t.Run("concurrent retries all return the same result", func(t *testing.T) {
		topUps := c.New(t)
		request := TopUpRequest{AccountID: "user-123", AmountPence: 1000}
		const attempts = 10
		ctx := t.Context()
		start := make(chan struct{})
		type outcome struct {
			result TopUpResult
			err    error
		}
		outcomes := make(chan outcome, attempts)
		for range attempts {
			go func() {
				<-start
				result, err := topUps.Apply(ctx, "same-top-up", request)
				outcomes <- outcome{result: result, err: err}
			}()
		}
		close(start)

		for range attempts {
			got := <-outcomes
			if got.err != nil {
				t.Errorf("could not apply top-up: %v", got.err)
				continue
			}
			assertResult(t, got.result, TopUpResult{AccountID: "user-123", BalancePence: 1000})
		}
		assertBalance(t, topUps, "user-123", 1000)
	})

	t.Run("concurrent distinct top-ups do not lose credit", func(t *testing.T) {
		topUps := c.New(t)
		ctx := t.Context()
		keys := []string{"first", "second", "third"}
		start := make(chan struct{})
		errors := make(chan error, len(keys))
		for _, key := range keys {
			go func() {
				<-start
				_, err := topUps.Apply(ctx, key, TopUpRequest{AccountID: "user-123", AmountPence: 1000})
				errors <- err
			}()
		}
		close(start)
		for range keys {
			if err := <-errors; err != nil {
				t.Errorf("could not apply top-up: %v", err)
			}
		}
		assertBalance(t, topUps, "user-123", 3000)
	})
}

func TestInMemoryTopUps(t *testing.T) {
	TopUpsContract{New: func(t testing.TB) TopUps {
		return NewInMemoryTopUps()
	}}.Test(t)
}

func applyTopUp(t testing.TB, topUps TopUps, key string, request TopUpRequest) TopUpResult {
	t.Helper()
	result, err := topUps.Apply(t.Context(), key, request)
	if err != nil {
		t.Fatalf("could not apply top-up: %v", err)
	}
	return result
}

func assertResult(t testing.TB, got, want TopUpResult) {
	t.Helper()
	if got != want {
		t.Errorf("got result %+v, want %+v", got, want)
	}
}

func assertBalance(t testing.TB, topUps TopUps, accountID string, want int) {
	t.Helper()
	got, err := topUps.Balance(t.Context(), accountID)
	if err != nil {
		t.Fatalf("could not read balance: %v", err)
	}
	if got != want {
		t.Errorf("got balance %d pence for account %q, want %d", got, accountID, want)
	}
}
