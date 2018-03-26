# Pointers (WIP)

We learned about structs in the last section which let us capture a number of values related around a concept. 

At some point you may wish to use structs to manage state, exposing methods to let users change the state in a way that you can control.

**Fintech loves Go** and uhhh bitcoins? so let's show what an amazing banking system we can make. 

Let's make a `Wallet` struct which let's us deposit `Bitcoin`

## Write the test first

```go
func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()
	want := 10

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

In the previous example we accessed fields directly with the field name, however in our _very secure wallet_ we dont want to expose our inner state to the rest of the world. We want to control access via methods.

## Try and run the test

`./wallet_test.go:7:12: undefined: Wallet`

## Write the minimal amount of code for the test to run and check the failing test output

The compiler doesn't know what a `Wallet` is so let's tell it.

```go
type Wallet struct {
	balance int
}
```

You may have noticed that `balance` starts with a lowercase letter. In Go if a name of a value starts with a lowercase that means it is private outside the package it was defined in. 

This means that we can control our state and anyone importing our package can only use the methods we decide to expose

Now we've made our wallet, try and run the test again

```go
./wallet_test.go:9:8: wallet.Deposit undefined (type Wallet has no field or method Deposit)
./wallet_test.go:11:15: wallet.Balance undefined (type Wallet has no field or method Balance, but does have balance)
```

As expected, we need to define these methods to make the test pass. 

Remember to only do enough to make the tests run. We need to make sure our test fails correctly with a clear error message

```go
func (w Wallet) Deposit(amount int) {

}

func (w Wallet) Balance() int {
	return 0
}
```

If this syntax is unfamiliar go back and read the structs section. 

The tests should now compile and run

`wallet_test.go:15: got 0 want 10`

## Write enough code to make it pass

Remember we can access the internal `balance` field in the struct using the "receiver" variable.

```go
func (w Wallet) Deposit(amount int) {
	w.balance += amount
}

func (w Wallet) Balance() int {
	return w.balance
}
```

With our career in fintech secured, run our tests and bask in the passing test

`wallet_test.go:15: got 0 want 10`

### ????

Well this is confusing, our code looks like it should work, we add the new amount onto our balance and then the balance method should return the current state of it.

In Go, when you call a function or a method the arguments are _copied_. 

So when we call `func (w Wallet) Deposit(amount int)` the `w` is a copy of whatever we called the method from. 

Without getting too computer-sciency, when you create a value - like a wallet, it is stored somewhere in memory. You can find out what the _address_ of that bit of memory with `&myVal`

Experiment by adding some prints to your code

```go
func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(10)

	got := wallet.Balance()
	
	fmt.Println("address of balance in test is", &wallet.balance)
	
	want := 10

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

```go
func (w Wallet) Deposit(amount int) {
	fmt.Println("address of balance in Deposit is", &w.balance)
	w.balance += amount
}
``` 

Now re-run the test

```
address of balance in Deposit is 0xc420012268
address of balance in test is 0xc420012260
```

You can see that the addresses of the two balances are different. So when we change the value of the balance inside the code, we are working on a copy of what came from the test. Therefore the balance in the test is unchanged.

We can fix this with _pointers_. Pointers let us _point_ to some values and then let us change them. So rather than taking a copy of the Wallet, we take a pointer to the wallet so we can change it.

```go
func (w *Wallet) Deposit(amount int) {
	w.balance += amount
}

func (w *Wallet) Balance() int {
	return w.balance
}
```

The difference is the receiver type is `*Wallet` rather than `Wallet` which you can read as "a pointer to a wallet".

Try and re-run the tests and they should pass. 

## Refactor
