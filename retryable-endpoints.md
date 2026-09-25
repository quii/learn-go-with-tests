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

## Write the test first

So far, our client retries after the first request has finished. But what if it hasn't received a response and retries while the first request is still running? The account might already have been credited.

We could start two goroutines and hope they hit the handler at just the wrong moment. I'd rather have a test that fails every time we have this bug, not just when my laptop is having a busy day.

Let's arrange this sequence deliberately:

1. The first request checks its key and adds credit to the account.
2. We pause it before `AddCredit` returns, so the handler hasn't recorded the key or sent a response yet.
3. A second request arrives with the same key.
4. We let the first request finish and check that the account was only credited once.

There's also a response to consider. Previously, we returned `200 OK` for a retry because the original request had finished. For an operation still in progress, we'll choose to return `409 Conflict`, allowing the caller to retry later. Waiting for the first request to finish is another option, but we'll use the fail-fast approach here.

We'll use [`testing/synctest`](revisiting-time-with-synctest.md) to control the overlap. Remember that `synctest.Test` creates a bubble containing our test and its goroutines. `synctest.Wait()` waits until all the other goroutines in that bubble have finished or are *durably blocked*, such as waiting to receive from a channel created inside the bubble.

Add `"testing/synctest"` to your imports, then add this subtest. We'll implement `PausingAccounts` and `newTopUpRequest` next.

```go
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
```

The two calls to `synctest.Wait()` do different jobs:

1. The first lets the handler reach the pause. We can then check the balance and send the retry while the first request is still in progress.
2. After we close `resume`, the second lets the first handler finish. We can then safely inspect its response and the final balance.

No sleeps required. We're controlling the order of events, rather than guessing how long they take.

## Write the minimal amount of code for the test to run and check the failing test output

We need a way to pause the account operation without adding test-specific behaviour to our handler. As in [Working Without Mocks](working-without-mocks.md#off-the-happy-path-with-decorators), we can decorate our fake.

```go
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
```

Embedding `Accounts` lets us keep using the fake's `Balance` method. We only change how `AddCredit` behaves:

1. Delegate to the fake to add the credit.
2. If this call should pause, set `pauseNext` to `false` so the retry won't pause too.
3. Wait for the test to close `resume` before returning to the handler.

Notice that we're pausing *after* changing the balance. The work has happened, but the handler hasn't yet recorded that fact. That's the gap we want to expose.

This decorator relies on our controlled ordering: the first `synctest.Wait()` lets the first request reach the channel receive before we start the second request. It isn't a general-purpose, concurrency-safe accounts implementation.

We also need to separate constructing a request from sending it. Our existing helper can call `t.Fatalf` if JSON encoding fails, and fatal test methods must be called from the test goroutine. Let's do that preparation before starting the handler's goroutine.

Extract the request construction into `newTopUpRequest`:

```go
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
```

Our existing tests can keep using `postTopUp`, which now delegates to the new helper:

```go
func postTopUp(t testing.TB, handler http.Handler, topUp TopUpRequest, idempotencyKey string) *httptest.ResponseRecorder {
	t.Helper()

	req := newTopUpRequest(t, topUp, idempotencyKey)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
}
```

Run the new test:

```sh
go test ./retryable-endpoints -run 'TestCreditAccount/does_not_credit_twice'
```

```text
--- FAIL: TestCreditAccount (0.00s)
    --- FAIL: TestCreditAccount/does_not_credit_twice_when_a_retry_arrives_before_the_first_request_completes (0.00s)
        endpoint_test.go:161: got status 200, want 409; response body:
        endpoint_test.go:162: got balance 2000 pence for account "user-123", want 1000
```

Both requests added credit. When the retry checked the map, the first request hadn't recorded its key yet, so the retry went ahead and returned `200 OK` too.

You can also run this with `-race`. This test fails its assertions without a race-detector report: we've deliberately ordered the accesses using `synctest.Wait()` and the channel. The requests overlap, but they don't access the maps at the same instant. We've exposed a check-then-act bug, which the race detector alone can't catch.

Our handler needs to distinguish a key that's *in progress* from one that's *completed*, and checking and claiming a key must happen atomically. That's our next step.

## Write enough code to make it pass

We now need to check and claim a key as one operation. Let's put the map and its mutex together in a type that handles this for us. The handler can then ask to claim a key without managing locks itself.

There are three possible outcomes when we try to claim a key:

```go
type claimResult int

const (
	claimed claimResult = iota
	alreadyInProgress
	alreadyCompleted
)
```

`claimed` means this caller gets to do the work. The other two results tell it why it shouldn't. We'll store the result that *subsequent* callers should receive, so `claimed` itself never goes into the map.

Add `"sync"` to your imports and create the store:

```go
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
```

The important part of `Claim` is the scope of the lock:

1. Acquire the mutex before looking up the key.
2. If the key exists, return its recorded result.
3. Otherwise, record it as in progress and return `claimed`.

The deferred unlock happens as the method returns. No other caller can check the map between our lookup and insertion, so two requests can't both claim the same key.

The mutex is only held while we inspect or update the map. **The claim remains recorded after the mutex is released.** That lets a retry discover that the work is in progress without waiting for the work to finish.

Now replace the handler's map and inline checks with our store:

```go
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
```

Run the tests again, including with the race detector:

```sh
go test ./retryable-endpoints -race -count=10
```

They should pass. When our first request pauses, its key is already recorded as in progress, so the retry gets `409 Conflict` without adding any credit. Once the first request finishes, subsequent retries get `200 OK`.

The handler assumes its `Accounts` dependency supports concurrent calls; how it achieves that is the dependency's responsibility. Our concern here is preventing two requests with the same key from performing the same top-up twice.

The idempotency store is still local to one handler instance; we haven't solved persistence or horizontal scaling yet.
