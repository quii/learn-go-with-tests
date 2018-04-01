# Pointers and errors (WIP)

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

In the previous example we accessed fields directly with the field name, however in our _very secure wallet_ we don't want to expose our inner state to the rest of the world. We want to control access via methods.

## Try and run the test

`./wallet_test.go:7:12: undefined: Wallet`

## Write the minimal amount of code for the test to run and check the failing test output

The compiler doesn't know what a `Wallet` is so let's tell it.

```go
type Wallet struct { }
```

Now we've made our wallet, try and run the test again

```go
./wallet_test.go:9:8: wallet.Deposit undefined (type Wallet has no field or method Deposit)
./wallet_test.go:11:15: wallet.Balance undefined (type Wallet has no field or method Balance)
```

As expected, we need to define these methods to make the test pass. 

Remember to only do enough to make the tests run. We need to make sure our test fails correctly with a clear error message.

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

We will need some kind of _balance_ variable in our struct to store the state

```go
type Wallet struct {
	balance int
}
```

In Go if a symbol (so variables, types, functions et al) starts with a lowercase symbol then it is private _outside the package it's defined in_.

In our case we want our methods to be able to manipulate this value but no one else.

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

In Go, **when you call a function or a method the arguments are _copied_**. 

When calling `func (w Wallet) Deposit(amount int)` the `w` is a copy of whatever we called the method from. 

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

We said we were making a Bitcoin wallet but we have not mentioned them so far. We've been using `int` because they're a good type for counting things!

It seems a bit overkill to create a `struct` for this. `int` is fine in terms of the way it works but it's not descriptive.

Go lets you create **type aliases** which let you effectively create a new type out of an existing one.

The syntax is `type MyName OriginalType` 

```go
type Bitcoin int

type Wallet struct {
	balance Bitcoin
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}
```

```go
func TestWallet(t *testing.T) {

	wallet := Wallet{}

	wallet.Deposit(Bitcoin(10))

	got := wallet.Balance()

	want := Bitcoin(10)

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
```

To make `Bitcoin` you just use the syntax `Bitcoin(999)`

An interesting property of type aliasing is that you can also declare _methods_ on them. This can be very useful when you want to add some domain specific functionality on top of existing types.

[Let's implement Stringer on Bitcoin](https://golang.org/pkg/fmt/#Stringer)

```go
type Stringer interface {
        String() string
}
```

This interface is defined in the `fmt` package and let's you define how your type is printed when used with the `%s` format string in prints.

```go
func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}
```

As you can see, the syntax for creating a method on a type alias is the same as it is on a struct.

Next we need to update our test format strings so they will use `String()` instead.

```go
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
```

To see this in action, deliberately break the test so we can see it

`wallet_test.go:18: got 10 BTC want 20 BTC`

This makes it clearer what's going on in our test. 

The next requirement is for a `Withdraw` function.

## Write the test first

Pretty much the opposite of `Deposit()`

```go
func TestWallet(t *testing.T) {

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}

		wallet.Deposit(Bitcoin(10))

		got := wallet.Balance()

		want := Bitcoin(10)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}

		wallet.Withdraw(10)

		got := wallet.Balance()

		want := Bitcoin(10)

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

}
```

## Try and run the test

`./wallet_test.go:26:9: wallet.Withdraw undefined (type Wallet has no field or method Withdraw)`

## Write the minimal amount of code for the test to run and check the failing test output

```go
func (w *Wallet) Withdraw(amount Bitcoin) {
	
}
```

`wallet_test.go:33: got 20 BTC want 10 BTC`

## Write enough code to make it pass

```go
func (w *Wallet) Withdraw(amount Bitcoin) {
	w.balance -= amount
}
```

## Refactor

There's some duplication in our tests, let's refactor that out.

```go
func TestWallet(t *testing.T) {

	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		wallet.Withdraw(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

}
```

What should happen if you try and `Withdraw` more than is left in the account? For now, our requirement is to assume there is not an overdraft facility. 

How do we signal a problem when using `Withdraw` ? 

In Go, if you want to indicate an error it is idiomatic for your function to return an `err` for the caller to check and act on.

Let's try this out in a test

## Write the test first

```go
t.Run("Withdraw over balance limit", func(t *testing.T) {
    wallet := Wallet{balance: Bitcoin(20)}
    err := wallet.Withdraw(Bitcoin(100))
    
    if err == nil {
        t.Errorf("expected an error to be returned when withdrawing too much")
    }
})
```

We now want `Withdraw` to return an error _if_ you try and take out more than you have. 

We then check it has returned it by failing the test if it is `nil`

`nil` is synonomous with `null` from other programming languages. Errors can be `nil` because the return type of `Widthdraw` will be `error`, which is an interface. If you see a function that takes arguments or returns values that are interfaces, they can be nillable. 

Like `null` if you try and access a value that is `nil` it will throw a **runtime panic**. This is bad! You should make sure that you check for nils. 

## Try and run the test

`./wallet_test.go:31:25: wallet.Withdraw(Bitcoin(100)) used as value`

The wording is perhaps a little unclear, but our previous intent with `Withdraw` was just to call it, it will never return a value. To make this compile we will need to change it so it has a return type.

## Write the minimal amount of code for the test to run and check the failing test output

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
	w.balance -= amount
	return nil
}
```

Again, it is very important to just write enough code to satisfy the compiler. We correct our `Withdraw` method to return `error` and for now we have to return _something_ so let's just return `nil`

## Write enough code to make it pass

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return errors.New("oh no")
	}
	
	w.balance -= amount
	return nil
}
```

