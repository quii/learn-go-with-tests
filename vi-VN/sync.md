# Sync

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/main/sync)**

We want to make a counter which is safe to use concurrently.

We'll start with an unsafe counter and verify its behaviour works in a single-threaded environment.

Then we'll exercise it's unsafeness, with multiple goroutines trying to use the counter via a test, and fix it.

## Write the test first

We want our API to give us a method to increment the counter and then retrieve its value.

```go
func TestCounter(t *testing.T) {
	t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
		counter := Counter{}
		counter.Inc()
		counter.Inc()
		counter.Inc()

		if counter.Value() != 3 {
			t.Errorf("got %d, want %d", counter.Value(), 3)
		}
	})
}
```

## Try to run the test

```
./sync_test.go:9:14: undefined: Counter
```

## Write the minimal amount of code for the test to run and check the failing test output

Let's define `Counter`.

```go
type Counter struct {
}
```

Try again and it fails with the following

```
./sync_test.go:14:10: counter.Inc undefined (type Counter has no field or method Inc)
./sync_test.go:18:13: counter.Value undefined (type Counter has no field or method Value)
```

So to finally make the test run we can define those methods

```go
func (c *Counter) Inc() {

}

func (c *Counter) Value() int {
	return 0
}
```

It should now run and fail

```
=== RUN   TestCounter
=== RUN   TestCounter/incrementing_the_counter_3_times_leaves_it_at_3
--- FAIL: TestCounter (0.00s)
    --- FAIL: TestCounter/incrementing_the_counter_3_times_leaves_it_at_3 (0.00s)
    	sync_test.go:27: got 0, want 3
```

## Write enough code to make it pass

This should be trivial for Go experts like us. We need to keep some state for the counter in our datatype and then increment it on every `Inc` call

```go
type Counter struct {
	value int
}

func (c *Counter) Inc() {
	c.value++
}

func (c *Counter) Value() int {
	return c.value
}
```

## Refactor

There's not a lot to refactor but given we're going to write more tests around `Counter` we'll write a small assertion function `assertCount` so the test reads a bit clearer.

```go
t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
	counter := Counter{}
	counter.Inc()
	counter.Inc()
	counter.Inc()

	assertCounter(t, counter, 3)
})
```
```go
func assertCounter(t testing.TB, got Counter, want int) {
	t.Helper()
	if got.Value() != want {
		t.Errorf("got %d, want %d", got.Value(), want)
	}
}
```

## Next steps

That was easy enough but now we have a requirement that it must be safe to use in a concurrent environment. We will need to write a failing test to exercise this.

## Write the test first

```go
t.Run("it runs safely concurrently", func(t *testing.T) {
	wantedCount := 1000
	counter := Counter{}

	var wg sync.WaitGroup
	wg.Add(wantedCount)

	for i := 0; i < wantedCount; i++ {
		go func() {
			counter.Inc()
			wg.Done()
		}()
	}
	wg.Wait()

	assertCounter(t, counter, wantedCount)
})
```

This will loop through our `wantedCount` and fire a goroutine to call `counter.Inc()`.

