# Retryable endpoints

We are building an API that allows a system to add credit to an account. Bread and butter Go, at least on the happy path. But we're building a resilient [distributed system](https://en.wikipedia.org/wiki/Distributed_computing), so we need to think carefully about what happens when things go wrong.

Imagine a client is using a _[transactional outbox](https://microservices.io/patterns/data/transactional-outbox.html)_. If you're not familiar, it's enough to know that it's a list of pending work stored in a database. A worker picks up an item from the outbox, calls our API and then marks the item as done if the call succeeds.

What happens if the worker calls our API successfully, and then crashes before it can mark the item as done? When the worker restarts, it'll pick up the item, call the API again, and we will now accidentally increase the credit on the account again.

What we need to do is make our API friendly for retries. A fancier term for this is **[idempotency](https://en.wikipedia.org/wiki/Idempotence#Computer_science_meaning)**: repeating the same logical operation has the same intended effect as doing it once. Retrying a request should not add more to the account, but two separate requests should still deposit.

Sometimes we can design our APIs to be idempotent, almost out of the box. An [HTTP PUT](https://en.wikipedia.org/wiki/HTTP#Idempotent_method)'s semantics mean you update a resource in place, and if you do the same call again, the state of the resource is the same. GET should also follow these retryable semantics.

For operations that musn't be repeated on retry we need some help: **idempotency keys**. Let’s use TDD to see how they work.

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

## Refactor

Our tests deserve some attention too. We're about to exercise several variations of topping up an account, and I don't fancy copying all that JSON and HTTP setup into each test.

Let's extract a helper that takes our handler and payload, makes the request and returns the response.

```go
func postTopUp(t testing.TB, handler http.Handler, topUp TopUpRequest) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(topUp)
	if err != nil {
		t.Fatalf("could not marshal top-up request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/top-up", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
}
```

The helper handles the boring detail, including checking the error from `json.Marshal` that we ignored earlier. It creates a fresh request each time, which will be useful when we call the handler more than once: reading a request body consumes it.

We can extract our assertions too.

```go
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
```

Remember `t.Helper()` tells Go to report failures at the line calling the helper, so we can find the failing assertion in our scenario. Including the account ID and response body in our error messages should also help us understand what went wrong.

Let's also give our fake a constructor. Tests shouldn't need to know that it stores balances in a map, or remember to initialise it.

```go
func NewInMemoryAccounts() *InMemoryAccounts {
	return &InMemoryAccounts{balances: make(map[string]int)}
}
```

Now our test can focus on the top-up and its effect on the account.

```go
func TestCreditAccount(t *testing.T) {
	t.Run("adds credit to an account", func(t *testing.T) {
		accounts := NewInMemoryAccounts()
		handler := RetryableEndpoint(accounts)

		topUp := TopUpRequest{
			AccountID:   "user-123",
			AmountPence: 1000,
		}

		response := postTopUp(t, handler, topUp)

		assertStatus(t, response, http.StatusOK)
		assertBalance(t, accounts, "user-123", 1000)
	})
}
```

Run the test again to check it still passes. We can now write the next scenario without repeating the request setup and assertion details.

## Write the test first

Our clients are asking if we can support idempotency keys. They will generate one per logical top-up, and they expect us to use that key to make the endpoint idempotent. In practice this means if they retry with the same idempotency key, we will not top up again. 

```go
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
```

## Try to run the test

For this to compile, you'll need to update `postTopUp` to accept the key, and update the first test to pass it. You can just send an empty one for now. 

You'll also need Go 1.27 or above to have access to the `uuid` package. 

```go
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
```

## Try and run the test
```
--- FAIL: TestCreditAccount (0.00s)
    --- FAIL: TestCreditAccount/uses_idempotency_keys_given_from_client (0.00s)
        endpoint_test.go:84: got balance 2000 pence for account "user-123", want 1000
```

As expected, our endpoint does not use the idempotency key, so it tops up twice.

## Write enough code to make it pass

We need to remember which keys we've already handled. A map will do for now.

```go
func RetryableEndpoint(accounts Accounts) http.Handler {
	handledKeys := make(map[string]bool)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var topUp TopUpRequest
		if err := json.NewDecoder(r.Body).Decode(&topUp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idempotencyKey := r.Header.Get("Idempotency-Key")

		if handledKeys[idempotencyKey] {
			w.WriteHeader(http.StatusOK)
			return
		}

		accounts.AddCredit(topUp.AccountID, topUp.AmountPence)
		handledKeys[idempotencyKey] = true
		w.WriteHeader(http.StatusOK)
	})
}
```

Our tests pass, but this simple implementation has some gaps:

- Empty keys are not handled well at all. After the first empty-key request, all others will be ignored! For this endpoint, we'll require a non-empty key.
- Checking the key and claiming it should be _atomic_: no other request should be able to slip between those steps. Two requests with the same key could both find it missing and both add credit before either records it. Accessing the map concurrently without synchronisation is also a data race.
- The map belongs to one handler instance. If we run several instances of our service, a retry could reach a different instance that hasn't seen the key and add credit again. Restarting the service also loses the keys. We'll need to share and persist this information to support horizontal scaling, which we'll explore later in the chapter.

Let's tackle these one at a time.

## Write the test first

This is the simplest one to make pass. To prepare, update the first test to pass in a UUID as the key rather than an empty string, so we don't end up with two failing tests. Then, let's write a test to check for the key is sent properly

```go
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
```

## Try to run the test

```
--- FAIL: TestCreditAccount (0.00s)
    --- FAIL: TestCreditAccount/bad_request_when_idempotency_key_is_missing (0.00s)
        endpoint_test.go:107: got status 200, want 400;

```

Fails as expected

## Write enough code to make it pass

We just need to add a bit of validation to the header after we've extracted it from the request.

```go
if idempotencyKey == "" {
    http.Error(w, "missing idempotency key", http.StatusBadRequest)
    return
}
```

The test will now pass. 
