# Locks (WIP)

We want to make a counter which is safe to use concurrently. 

We'll start with an unsafe counter with a test, see if we can exercise it's unsafeness via a test and then fix it using locks 

## Write the test first

We want our API to give us a method to increment the counter and then retrieve its value. 

```go
func TestCounter(t *testing.T) {

	t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {

		counter := Counter{}
		counter.Inc()
		counter.Inc()
		counter.Inc()

		if counter.Value != 3 {			
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
./sync_test.go:15:10: counter.Inc undefined (type Counter has no field or method Inc)
./sync_test.go:16:10: counter.Inc undefined (type Counter has no field or method Inc)
./sync_test.go:18:13: counter.Value undefined (type Counter has no field or method Value)
./sync_test.go:19:39: counter.Value undefined (type Counter has no field or method Value)
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

There's not a lot to refactor but given we're going to write more tests around `Counter` try writing a small assertion function `assertCount` so the test reads a bit clearer.

```go
t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
    counter := Counter{}
    counter.Inc()
    counter.Inc()
    counter.Inc()
    
    assertCounter(t, counter, 3)
})
```

## Next steps

That was easy enough but now we have a requirement that it must be safe to use in a concurrent environment. We will need to write a failing test to exercise this. 

## Write the test first

```go
	t.Run("it runs safely concurrently", func(t *testing.T) {
		wantedCount := 10
		counter := Counter{}

		for i:=0; i<wantedCount; i++ {
			go func() {
				counter.Inc()
			}()
		}

		assertCounter(t, counter, wantedCount)
	})
```

This will loop through our `wantedCount` and fire a go routine to call `counter.Inc()`. 

## Try to run the test

```
=== RUN   TestCounter/it_runs_safely_in_a_concurrent_envionment
--- FAIL: TestCounter (0.00s)
    --- FAIL: TestCounter/it_runs_safely_in_a_concurrent_envionment (0.00s)
    	sync_test.go:26: got 6, want 10
FAIL
```

The test will _probably_ fail with a different number, but nonetheless it demonstrates it does not work when multiple go routines are trying to work with it.

## Write enough code to make it pass

A simple solution is to add a lock to our `Counter`, a [`Mutex`](https://golang.org/pkg/sync/#Mutex)

>A Mutex is a mutual exclusion lock. The zero value for a Mutex is an unlocked mutex.

## Refactor
