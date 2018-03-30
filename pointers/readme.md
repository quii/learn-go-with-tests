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

In the previous example we accessed fields directly with the field name, however in our _very secure wallet_ we dont want to expose our inner state to the rest of the world. We want to control access via methods.

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

There is a temptation to add a table driven test for `Withdraw` but let's resist that temptation for now. 

Hopefully you may be thinking that the error of "oh no" could maybe be a little improved. 

## Write the test first

```go
t.Run("Withdraw over balance limit", func(t *testing.T) {
    wallet := Wallet{balance: Bitcoin(20)}
    err := wallet.Withdraw(Bitcoin(100))

    if err == nil {
        t.Errorf("expected an error to be returned when withdrawing too much")
    }

    expectedErrorMessage := "cannot withdraw 100 BTC, insufficient funds (20 BTC)"
    if err.Error() != expectedErrorMessage {
        t.Errorf(`got error message of "%s", want "%s"`, err.Error(), expectedErrorMessage)
    }
})
```
## Try and run the test

`wallet_test.go:39: got error message of "oh no", want "cannot withdraw 100 BTC, insufficient funds (20 BTC)"`

## Write enough code to make it pass

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return fmt.Errorf("cannot withdraw %s, insufficient funds (%s)", amount, w.balance)
	}

	w.balance -= amount
	return nil
}
```

Remember to remove the import of `errors` and add `fmt`. 

`fmt.Errorf` is like `t.Errorf` in that it takes a format string and some values but instead of failing a test it returns an error

## Refactor

The main problem we have is the potential of a flaky test. We are asserting on the exact wording of the error, which is not very important to us. If a developer decides to change the wording of the message, a test will fail which seems too heavy handed. **Remember our tests should not be a burden**


In addition, the useful data in the error is "trapped" inside a `string`. If a developer wants to do something useful with this data she is going to be quite stuck. 

What we really want to assert is that we have useful information in our error but still somehow return an `error`

As mentioned before, [`error` is an interface](https://golang.org/pkg/builtin/#error).

```go
type error interface {
        Error() string
}
```

From the previous sections we learned how to implement interfaces. So what we can do is create a custom error type, which has raw values accessible to the caller of `Withdraw`.

This gives the users of our library some flexibility in their error handling:

- They can extract out the pertinent values of the error and do _something different_
- Simply use the `Error()` as is, perhaps logging it or printing it to the user

*Plus* it makes our tests more useful and less prone to error due to wording changes.

## Write the test first

```go
t.Run("Withdraw over balance limit", func(t *testing.T) {
    wallet := Wallet{balance: Bitcoin(20)}
    err := wallet.Withdraw(Bitcoin(100))

    if err == nil {
        t.Fatalf("expected an error to be returned when withdrawing too much")
    }

    got, isWithdrawErr := err.(WithdrawError)

    if !isWithdrawErr {
        t.Fatalf("did not get a withdraw error %#v", err)
    }

    want := WithdrawError{
        AmountToWithdraw: Bitcoin(100),
        CurrentBalance:   Bitcoin(20),
    }

    if want != got {
        t.Errorf("got %#v, want %#v", got, want)
    }

})
```

_We will probably refactor this!_ Hold your nose while we go through some concepts. 

As we discussed earlier, by design the concept of interfaces hides the concrete type. Most of the time this is super nice but in this case we do need to check the _particular type_ is returned because we want to make sure the information in the concrete type is returned. 

We can do this with a _type assertion_. When you have a value and all you know is it's an interface, you can ask Go if a particular thing is _actually_ type `Foo`. It returns two values, the value _cast to the type_ and a boolean telling you the result of the check.

The syntax is `theThingCastToThetype, booleanConfirmingItIsTheType := thing.(MyType)`. You can see this in action in the test code. 

- First important change is we change the nil check to use `t.Fatalf` rather than `t.Errorf`. `Fatalf` is helpful if you want the test to stop. `Errorf` will fail the test but the rest of the code will continue. In our case if we don't get an error there's no point in carrying on.
- We then do our type assertion. We check that it is the type we want and if not we fail the test.
- If it _is_ a `WithdrawError` then we check it's values with a normal assertion.

## Try and run the test

`./wallet_test.go:37:30: undefined: WithdrawError`

## Write the minimal amount of code for the test to run and check the failing test output

We have not defined our new error type yet

```go
type WithdrawError struct {
	CurrentBalance   Bitcoin
	AmountToWithdraw Bitcoin
}
```

Try again.

```
./wallet_test.go:37:28: impossible type assertion:
	WithdrawError does not implement error (missing Error method)
```

Go knows that our current type cannot possibly be an `error` due to the missing method, so lets implement the `error` interface on our new type.

```go
func (w WithdrawError) Error() string {
	return fmt.Sprintf("cannot withdraw %s, insufficient funds (%s)", w.AmountToWithdraw, w.CurrentBalance)
}
```

Finally the test runs and fails as we'd expect

`wallet_test.go:40: did not get a withdraw error &errors.errorString{s:"cannot withdraw 100 BTC, insufficient funds (20 BTC)"}`

Our `Withdraw` method is failing the type assertion

## Write enough code to make it pass

Make the method use our new type instead

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

	if amount > w.balance {
		return WithdrawError{
			CurrentBalance:   w.balance,
			AmountToWithdraw: amount,
		}
	}

	w.balance -= amount
	return nil
}
```