Remember to import `errors` into your code. 

`errors.New` creates a new `error` with a message of your choosing

## Refactor

We have a couple of tests around the same method so let's refactor it into a table test.

```go
t.Run("Withdraw", func(t *testing.T) {
    cases := []struct {
        name             string
        wallet           Wallet
        amountToWithdraw Bitcoin
        wantedBalance    Bitcoin
        wantedErr        bool
    }{
        {
            name:             "sufficient funds",
            wallet:           Wallet{Bitcoin(20)},
            amountToWithdraw: Bitcoin(10),
            wantedBalance:    Bitcoin(10),
            wantedErr:        false,
        },
        {
            name:             "insufficient funds",
            wallet:           Wallet{Bitcoin(20)},
            amountToWithdraw: Bitcoin(100),
            wantedBalance:    Bitcoin(20),
            wantedErr:        true,
        },
    }

    for _, tt := range cases {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.wallet.Withdraw(tt.amountToWithdraw)

            if tt.wallet.Balance() != tt.wantedBalance {
                t.Errorf("got balance %s want %s", tt.wallet.Balance(), tt.wantedBalance)
            }

            if tt.wantedErr && err == nil {
                t.Error("wanted an error but didn't get one")
            }

            if !tt.wantedErr && err != nil {
                t.Errorf("didnt want an error but got one %s", err)
            }
        })
    }
})
```

Our test cases are describing our intent a little clearer now.

Hopefully when returning an error of "oh no" you were thinking that we _might_ iterate on that because it doesn't seem that useful to return.

Assuming that the error ultimately gets returned to the user, let's update our test to assert on some kind of error message rather than just the existence of an error

## Write the test first

```go
t.Run("Withdraw", func(t *testing.T) {
    cases := []struct {
        name             string
        wallet           Wallet
        amountToWithdraw Bitcoin
        wantedBalance    Bitcoin
        wantedErr        error
    }{
        {
            name:             "sufficient funds",
            wallet:           Wallet{Bitcoin(20)},
            amountToWithdraw: Bitcoin(10),
            wantedBalance:    Bitcoin(10),
            wantedErr:        nil,
        },
        {
            name:             "insufficient funds",
            wallet:           Wallet{Bitcoin(20)},
            amountToWithdraw: Bitcoin(100),
            wantedBalance:    Bitcoin(20),
            wantedErr:        errors.New("cannot withdraw, insufficient funds"),
        },
    }

    for _, tt := range cases {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.wallet.Withdraw(tt.amountToWithdraw)

            assertBalance(t, tt.wallet, tt.wantedBalance)

            if tt.wantedErr != nil {
                if err == nil {
                    t.Error("wanted an error but didn't get one")
                }

                if err.Error() != tt.wantedErr.Error() {
                    t.Errorf("got err '%s' want '%s'", err.Error(), tt.wantedErr)
                }
            }

            if tt.wantedErr == nil && err != nil {
                t.Errorf("didn't want an error but got one %s", err)
            }
        })
    }
})
```

