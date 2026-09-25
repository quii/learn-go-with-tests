# Retryable endpoints

We are building an API that allows a system to add credit to an account. Bread and butter Go, at least on the happy path. But we're building a resilient [distributed system](https://en.wikipedia.org/wiki/Distributed_computing), so we need to think carefully about what happens when things go wrong.

Imagine a client is using a _[transactional outbox](https://microservices.io/patterns/data/transactional-outbox.html)_. If you're not familiar, it's enough to know that it's a list of pending work stored in a database. A worker picks up an item from the outbox, calls our API and then marks the item as done if the call succeeds.

What happens if the worker calls our API successfully, and then crashes before it can mark the item as done? When the worker restarts, it'll pick up the item, call the API again, and we will now accidentally increase the credit on the account again.

What we need to do is make our API friendly for retries. A fancier term for this is **[idempotency](https://en.wikipedia.org/wiki/Idempotence#Computer_science_meaning)**: repeating the same logical operation has the same intended effect as doing it once. Retrying a request should not add more to the account, but two separate requests should still deposit.

Sometimes we can design our APIs to be idempotent, almost out of the box. An [HTTP PUT](https://en.wikipedia.org/wiki/HTTP#Idempotent_method)'s semantics mean you update a resource in place, and if you do the same call again, the state of the resource is the same. GET should also follow these retryable semantics.

For operations like sending an email, we need some help: **idempotency keys**. Let’s use TDD to see how they work.

## Write the test first

We'll start with the happy path to get the scaffolding in. We'll be testing an `http.Handler`.

As discussed in [Working Without Mocks](working-without-mocks.md), we prefer to model test-doubles as fakes rather than spies. We'll model the account with an interface so we can pass in a fake to our HTTP handler. We can then call the endpoint, and verify the user's account balance.

```go
func TestCreditAccount(t *testing.T) {
	t.Run("adds credit to an account", func(t *testing.T) {
		accounts := &InMemoryAccounts{balances: make(map[string]int)}
		handler := RetryableEndpoint(accounts)

		topUp := TopUpRequest{
			AccountID:   "user-123",
			AmountPence: 1000,
		}

		topUpPayload, _ := json.Marshal(topUp)

		req := httptest.NewRequest("POST", "/top-up", bytes.NewReader(topUpPayload))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		if balance := accounts.Balance("user-123"); balance != 1000 {
			t.Errorf("got balance %d pence, want 1000", balance)
		}
	})
}
```

## Write the minimal amount of code for the test to run and check the failing test output

Let's add the account interface, its in-memory implementation, the request type and an empty handler so we can run our test.

```go
type Accounts interface {
	AddCredit(accountID string, amountPence int)
	Balance(accountID string) int
}

type InMemoryAccounts struct {
	balances map[string]int
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
		// TODO: add credit to the account
	})
}
```

If you run this test, you'll see it fails because the balance of the account has not been updated.

## Write enough code to make it pass

We just need to parse the JSON into the `TopUpRequest` and call the service.

```go
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
```

For brevity, we'll leave testing invalid JSON out of this example.

It should now pass.