We are using [`sync.WaitGroup`](https://golang.org/pkg/sync/#WaitGroup) which is a convenient way of synchronising concurrent processes.

> A WaitGroup waits for a collection of goroutines to finish. The main goroutine calls Add to set the number of goroutines to wait for. Then each of the goroutines runs and calls Done when finished. At the same time, Wait can be used to block until all goroutines have finished.

By waiting for `wg.Wait()` to finish before making our assertions we can be sure all of our goroutines have attempted to `Inc` the `Counter`.

## Try to run the test

```
=== RUN   TestCounter/it_runs_safely_in_a_concurrent_envionment
--- FAIL: TestCounter (0.00s)
    --- FAIL: TestCounter/it_runs_safely_in_a_concurrent_envionment (0.00s)
    	sync_test.go:26: got 939, want 1000
FAIL
```

The test will _probably_ fail with a different number, but nonetheless it demonstrates it does not work when multiple goroutines are trying to mutate the value of the counter at the same time.

## Write enough code to make it pass

A simple solution is to add a lock to our `Counter`, ensuring only one goroutine can increment the counter at a time. Go's [`Mutex`](https://golang.org/pkg/sync/#Mutex) provides such a lock:

>A Mutex is a mutual exclusion lock. The zero value for a Mutex is an unlocked mutex.

```go
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}
```

What this means is any goroutine calling `Inc` will acquire the lock on `Counter` if they are first. All the other goroutines will have to wait for it to be `Unlock`ed before getting access.

If you now re-run the test it should now pass because each goroutine has to wait its turn before making a change.

## I've seen other examples where the `sync.Mutex` is embedded into the struct.

You may see examples like this

```go
type Counter struct {
	sync.Mutex
	value int
}
```

It can be argued that it can make the code a bit more elegant.

```go
func (c *Counter) Inc() {
	c.Lock()
	defer c.Unlock()
	c.value++
}
```

This _looks_ nice but while programming is a hugely subjective discipline, this is **bad and wrong**.

Sometimes people forget that embedding types means the methods of that type becomes _part of the public interface_; and you often will not want that. Remember that we should be very careful with our public APIs, the moment we make something public is the moment other code can couple themselves to it. We always want to avoid unnecessary coupling.

Exposing `Lock` and `Unlock` is at best confusing but at worst potentially very harmful to your software if callers of your type start calling these methods.

![Showing how a user of this API can wrongly change the state of the lock](https://i.imgur.com/SWYNpwm.png)

_This seems like a really bad idea_

## Copying mutexes

Our test passes but our code is still a bit dangerous

If you run `go vet` on your code you should get an error like the following

```
sync/v2/sync_test.go:16: call of assertCounter copies lock value: v1.Counter contains sync.Mutex
sync/v2/sync_test.go:39: assertCounter passes lock by value: v1.Counter contains sync.Mutex
```

A look at the documentation of [`sync.Mutex`](https://golang.org/pkg/sync/#Mutex) tells us why

> A Mutex must not be copied after first use.

When we pass our `Counter` (by value) to `assertCounter` it will try and create a copy of the mutex.

To solve this we should pass in a pointer to our `Counter` instead, so change the signature of `assertCounter`

```go
func assertCounter(t testing.TB, got *Counter, want int)
```

Our tests will no longer compile because we are trying to pass in a `Counter` rather than a `*Counter`. To solve this I prefer to create a constructor which shows readers of your API that it would be better to not initialise the type yourself.

```go
func NewCounter() *Counter {
	return &Counter{}
}
```

Use this function in your tests when initialising `Counter`.

## Wrapping up

We've covered a few things from the [sync package](https://golang.org/pkg/sync/)

- `Mutex` allows us to add locks to our data
- `WaitGroup` is a means of waiting for goroutines to finish jobs

### When to use locks over channels and goroutines?

[We've previously covered goroutines in the first concurrency chapter](concurrency.md) which let us write safe concurrent code so why would you use locks?
[The go wiki has a page dedicated to this topic; Mutex Or Channel](https://github.com/golang/go/wiki/MutexOrChannel)

> A common Go newbie mistake is to over-use channels and goroutines just because it's possible, and/or because it's fun. Don't be afraid to use a sync.Mutex if that fits your problem best. Go is pragmatic in letting you use the tools that solve your problem best and not forcing you into one style of code.

Paraphrasing:

- **Use channels when passing ownership of data**
- **Use mutexes for managing state**

### go vet

Remember to use go vet in your build scripts as it can alert you to some subtle bugs in your code before they hit your poor users.

### Don't use embedding because it's convenient

- Think about the effect embedding has on your public API.
- Do you _really_ want to expose these methods and have people coupling their own code to them?
- With respect to mutexes, this could be potentially disastrous in very unpredictable and weird ways, imagine some nefarious code unlocking a mutex when it shouldn't be; this would cause some very strange bugs that will be hard to track down.
