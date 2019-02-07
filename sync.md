# Sync (WIP)

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/master/sync)**

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
    wantedCount := 1000
    counter := Counter{}

    var wg sync.WaitGroup
    wg.Add(wantedCount)

    for i:=0; i<wantedCount; i++ {
        go func(w *sync.WaitGroup) {
            counter.Inc()
            w.Done()
        }(&wg)
    }
    wg.Wait()

    assertCounter(t, counter, wantedCount)
})
```

This will loop through our `wantedCount` and fire a go routine to call `counter.Inc()`. 

We are using [`sync.WaitGroup`](https://golang.org/pkg/sync/#WaitGroup) which is a convenient way of synchronising concurrent processes.

> A WaitGroup waits for a collection of goroutines to finish. The main goroutine calls Add to set the number of goroutines to wait for. Then each of the goroutines runs and calls Done when finished. At the same time, Wait can be used to block until all goroutines have finished.

By waiting for `wg.Wait()` to finish before making our assertions we can be sure all of our go routines have attempted to `Inc` the `Counter`,

## Try to run the test

```
=== RUN   TestCounter/it_runs_safely_in_a_concurrent_envionment
--- FAIL: TestCounter (0.00s)
    --- FAIL: TestCounter/it_runs_safely_in_a_concurrent_envionment (0.00s)
    	sync_test.go:26: got 939, want 1000
FAIL
```

The test will _probably_ fail with a different number, but nonetheless it demonstrates it does not work when multiple go routines are trying to work with it.

## Write enough code to make it pass

A simple solution is to add a lock to our `Counter`, a [`Mutex`](https://golang.org/pkg/sync/#Mutex)

>A Mutex is a mutual exclusion lock. The zero value for a Mutex is an unlocked mutex.

```go
type Counter struct {
	value int
	lock sync.Mutex
}

func (c *Counter) Inc() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.value++
}
```

What this means is any go routine calling `Inc` will acquire the lock on `Counter` if they are first. All the other go routines will have to wait for it to be `Unlock`ed before getting access. 

## I've seen other examples where the `sync.Mutex` is embedded into the struct. 

You may see examples like this 

```go
type Counter struct {
	value int
	sync.Mutex
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

![This seems like a really bad idea](https://i.imgur.com/SWYNpwm.png)

_This seems like a really bad idea_

## Wrapping up

We've covered a few things from the [sync package](https://golang.org/pkg/sync/)

- `Mutex` allows us to add locks to our data
- `Waitgroup` is a means of waiting for go routines to finish jobs

### When to use locks over channels and go routines?

[We've previously covered go routines in the first concurrency chapter](concurrency.md) which let us write safe concurrent code so why would you use locks?   
[The go wiki has a page dedicated to this topic; Mutex Or Channel](https://github.com/golang/go/wiki/MutexOrChannel)

> A common Go newbie mistake is to over-use channels and goroutines just because it's possible, and/or because it's fun. Don't be afraid to use a sync.Mutex if that fits your problem best. Go is pragmatic in letting you use the tools that solve your problem best and not forcing you into one style of code.

Paraphrasing:

- **Use channels when passing ownership of data** 
- **Use mutexes for managing state**