The test should now pass. 

## Refactor

Our test is not reading great, there's a lot of mechanics around type assertions that makes the intent of the test a little unclear. We can refactor out our type assertion code into a helper like usual.

```go
func TestWallet(t *testing.T) {

	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}
	
	assertWithdrawError := func(t *testing.T, err error, want WithdrawError) {
		got, isWithdrawErr := err.(WithdrawError)

		if !isWithdrawErr {
			t.Fatalf("did not get a withdraw error %#v", err)
		}

		if want != got {
			t.Errorf("got %#v, want %#v", got, want)
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

	t.Run("Withdraw over balance limit", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(100))

		if err == nil {
			t.Fatalf("expected an error to be returned when withdrawing too much")
		}

		assertWithdrawError(t, err, WithdrawError{
			AmountToWithdraw: Bitcoin(100),
			CurrentBalance:   Bitcoin(20),
		})


	})

}
```

This is a bit of an improvement, but we can do better. 

There's a few things around `Withdraw` we have not tested 

- If you fail to `Withdraw` it should have no effect on the `balance`
- If you successfully `Withdraw` an error should not be returned. 

This _feels_ like it might be well suited to a table based test. We want to having some kind of wallet, do `Withdraw` and see what happens with differing amounts of Bitcoin.

```go
func TestWallet(t *testing.T) {

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {

		cases := []struct {
			description      string
			wallet           Wallet
			amountToWithdraw Bitcoin
			wantedBalance    Bitcoin
			wantedErr        *WithdrawError
		}{
			{
				description:      "happy withdraw",
				wallet:           Wallet{balance: Bitcoin(10)},
				amountToWithdraw: Bitcoin(5),
				wantedBalance:    Bitcoin(5),
				wantedErr:        nil,
			},
			{
				description:      "not enough funds",
				wallet:           Wallet{balance: Bitcoin(10)},
				amountToWithdraw: Bitcoin(20),
				wantedBalance:    Bitcoin(10),
				wantedErr:        &WithdrawError{AmountToWithdraw: Bitcoin(20), CurrentBalance: Bitcoin(10)},
			},
		}

		for _, tt := range cases {
			t.Run(tt.description, func(t *testing.T) {
				err := tt.wallet.Withdraw(tt.amountToWithdraw)

				assertBalance(t, tt.wallet, tt.wantedBalance)

				if tt.wantedErr != nil {
					assertWithdrawError(t, err, *tt.wantedErr)
				}
			})
		}
	})

}

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertWithdrawError(t *testing.T, err error, want WithdrawError) {
	got, isWithdrawErr := err.(WithdrawError)

	if !isWithdrawErr {
		t.Fatalf("did not get a withdraw error %#v", err)
	}

	if want != got {
		t.Errorf("got %#v, want %#v", got, want)
	}
}
```

There's a few new things here

- We talked about `nil` earlier and how pointers to things can be nil. In our cases we want to be able to control whether we want an error or not so the type of `wantedErr` to `*WithdrawError`. So we can check for an error if we want to.
- In order to get a _pointer to a value_ you use the `&` operator on the value. You can see that in our second case for the error we're looking for
- Our helper function doesn't take a pointer, it takes the full value so we _dereference the pointer_ with `*` when calling `assertWithdrawError`
- I decided that our helper functions were cluttering up our test code. They're not very interesting so I moved them out and into the bottom of the file.

# Wrapping up

## Pointers

- Go copies values when you pass them to functions/methods so if you're writing a function that needs to mutate state you'll need it to take a pointer to the thing you want to change.
- The fact that Go takes a copy of values is useful a lot of the time but sometimes you wont want your system to make a copy of something, in which case you need to pass a reference. Examples could be very large data or perhaps things you intend only to have one instance of (like database connection pools)

## nil

- Pointers can be nil
- When a function returns a pointer to something, you need to make sure you check if it's nil or not or you will get a runtime exception, the compiler wont help you here.
- Useful for when you want to describe a value that could be missing

## Errors

- Errors are the way to signify failure when calling a function/method
- You _may_ wish to introduce your own type of error to let developers work with whatever problems come up
    - In our case we _listened to our tests_ and came to the conclusion that we would rather assert on some useful data, rather than a string check
    - Listening to your tests is important. Often if your tests are hard to write/read then the users of your code are also going to have a tough time. This is why TDD is praised as being helpful as a design tool.
- [Don’t just check errors, handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully) 

## Type aliases

- Useful for adding more domain specific meaning to values
- Can let you implement interfaces

Pointers and errors are a big part of writing Go that you need to get comfortable with. Thankfully the compiler will _usually_ help you out if you do something wrong, just take your time and read the error.