- We have changed the table so that `err` is now an `error`. This lets us define a particular kind of error to look for in the test and also let us put the value `nil` if we have a case where we don't want an error.
- Introduced `t.Fatal` which will stop the test if it is called. This is because we don't want to make any more assertions on the error returned if there isn't one around. Without this the test would carry on to the next step and panic because of a nil pointer.

## Try and run the test

`wallet_test.go:61: got err 'oh no' want 'cannot withdraw, insufficient funds'`

## Write enough code to make it pass

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return errors.New("cannot withdraw, insufficient funds")
	}

	w.balance -= amount
	return nil
}
```

## Refactor

We have duplication of the error message in both the test code and the `Withdraw` code. 

It would be really annoying for the test to fail if someone wanted to re-word the error and it's just too much detail for our test. We don't _really_ care what the exact wording is, just that some kind of meaningful error around withdrawing is returned given a certain condition. 

In Go, errors are values, so we can refactor it out into a variable and have a single source of truth for it. 

```go
var InsufficientFundsError = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return InsufficientFundsError
	}

	w.balance -= amount
	return nil
}
```

This is a positive change in itself because now our `Withdraw` function looks very clear.

```go
t.Run("Withdraw", func(t *testing.T) {
    cases := []struct {
        name             string
        wallet           Wallet
        amountToWithdraw Bitcoin
        wantedBalance    Bitcoin
        wantedErr        error
    }{
        {
            name:             "sufficient funds",
            wallet:           Wallet{Bitcoin(20)},
            amountToWithdraw: Bitcoin(10),
            wantedBalance:    Bitcoin(10),
            wantedErr:        nil,
        },
        {
            name:             "insufficient funds",
            wallet:           Wallet{Bitcoin(20)},
            amountToWithdraw: Bitcoin(100),
            wantedBalance:    Bitcoin(20),
            wantedErr:        InsufficientFundsError,
        },
    }

    for _, tt := range cases {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.wallet.Withdraw(tt.amountToWithdraw)

            assertBalance(t, tt.wallet, tt.wantedBalance)

            gotAnError := err != nil
            wantAnError := tt.wantedErr != nil

            if gotAnError != wantAnError {
                t.Fatalf("got error '%s' want '%s'", err, tt.wantedErr)
            }

            if wantAnError && err != tt.wantedErr {
                t.Errorf("got err '%s' want '%s'", err.Error(), tt.wantedErr)
            }
        })
    }
})
```

And now the test is easier to follow too.

Another useful property of tests is that they help us understand the _real_ usage of our code so we can make sympathetic code. We can see here that a developer can simply call our code and do an equals check to `InsufficientFundsError` and act accordingly.

## Wrapping up

### Pointers

- Go copies values when you pass them to functions/methods so if you're writing a function that needs to mutate state you'll need it to take a pointer to the thing you want to change.
- The fact that Go takes a copy of values is useful a lot of the time but sometimes you wont want your system to make a copy of something, in which case you need to pass a reference. Examples could be very large data or perhaps things you intend only to have one instance of (like database connection pools)

### nil

- Pointers can be nil
- When a function returns a pointer to something, you need to make sure you check if it's nil or you might raise a runtime exception, the compiler wont help you here.
- Useful for when you want to describe a value that could be missing

### Errors

- Errors are the way to signify failure when calling a function/method
- By listening to our tests we concluded that checking for a string in an error would result in a flaky test. So we refactored to use a meaningful value instead and this resulted in easier to test code and concluded this would be easier for users of our API too. 
- This is not the end of the story with error handling, you can do more sophisticated things but this is just an intro. Later sections will cover more strategies.
- [Don’t just check errors, handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully) 

### Type aliases

- Useful for adding more domain specific meaning to values
- Can let you implement interfaces

Pointers and errors are a big part of writing Go that you need to get comfortable with. Thankfully the compiler will _usually_ help you out if you do something wrong, just take your time and read the error.