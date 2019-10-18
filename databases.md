# Databases

<!-- TODO: add links to start:v1, end:v2 -->

Oftentimes when creating software, it's necessary to save (or, more precisely, _persist_) some application state.

As an example, when you log into your online banking system, the system has to:

1. Check that it's really you accessing the system (this is called _authentication_, and is beyond the scope of this chapter)
2. Retrieve some information from _somewhere_ and show it to the user (you).

Information that is stored and meant to be long-lived is said to be [_persisted_](<https://en.wikipedia.org/wiki/Persistence_(computer_science)>), usually on a medium that can reliably reproduce the data stored.

Some storage systems, like the filesystem, can be effective for one-off or small amounts of storage, but they fall short for larger application, for a number of reasons.

This is why most software applications, large and small, opt for storage systems that can provide:

-   Reliability: The data you want is there when you need it.
-   Concurrency: Imagine thousands of users accessing simultaneously.
-   Consistent: You expect the same inputs to produce the same results.
-   Durable: Data should remain there even in case of a system failure (power outage or system crash).

NOTE: The above bullet points are a rewording of the [_ACID principles_](https://en.wikipedia.org/wiki/ACID), it's a set of properties often expressed and used in database design.

_Databases_ are storage mediums that can provide these properties, and much much more.

Also note that, in general, there are two large branches of database types, [SQL](https://en.wikipedia.org/wiki/SQL) and [NOSQL](https://en.wikipedia.org/wiki/NoSQL). In this chapter we will be focusing on SQL databases, using the [database/sql](https://golang.org/pkg/database/sql) package and the `postgres` driver [pq](_ "github.com/lib/pq").

There is a fair bit of CLI usage in this chapter (mainly setting up the database). For the sake of simplicity we will assume that you are running a system `ubuntu` on your machine, with `bash` installed. In the near future, look into the appendix for installation on other systems.

## A note on RDBMS choice

RDBMS (**R**elational **D**ata**B**ase **M**anagement **S**ystem) is a software program that allows users to operate on the storage engine underneath (namely, the database itself).

There are many choices and many capable systems, each one with its strenghts and weaknesses. I encourage you to do some research on the subject in case you're not familiar with the different options.

In this chapter we will be using [PostgreSQL](https://www.postgresql.org/): a mature, production ready relational database that has been proven to be extremely reliable.

The reasons for this choice, include, but are not limited to:

-   Postgres doesn't hold your hand.

    While there are GUI tools for visual exploration, Postgres comes by default with only a CLI. This makes for a better understanding of the SQL commands itself, and also makes scripting much easier (we won't be covering database scripting in this guide.)

-   It's production ready

    The default settings for `PostgreSQL` are good enough to be used in a production environment (with some caveats). Using it during development helps us close the gap between the testing, staging and production environments (also referred to as [dev/prod parity](https://www.12factor.net/dev-prod-parity)). As you will soon see, this will present a challenge during development, that, when overcome, renders your entire application more reliable (hint: integration tests).

## Getting a PostgreSQL instance running

### Docker

The easiest (and cleanest) way of getting `PostgreSQL` up and running is by using `docker`. This will create the database and user:

-   [`Docker` installation instructions](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
-   See [https://hub.docker.com/\_/postgres](https://hub.docker.com/_/postgres) for more details on how to use this image.

```sh
~$ docker run \
    -e POSTGRES_DB=bookshelf_db \ # name for the database
    -e POSTGRES_USER=bookshelf_user \ # name for the database user
    -e POSTGRES_PASSWORD=secret-password \ # database password
    -p 5432:5432 \ # map port 5432 on host to docker container's 5432
	-d \ # detach process
    postgres:11.5 # get official postgres image, version 11.5
```

You may need to run the above command with elevation (prepend it with `sudo`).

### Manual installation

Install `PostgreSQL` with the package manager

```sh
~$ sudo apt-get upgrade
~$ sudo apt-get install postgresql postgresql-contrib
```

PostgreSQL installs and initializes a database called `postgres`, and a user also called `postgres`. Since this is a system-wide install, we don't want to pollute this main database with this application's tables (`PostgreSQL` uses these to store administrative data), so we will have to create a user and a database.

```sh
~$ sudo -i -u postgres # this will switch you to the postgres user
~$ psql
```

```
psql (10.10 (Ubuntu 10.10-0ubuntu0.18.04.1))
Type "help" for help

postgres=# CREATE USER bookshelf_user WITH CREATEDB PASSWORD 'secret-password';
CREATE ROLE
postgres=# CREATE DATABASE bookshelf_db OWNER bookshelf_user;
CREATE DATABASE
```

You can view users and databases with the commands \du and \l respectively.

## Database migrations

[Database migrations](https://en.wikipedia.org/wiki/Schema_migration) refers to the management of incremental, reversible changes and version control to the database. This is an important part of any system that implements a database, as it allows developers to reproduce the database in a deterministic manner.

Migrations are a complicated subject, well beyond the scope of this book. If I were to oversimplify, I would classify migrations into two large subcategories:

1. Those that only modify the `schema` (the tables, indexes and configuration of the database).
2. Those that modify data and/or the `schema` (a migration that adds/removes columns or rows from a table).

We will only be dealing with type no. 1 in this chapter.

### Justification

If you are wondering "What's all this about migrations?", it's because we will need them later in this chapter when we start writing integration tests for our sample application.

### Migration tools

There are excellent, open source migration management tools out there ([exhibit 1](https://github.com/golang-migrate/migrate), [exhibit 2](https://github.com/mattes/migrate), [exhibit 3](https://github.com/go-pg/migrations)).

But we will not be using those. We will write our own tool (although simpler than those linked above) instead, using the `go` toolchain. We will also try to adhere to best practices, these being:

-   Migrations should be _ordered_, so there is a definitive order each time they are run. We will simply prepend each filename with a number, allowing enough digits for a large number of files (although realistically, migrations get "squashed" into a smaller number of files on real projects, if it makes sense to do so).
-   Migrations should be _reversible_. Applying a change to a database should simple to perform, and reversing said change should be easy as well. We will write two migration files for every change, appropriately suffixed `up` and `down`.
-   Migrations should be _idempotent_. Re-applying a change that has previously been applied should yield no effects beyond those of the initial application.

With these in mind, let's dive in...

## Project

We will initially write a (simple) tool to handle our (also simple) migrations, then we will be creating a `CRUD` program to interact with our spiffy, real database.

## Write the test first

```go
// migrate_test.go
package main

import (
	"testing"
)

func TestMigrateUp(t *testing.T) {
	store, removeStore := NewStore()
	defer removeStore()

	const numberOfMigrations = 1
	err := MigrateUp(store, "migrations-directory", numberOfMigrations)
	if err != nil {
		t.Errorf("received error but didn't want one: %v", err)
	}
}
```

## Try to run the test

Fails, as expected.

```sh
# command-line-arguments [command-line-arguments.test]
.\migrate_test.go:9:10: undefined: MigrateUp
FAIL    command-line-arguments [build failed]
FAIL
```

## Write enough code to make it pass

We will take some liberties with writing code, keeping it in small batches and testing as we go.

Below is some boilerplate that will assist in creating the database connection. We need an `*sql.DB` instance that we can pass on to our tests first. As we learned in the `Dependency Injection` chapter, we should make a helper method so we can acquire an instance from anywhere in our application, while passing in the dependencies, in our case, the database connection string.

Before we continue down the happy-path, we need to make sure our `MigrateUp` function works as expected (unit-test) before we attempt to call it on the real database (integration-test).

The code directly below defines an interface, which will allow us to mock the database functionality, as well as a `NewStore` function that will simplify its creation.

```go
// bookshelf-store.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Storer interface {
	ApplyMigration(name, stmt string) error
}

type Store struct {
	db *sql.DB
}

const (
	removeTimeout = 10 * time.Second
)

func NewStore() (*Store, func()) {
	// remember to change 'secret-password' for the password you set earlier
	const connStr = "postgres://bookshelf_user:secret-password@localhost:5432/bookshelf_db"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	// exponential backoff
	remove := func() {
		deadline := time.Now().Add(removeTimeout)
		for tries := 0; time.Now().Before(deadline); tries++ {
			err := db.Close()
			retryIn := time.Second << uint(tries)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error closing connection to database, retrying in %v: %v\n", retryIn, err)
				time.Sleep(retryIn)
				continue
			}
			return
		}
		log.Fatalf("timeout of %v exceeded", removeTimeout)
	}

	return &Store{db: db}, remove
}

func (s *Store) ApplyMigration(name, stmt string) error {
	_, err := s.db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func main(){}
```

---

### Note on exponential backoff (`remove` anonymous function)

Generally when calling external services you want to account for the possibility of failure, yet the reasons why the service could fail are numerous. The database is an external service because (generally) databases run in a client-server paradigm. In our case, the database connection could fail to close for a number of reasons. This pattern is useful to give the service some time to correct itself or reset before attempting our call again. This snippet does so by waiting on powers of 2 incrementally: 1, 2, 4, 8, 16, ..., N seconds, with a hard limit set at `removeTimeout` seconds.

This is here as a demonstration mostly, where it will prove useful is during the integration tests. The test database that we create will have be destroyed after running tests, thus, it can't have read or write operations running.

Code explained:

```go
remove := func() {
	// sets the deadline to Now + removeTimeout seconds
	deadline := time.Now().Add(removeTimeout)
	// it's a normal for loop, but the failure condition has been
	// swapped by `time.Now().Before(deadline)`, which returns a boolean
	// see https://golang.org/pkg/time/#Time.Before
	for tries := 0; time.Now().Before(deadline); tries++ {
		err := db.Close()
		// as `tries` increases every loop, `retryIn` becomes a
		// time.Second` unit to the `tries` of 2
		// https://play.golang.org/p/ubyLNhxE31K has an illustrative example
		retryIn := time.Second << uint(tries)
		if err != nil {
			...
			time.Sleep(retryIn)
			continue
		}
		return
	}
	// panic and log the failure
	log.Fatalf("timeout of %v exceeded", removeTimeout)
}
```

---

We could add more methods to the `Storer` interface, but at this point, we don't need them. So we won't for now.

The `ApplyMigration` method is merely a wrapper around the `sql.DB.Exec` method, but this allows us to abstract it on our unit tests, testing that the `MigrateUp`, function does what it's intended

Here is the signature for our `MigrateUp` function, with our `Storer` interface:

```go
MigrateUp(store Storer, dir string, num int)
```

There's a lot to our `MigrateUp` function, so let's break it down and test each step.

1. It needs to check whether the `dir` passed in exists.
2. It needs to get all the filenames inside `dir`, and allow us to iterate over them.

    2.1 It should allow for ordered iteration through the files.

3. It needs to run only `up` migrations.
4. It needs to run migrations `1` through `num`, or all of them if `num == -1`.
5. Lastly, it needs to report on the success of each migration run, if a migration fails, the entire process should be halted.

Seeing as we will have a `MigrateDown` as well, and so far they seem only to differ on step `3`, we can use a little foresight and create a utility function `migrate`, which will be used by both variants.

## Write the test first

Let's write our mock store (and test) first, so we can test the `migrate` function in isolation.

```go
// migrate_test.go
import (
	"time"
	...
)

type migration struct {
	created time.Time
	name string
	stmt string
	called int
}

type SpyStore struct {
	migrations map[string]migration
}

func (s *SpyStore) ApplyMigration(name, stmt string) error {
	var m migration
	if mig, ok := s.migrations[name]; ok {
		m = mig
		m.called++
		return nil
	}
	m := migration{
		name: name,
		stmt: stmt,
		created: time.Now(),
	}
	m.called++
	s.migrations[name] = m
	return nil
}
func NewSpyStore() {
	return &SpyStore{map[string]migration{}}
}

```

And our migration files, which will be used in the tests

```sh
~$ mkdir migrations
```

Then create the `.sql` files inside this dir, name them `0001_create_bookshelf_table.up.sql` and `0001_create_bookshelf_table.down.sql`, with the SQL below.

```sql
-- migrations/0001_create_bookshelf_table.up.sql
BEGIN;
CREATE TABLE IF NOT EXISTS books (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	author VARCHAR(255) NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS  books_id_uindex ON books (id);
COMMIT;

-- migrations/0001_create_bookshelf_table.down.sql
BEGIN;
DROP INDEX IF EXISTS books_id_uindex CASCADE;
DROP TABLE IF EXISTS books CASCADE;
COMMIT;
```

Now we're ready to address the different points, one by one.

1. It needs to check whether the `dir` passed in exists.

```go
...
func TestMigrate(t *testing.T) {
	store := NewSpyStore()
	t.Run("error on nonexistent directory", func(t *testing.T){
		err := migrate(store, "i-do-not-exist", -1)
		if err != nil {
			t.Errorf("got an error but didn't want one: %v", err)
		}
	})

	t.Run("no error on existing directory", func(t *testing.T){
		err := migrate(store, "migrations", -1)
		if err == nil {
			t.Error("wanted an error but didn't get one")
		}
	})
}
```

## Try to run the test

As expected, it fails.

```sh
# github.com/quii/learn-go-with-tests/databases/v1 [github.com/quii/learn-go-with-tests/databases/v1.test]
.\migrate_test.go:37:10: undefined: migrate
.\migrate_test.go:44:10: undefined: migrate
FAIL    github.com/quii/learn-go-with-tests/databases/v1 [build failed]
```

## Write enough code to make it pass

```go
// bookshelf-store.go
func migrate(store Storer, dir string, num int) error {
	return nil
}
```

```sh
--- FAIL: TestMigrate (0.00s)
    --- FAIL: TestMigrate/no_error_on_existing_directory (0.00s)
        migrate_test.go:45: wanted an error but didn't get one
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v1        0.608s
```

We need an error for non-existent directories. Thankfully, this is easy with `os.Stat` and `os.IsNotExist`

```go
func migrate(store Storer, dir string, num int) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory %q does not exist", dir)
	}
	return nil
}
```

```sh
PASS
ok      github.com/quii/learn-go-with-tests/databases/v1        0.484s
```

On to the next point

2. It needs to get all the filenames inside `dir`, and allow us to iterate over them.

We can use [`ioutil.ReadDir`](https://golang.org/pkg/io/ioutil/#ReadDir) to implement the desired functionality. Since we would also like to prevent an empty directory from breaking our code, we should check for that as well; we can use [`ioutil.TempDir`](https://golang.org/pkg/io/ioutil/#TempDir) to create temporary directories for our tests.

## Write the test first

```go
// migrate-test.go
import (
	...
	"os"
	"io/ioutil"
)
func TestMigrate(t \*testing.T) {
	...
		t.Run("error on empty directory", func(t *testing.T) {
		// create temporary directory
		tmpdir, err := ioutil.TempDir("", "test-migrations")
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		defer os.RemoveAll(tmpdir) // cleanup when done

		err = migrate(store, tmpdir, -1)
		if err == nil {
			t.Error("wanted an error but didn't get one")
		}
	})

	t.Run("non-empty directory attempts to migrate", func(t *testing.T) {
		// create temporary directory
		tmpdir, err := ioutil.TempDir("", "test-migrations")
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		defer os.RemoveAll(tmpdir) // cleanup when done

		// create temporary files
		for _, filename := range []string{
			"01.*.up.sql",
			"01.*.down.sql",
			"02.*.up.sql",
			"02.*.down.sql",
		}{
			tmpfile, err := ioutil.TempFile(tmpdir, filename)
			if err != nil {
				fmt.Println(err)
				t.FailNow()
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.Write([]byte(filename + " SQL content")); err != nil {
				tmpfile.Close()
				fmt.Println(err)
				t.FailNow()
			}
			if err := tmpfile.Close(); err != nil {
				fmt.Println(err)
				t.FailNow()
			}
		}

		err = migrate(store, tmpdir, -1)
		if err != nil {
			t.Errorf("got an error but didn't want one: %v", err)
		}
	})
}
```

## Try to run the tests

```sh
--- FAIL: TestMigrate (0.00s)
    --- FAIL: TestMigrate/error_on_empty_directory (0.00s)
        migrate_test.go:72: wanted an error but didn't get one
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v2        0.389s
```

Only one failure, and no migration input. This is because our store is completely empty (as no code is implemented yet). Let's add a check for it before we move on.

```go
// migrate-test.go
...
func TestMigrate(t \*testing.T) {
	...
	t.Run("non-empty directory attempts to migrate", func(t *testing.T) {
		...
		if len(store.migrations) == 0 {
			t.Error("no migrations in store")
		}
		for _, m  := range store.migrations {
			want := 1
			if m.called != want {
				t.Errorf("wanted %d call got %d calls for %s migration", want,  m.called, m.name)
			}
		}
	})
}
```

```sh
--- FAIL: TestMigrate (0.03s)
    --- FAIL: TestMigrate/error_on_empty_directory (0.00s)
        migrate_test.go:72: wanted an error but didn't get one
    --- FAIL: TestMigrate/non-empty_directory_attempts_to_migrate (0.03s)
        migrate_test.go:115: no migrations in store
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v2        0.521s
```

That's better.

## Write enough code to make it pass

We finally get to use the interface! Here, we call `ApplyMigration` inside our `migrate` function

```go
// bookshelf-store.go
import (
	...
	"errors"
	"io/ioutil"
	"path/filepath"
)
func migrate(store Storer, dir string, num int) error {
	...
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errors.New("empty migration file")
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read migration file %s, %v", file.Name(), err)
			return err
		}
		err = store.ApplyMigration(file.Name(), string(content))
		if err != nil {
			return err
		}
	}
	return nil
}
```

```sh
PASS
ok      github.com/quii/learn-go-with-tests/databases/v2        0.512s
```

Just to see what the failing output would look like, change the `want` variable value to 2 (so it fails) in the `non-empty directory attempts to migrate` and run the tests.

```sh
--- FAIL: TestMigrate (0.05s)
    --- FAIL: TestMigrate/non-empty_directory_attempts_to_migrate (0.05s)
        migrate_test.go:120: wanted 2 call got 1 calls for 02.017392150.down.sql migration
        migrate_test.go:120: wanted 2 call got 1 calls for 02.540392659.up.sql migration
        migrate_test.go:120: wanted 2 call got 1 calls for 0001_create_books_table.down.sql migration
        migrate_test.go:120: wanted 2 call got 1 calls for 0001_create_books_table.up.sql migration
        migrate_test.go:120: wanted 2 call got 1 calls for 01.434163769.up.sql migration
        migrate_test.go:120: wanted 2 call got 1 calls for 01.969306692.down.sql migration
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v2        0.556s
```

Uh oh. We have several problems that this output reveals:

-   Our `migrate` function is running `up` and `down` migrations indiscriminately. We need to add a testcase so that it only runs one kind at a time.
-   It's including cases from previous tests. This is an easy fix, just need to instantiate a new `SpyStore` on each test.

## Refactor

Let's take this opportunity to clean up our tests, by making some assertions and other helper functions to make the actual tests more succint.

We'll also make some error variables to better define our expected errors.

```go
// bookshelf-store.go
...
var (
	ErrMigrationDirEmpty = errors.New("empty migration directory")
	ErrMigrationDirNoExist = errors.New("migration directory does not exist")
)
...
func migrate(store Storer, dir string, num int) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return ErrMigrationDirNoExist
	}
	...
	if len(files) == 0 {
		return ErrMigrationDirEmpty
	}
	...
}
```

```go
// migrate_test.go
import (
	...
	"path/filepath"
)
...
func AssertError(t *testing.T, got, want error) {
	t.Helper()
	if got == nil {
		t.Error("wanted an error but didn't get one")
	}
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertNoError(t *testing.T, got error) {
	t.Helper()
	if err != nil {
		t.Errorf("got an error but didn't want one: %v", err)
	}
}

func AssertStoreMigrationCalls(t *testing.T, store *SpyStore, name string, num int) {
	t.Helper()

	m, ok := store.migrations[m]
	if !ok {
		t.Errorf("migration %q does not exist in store", name)
	}
	if m.called != num {
		t.Errorf("got %d want %d calls migration %q", num, m.called, name)
	}
}

func AssertAllStoreMigrationCalls(t *testing.T, store *SpyStore, num int, direction string) {
	t.Helper()

	for _, m  := range store.migrations {
		if !strings.HasSuffix(m.name, direction+".sql") {
			continue
		}
		AssertStoreMigrationCalls(t, store, m.name, num)
	}
}

func CreateTempDir(
	t *testing.T,
	name string,
	empty bool,
) (string, []string, func()) {
	t.Helper()

	tmpdir, err := ioutil.TempDir("", name)
	if err != nil {
		fmt.Println(err)
		os.RemoveAll(tmpdir)
		t.FailNow()
	}
	filenames := make([]string, 0)
	if !empty {
		for _, filename := range []string{
			"01.*.up.sql",
			"01.*.down.sql",
			"02.*.up.sql",
			"02.*.down.sql",
		}{
			tmpfile, err := ioutil.TempFile(tmpdir, filename)
			if err != nil {
				fmt.Println(err)
				os.Remove(tmpfile.Name())
				t.FailNow()
			}
			filenames = append(filenames, filepath.Base(tmpfile.Name()))

			if _, err := tmpfile.Write([]byte(filename + " SQL content")); err != nil {
				tmpfile.Close()
				fmt.Println(err)
				os.Remove(tmpfile.Name())
				t.FailNow()
			}
			if err := tmpfile.Close(); err != nil {
				fmt.Println(err)
				os.Remove(tmpfile.Name())
				t.FailNow()
			}
		}

	}
	cleanup := func() {
		os.RemoveAll(tmpdir)
	}
	return tmpdir, filenames, cleanup
}
```

Our second test `no error on existing directory` can be removed, as an empty, existing directory raises an error as well. We've added the `direction` test as well.

```go
...
import (
	...
	"strings"
)
func TestMigrate(t *testing.T) {
	t.Run("error on nonexistent directory", func(t *testing.T) {
		store := NewSpyStore()
		err := migrate(store, "i-do-not-exist", -1)

		AssertError(t, err, ErrMigrationDirNoExist)
	})
	t.Run("error on empty directory", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", true)
		defer cleanup()

		err := migrate(store, tmpdir, -1)
		AssertError(t, err, ErrMigrationDirEmpty)
	})
	t.Run("non-empty directory attempts to migrate", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		err := migrate(store, tmpdir, -1, "up")
		AssertNoError(t, err)
		AssertAllStoreMigrationCalls(t, store, 1, "up")
	})
	t.Run("only apply migrations in one direction", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		err := migrate(store, tmpdir, -1, "up")
		AssertNoError(t, err)
		AssertAllStoreMigrationCalls(t, store, 1, "up")
		for name := range store.migrations {
			if strings.HasSuffix(name, "down.sql") {
				t.Errorf("Wrong direction migration applied: %s", name)
			}
		}
	})
}
```

```sh
# github.com/quii/learn-go-with-tests/databases/v2 [github.com/quii/learn-go-with-tests/databases/v2.test]
.\migrate_test.go:81:17: too many arguments in call to migrate
        have (*SpyStore, string, number, string)
        want (Storer, string, int)
FAIL    github.com/quii/learn-go-with-tests/databases/v2 [build failed]
```

The compiler is complaining, because migrate does not yet accept a direction. Let's `DRY` things a little preemtively this time. Add the following to `bookshelf-store.go`.

```go
// bookshelf-store.go
...
const (
	UP uint = iota
	DOWN
)
...
var (
	Directions = [...]string{UP: "up", DOWN: "down"}
)
...
```

Now you can access the directions by a very explicit `Directions[UP]` or `Directions[DOWN]`.

Change the signature of `migrate` to include a direction, and a check using [`strings.HasSuffix`](https://golang.org/pkg/strings/#HasSuffix).

```go
...
import (
	"strings"
)
...
func migrate(store Storer, dir string, num int, direction string) {
	...
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), direction+".sql") {
			continue
		}
		...
	}
	...
```

```sh
ok      github.com/quii/learn-go-with-tests/databases/v2        0.437s
```

Remember to modify other occurrences of `migrate` in the test file.

Looks like we accidentally covered point no. 3:

    3. It needs to run only `up` migrations.

But we are missing the 2.1 annex: it needs to be ordered. Thanfully, we can solve this with a couple of assertions and the fact that in `go`, strings are comparable.

```go
// migrate_test.go
...
import (
	...
	"sort"
)
...
func AssertOrderAscending(t *testing.T, store *SpyStore, migrations []string) {
	t.Helper()
	for i := 0; i < len(migrations) - 1; i++ {
		m0, m1 := migrations[i], migrations[i+1]
		if m0 > m1 {
			t.Errorf("wrong migration order for asc: %q before %q)", m0, m1)
		}
	}
}

func AssertOrderDescending(t *testing.T, store *SpyStore, migrations []string) {
	t.Helper()
	for i := 0; i < len(migrations) - 1; i++ {
		m0, m1 := migrations[i], migrations[i+1]
		if m0 < m1 {
			t.Errorf("wrong migration order for desc: %q before %q)", m0, m1)
		}
	}
}
```

But we don't have a way to capture migrations _in order_, we only get them from the `SpyStore`. We need to implement a return value that captures the order. This is also, indirectly helping us with point no. 5:

5. Lastly, it needs to report on the success of each migration run, if a migration fails, the entire process should be halted.

```go
// bookshelf-store.go
...
func migrate(
	store Storer,
	dir string,
	num int,
	direction string,
) ([]string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, ErrMigrationDirNoExist
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrMigrationDirEmpty
	}

	migrations := make([]string, 0)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), direction+".sql") {
			continue
		}
		path := filepath.Join(dir, file.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read migration file %s, %v", file.Name(), err)
			return nil, err
		}
		err = store.ApplyMigration(file.Name(), string(content))
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, file.Name())
	}
	return migrations, nil
}
```

Remember to change implementation of `migrate` in the test, to now return two values, discarding the first one as we don't need it.

## Write the test first

Now we can test the order of the migrations:

```go
// migrate_test.go
t.Run("up migrations should be ordered ascending", func(t *testing.T){
	store := NewSpyStore()
	tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
	defer cleanup()

	migrations, _ := migrate(store, tmpdir, -1, Directions[UP])
	AssertOrderAscending(t, store, migrations)
})

t.Run("down migrations should be ordered descending", func(t *testing.T){
	store := NewSpyStore()
	tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
	defer cleanup()

	migrations, _ := migrate(store, tmpdir, -1, Directions[DOWN])
	AssertOrderDescending(t, store, migrations)
})
```

## Try to run the test

```sh
--- FAIL: TestMigrate (0.03s)
    --- FAIL: TestMigrate/down_migrations_should_be_ordered_descending (0.00s)
        migrate_test.go:98: wrong migration order for desc: "01.993770360.down.sql" before "02.441780074.down.sql")
        migrate_test.go:98: wrong migration order for desc: "02.441780074.down.sql" before "03.668063276.down.sql")
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v2        0.624s
```

The `up` migrations are in the correct order, by grace of `ioutil.Readall`. But we should implement it explicitly, as the API for `ioutil.Readall` is not in our control, and may change and break our application.

## Write the minimal amount of code to make it pass

It's a matter of using the [`sort`](https://golang.org/pkg/sort) package to sort the `files` returned by `ioutil.ReadAll`. Specifically, [`sort.SliceStable`](https://golang.org/pkg/sort/#SliceStable). Recall that our order is implemented by the filename, we need to use `file.Name()` inside our sorting functions

```go
// bookshelf-store.go
func migrate(
	store Storer,
	dir string,
	num int,
	direction string,
) ([]string, error) {

	...

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	switch direction {
	case Directions[DOWN]:
		sort.SliceStable(files, func(i, j int) bool { return files[j].Name() < files[i].Name() })
	default:
		sort.SliceStable(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	}

	...
}
```

Now our tests pass

```sh
PASS
ok      github.com/quii/learn-go-with-tests/databases/v2        0.653s
```

On to point no. 4:

4. It needs to run migrations `1` through `num`, or all of them if `num == -1`.

For `up` migrations, `1` through `num` makes sense: you will run the `num` first migrations in the dir. But what about `down` migrations? What does `num` represent in this case?

-   The migration number that it will go down to?
-   or the number of files to process before stopping?
-   What if the number of files change before applying `down` migrations?

These are all important questions, but since we are adhering to best practices, our migration files will be _idempotent_, so it does not matter if they are run repeatedly.

`num` should be the number of files to process on `down` migrations, reporting appropriately the state of the database when done.

## Write the test first

Our `CreateTempDir` function creates 3 `up` files, and 3 `down` files if the `empty` boolan param is `false`. Since our `SpyStore` is tracking how many calls each migration receives, we can check that only the first `num` have been called directly, so let's do that.

Since we have to write two tests (one for `up` and one for `down` migrations), let's our tests DRY and write an assertion function using [`reflect`](https://golang.org/pkg/reflect). We need to take care to account for migrations that may not exist in the store (with a `0` value), so we'll have to modify the `got` slice artificially.

```go
// migrate_test.go
...
import (
	...
	"reflect"
)
...
func TestMigrate(t *testing.T) {
	...
	t.Run("runs as many migrations as the num param, up", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		migrations, _ := migrate(store, tmpdir, 2, Directions[UP])
		AssertSliceCalls(t, store, migrations, []int{1, 1, 0})
	})

	t.Run("runs as many migrations as the num param, down", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		// `migrations` slice is reversed, so desired order is still (1,1,0)
		migrations, _ := migrate(store, tmpdir, 2, Directions[DOWN])
		AssertSliceCalls(t, store, migrations, []int{1, 1, 0})
	})
}
...
func AssertSliceCalls(t *testing.T, store *SpyStore, migrations []string,want []int) {
	t.Helper()
	got := make([]int, 0)
	for _, m := range migrations {
		got = append(got, store.migrations[m].called)
	}
	for len(got) < len(want) {
		got = append(got, 0)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v calls for migrations %v", got, want, migrations)
	}
}
```

## Try to run the test

We get an explicit report of what was called and shouldn't have been:

```sh
--- FAIL: TestMigrate (0.04s)
    --- FAIL: TestMigrate/runs_as_many_migrations_as_the_num_param,_up (0.02s)
        migrate_test.go:113: got [1 1 1] want [1 1 0] calls for migrations [01.813038608.up.sql 02.463630274.up.sql 03.431250500.up.sql]
    --- FAIL: TestMigrate/runs_as_many_migrations_as_the_num_param,_down (0.01s)
        migrate_test.go:123: got [1 1 1] want [1 1 0] calls for migrations [03.900955562.down.sql 02.186370488.down.sql 01.533641238.down.sql]
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v2        0.566s
```

## Write the minimal amount of code to make them pass

We'll introduce a `count` variable, increment it when the migrations are applied, and break the loop as soon as `count >= num`.

```go
// bookshelf-store.go
...
func migrate(
	store Storer,
	dir string,
	num int,
	direction string,
) ([]string, error) {
	...
	migrations := make([]string, 0)
	count := 0
	for _, file := range files {
		if count >= num {
			break
		}
		if !strings.HasSuffix(file.Name(), direction+".sql") {
			continue
		}
		path := filepath.Join(dir, file.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read migration file %s, %v", file.Name(), err)
			return nil, err
		}
		err = store.ApplyMigration(file.Name(), string(content))
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, file.Name())
		count++
	}
	return migrations, nil
}
```

## Try to run the tests

Our tests now pass

```sh
PASS
ok github.com/quii/learn-go-with-tests/databases/v2 0.475s
```

We're missing an explicit check for `num == -1` to run all migrations.

## Write the test first

```go
// migrate_test.go
...
	t.Run("runs all migrations if num == -1", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		migrations, _ := migrate(store, tmpdir, -1, Directions[UP])
		AssertSliceCalls(t, store, migrations, []int{1, 1, 1})
	})
...
```

## Try to run the test

It fails, as the `count` variable introduced before is always greater than `-1`.

```sh
--- FAIL: TestMigrate (0.04s)
    --- FAIL: TestMigrate/runs_all_migrations_if_num_==_-1 (0.00s)
        migrate_test.go:132: got [0 0 0] want [1 1 1] calls for migrations []
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v2        1.059s
```

## Write enough code to make it pass

It's a simple fix: prepend `num != -1` to the breaking condition.

```go
// bookshelf-store.go
func migrate(
	store Storer,
	dir string,
	num int,
	direction string,
) ([]string, error) {
	...
	for _, file := range files {
		if num != -1 && count >= num {
			break
		}
		...
}
```

And we're back to green

```sh
PASS
ok     github.com/quii/learn-go-with-tests/databases/v2        0.943s
```

## Moving on

There's a lot of nuance to point number 5, so let's break it into simpler pieces:

5. Lastly, it needs to report on the success of each migration run, if a migration fails, the entire process should be halted.

> ... it needs to report on the success of each migration run...

Our migrate function currently writes no output. A simple `fmt.Println` should get the job done. But it complicates our testing, as the default output would be `os.Stdout`. We could use `go`'s utilites (namely, `os.Pipe`) to capture `os.Stdout`'s output, but this could lead to a race condition (`migrate` would be writing to `os.Stdout` as well as the `testing`). It can get hairy.

Since `migrate` will be an internal function (remember `MigrateUp`), we can hardcode `os.Stdout` inside the `migrate` call inside `MigrateUp`. Let's add an `io.Writer` parameter to `migrate`, which then we can use to safely inspect output.

But, as always, test first!

## Write the test first

Wouldn't it be helpful if migrate also let you know how many migrations are in total? This is the right time to implement this.

```go
// migrate_test.go
...
import (
	...
	"bytes"
)
...
	t.Run("success output is expected", func(t *testing.T){
		store := NewSpyStore()

		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		gotBuf := &bytes.Buffer{}
		migrations, _ := migrate(gotBuf, store, tmpdir, -1, Directions[UP])
		got := gotBuf.String()

		total := len(migrations)

		wantBuf := &bytes.Buffer{}
		current := 1
		for _, m := range migrations {
			str := fmt.Sprintf("applying %d/%d: %s ...SUCCESS\n", current, total, m)
			wantBuf.WriteString(str)
			current++
		}
		want := wantBuf.String()

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
...
```

## Try to run the test

Test fails, as expected.

```sh
# github.com/quii/learn-go-with-tests/databases/v2 [github.com/quii/learn-go-with-tests/databases/v2.test]
.\migrate_test.go:143:27: too many arguments in call to migrate
        have (*bytes.Buffer, *SpyStore, string, number, string)
        want (Storer, string, int, string)
FAIL    github.com/quii/learn-go-with-tests/databases/v2 [build failed]
```

## Write the minimal amount of code for the test to run and check the failing test output

```go
// bookshelf-store.go
...
import (
	...
	"io"
)
...
func migrate(
	out io.Writer,
	store Storer,
	dir string,
	num int,
	direction string,
) ([]string, error) {
...
```

Run the tests again

```sh
# github.com/quii/learn-go-with-tests/databases/v2 [github.com/quii/learn-go-with-tests/databases/v2.test]
.\migrate_test.go:51:20: not enough arguments in call to migrate
        have (*SpyStore, string, number, string)
        want (io.Writer, Storer, string, int, string)
.\migrate_test.go:61:20: not enough arguments in call to migrate
        have (*SpyStore, string, number, string)
        want (io.Writer, Storer, string, int, string)
.\migrate_test.go:70:20: not enough arguments in call to migrate
        have (*SpyStore, string, number, string)
        want (io.Writer, Storer, string, int, string)
...
FAIL    github.com/quii/learn-go-with-tests/databases/v2 [build failed]
```

Whoops, looks like we forgot to modify current calls to `migrate` in our tests. Let's create a `dummyWriter` variable and insert it on all but the last test we wrote.

```go
// migrate_test.go
...
...
migrate(dummyWriter, store, "i-do-not-exist", -1, Directions[UP])
...
var dummyWriter = &bytes.Buffer{}
```

Now our tests run, and yield the expected failure

```sh
--- FAIL: TestMigrate (0.08s)
    --- FAIL: TestMigrate/success_output_is_expected (0.01s)
        migrate_test.go:158: got "" want "applying 1/3: 01.734206635.up.sql ...SUCCESS\napplying 2/3: 02.168107541.up.sql ...SUCCESS\napplying 3/3: 03.660970255.up.sql ...SUCCESS\n"
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v2        0.934s
```

## Write enough code to make it pass

If `store.ApplyMigration` doesn't return an error, we can assume it was successful. This is where we'll add our message.

```go
// bookshelf-store.go
...
func migrate(
	out io.Writer,
	store Storer,
	dir string,
	num int,
	direction string,
) ([]string, error) {
	...

	total := len(files)
	if total == 0 {
		return nil, ErrMigrationDirEmpty
	}

	migrations := make([]string, 0)
	count := 0
	for _, file := range files {
		...

		fmt.Fprintf(out, "applying %d/%d: %s ", count+1, total, file.Name())
		err = store.ApplyMigration(file.Name(), string(content))
		if err != nil {
			return nil, err
		}
		fmt.Fprint(out, "...SUCCESS\n")
		migrations = append(migrations, file.Name())
		count++
	}
	return migrations, nil
}
```

## Try to run the test

Success! oh wait...

```sh
--- FAIL: TestMigrate (0.11s)
    --- FAIL: TestMigrate/success_output_is_expected (0.00s)
        migrate_test.go:158: got "applying 1/6: 01.414028183.up.sql ...SUCCESS\napplying 2/6: 02.949529057.up.sql ...SUCCESS\napplying 3/6: 03.887873211.up.sql ...SUCCESS\n" want "applying 1/3: 01.414028183.up.sql ...SUCCESS\napplying 2/3: 02.949529057.up.sql ...SUCCESS\napplying 3/3: 03.887873211.up.sql ...SUCCESS\n"
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v2        0.807s
```

## Write enough code to make it pass

It's reporting `n/6`, it should be `n/3`. Filtering inside the main loop with `strings.HasSuffix(file.Name(), direction+".sql")` is not working for us anymore.

Let's fix that by filtering the `files` outside of the main loop.

```go
// bookshelf-store.go
...
func migrate(
	out io.Writer,
	store Storer,
	dir string,
	num int,
	direction string,
) ([]string, error) {
	...
	allFiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := make([]os.FileInfo, 0)
	for _, f := range allFiles {
		if strings.HasSuffix(f.Name(), direction+".sql") {
			files = append(files, f)
		}
	}

	for _, file := range files {
		if num != -1 && count >= num {
			break
		}
		path := filepath.Join(dir, file.Name())
		...
	}
	return migrations, nil
}
```

Now our tests are passing.

```sh
PASS
ok      github.com/quii/learn-go-with-tests/databases/v2        1.129s
```

## On to 5-2

5. Lastly, it needs to report on the success of each migration run, if a migration fails, the entire process should be halted.

> ... if a migration fails, the entire process should be halted.

This already happens, as we return the error if it fails. But how do we _test_ this?

This is prime material for the integration tests, as the database engine itself will tell you if the `sql` inside the migration file is bad or cannot be executed. And we will get to this.

We need to simulate a failure. We can do this within our `SpyStore`.

Given that we will test this anyway (via the integration tests), we don't have to parse any `sql` inside the migration files; this is a whole different beast that we fortunately don't have to deal with.

Since we used an interface (`Storer`) to abstract implementation, we can put whatever we want inside `SpyStore.ApplyMigration` for the failure condition.

Pies vs Cakes is a very serious debate that has raged on since forever. Pie is clearly superior. So we decided to forbid any cake-related SQL. Our failure condition will be the word `cake` inside the `sql` files.

Let's move on to our tests

## Write the test first

The failure should be reported with a helpful error message as well, so as not to leave our users helplessly looking through their files.

We will ignore errors for the sake of brevity.

```go
// migrate_test.go
...
import (
	...
	"errors"
)
...
var errNoCakeSQL = errors.New("cakeSQL is not allowed")
...
func (s *SpyStore) ApplyMigration(name, stmt string) error {
	if strings.Contains(strings.Lower(stmt), "cake") {
		return errNoCakeSQL
	}
	// the rest of the method is unchanged
	...
}
...
	t.Run("failure output is expected", func(t *testing.T){
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", true)
		defer cleanup()

		tmpfile, _ := ioutil.TempFile(tmpdir, "01.cake.*.up.sql")
		tmpfile.Write([]byte("cake is superior! end pie tyranny")); err != nil {
		tmpfile.Close()

		gotBuf := &bytes.Buffer{}
		_, err = migrate(gotBuf, store, tmpdir, -1, Directions[UP])
		got := gotBuf.String()

		wantBuf := &bytes.Buffer{}
		str := fmt.Sprintf(
			"applying 1/1: %s ...FAILURE: %v\n",
			filepath.Base(tmpfile.Name()),
			errNoCakeSQL,
		)
		wantBuf.WriteString(str)
		want := wantBuf.String()

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
		AssertError(t, err, errNoCakeSQL)
	})
```

## Try to run the test

As expected, since there is no code for the failure report yet.

```sh
--- FAIL: TestMigrate (0.07s)
    --- FAIL: TestMigrate/failure_output_is_expected (0.00s)
        migrate_test.go:191: got "applying 1/1: 01.cake.702313345.up.sql " want "applying 1/1: 01.cake.702313345.up.sql ...FAILURE: cakeSQL is not allowed\n"
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v2        0.647s
```

## Write enough code to make it pass

```go
// bookshelf-store.go
...
func migrate(
	out io.Writer,
	store Storer,
	dir string,
	num int,
	direction string,
) ([]string, error) {
	...
	for _, file := range files {
		...
		err = store.ApplyMigration(file.Name(), string(content))
		if err != nil {
			fmt.Fprintf(out, "...FAILURE: %v\n", err)
			return nil, err
		}
		...
	}
	return migrations, nil
}
```

And we are green again

```sh
PASS
ok      github.com/quii/learn-go-with-tests/databases/v2        1.555s
```

## Modifying the database

<!-- TODO: add links start:v2, end:v3 -->

The moment of truth, now that our `migrate` function behaves as expected, we can implement it for our integration tests!

But, like everything we've done before, we need to test it as well.

Integration tests are, in some aspects, easier to implement than unit tests, as you can rely on the service error reporting to test your application (though you want to maintain as much control as possible).

Before we move on, let's write two wrapper functions around our `migrate` powerhouse.

```go
// bookshelf-store.go
func MigrateUp(out io.Writer, store Storer, dir string, num int) ([]string, error) {
	return migrate(out, os.Stdout, store, dir, num, Directions[UP])
}

func MigrateDown(out io.Writer, store Storer, dir string, num int) ([]string, error) {
	return migrate(out, os.Stdout, store, dir, num, Directions[DOWN])
}
```

## Write the test first

We're going to rely on the failure of `sql.DB`'s [`Exec`](https://golang.org/pkg/database/sql/#DB.Exec) method, which returns an error if the operation could not complete. We can also use some `PostgreSQL` utilities to list the tables from within our application, and examine the output to verify our tables are there.

Recall our migration written earlier from `migrations/0001_create_books_table.up.sql`. This time we will add some comments to explain what each line does.

Remember, in `SQL`, everything after a double hyphen (`--`) is a comment.

```sql
BEGIN; -- BEGIN starts a transaction

-- create a table named `books`
CREATE TABLE IF NOT EXISTS books (
	-- add a column named `id`, which is the PRIMARY KEY
	-- type SERIAL means it's an auto-incrementing sequence
	-- of values
	id SERIAL PRIMARY KEY,
	-- add a column named `title`, which is of type VARCHAR (string)
	-- and a maximum allocation of 255 characters. It cannot be NULL
	title VARCHAR(255) NOT NULL,
	-- add a column named `author`, same as `title`
	author VARCHAR(255) NOT NULL
);

-- CREATEs an index named `books_id_uindex` on the table `books`
-- using the column `id`
CREATE UNIQUE INDEX IF NOT EXISTS  books_id_uindex ON books (id);

COMMIT; -- COMMIT saves the transaction to the database
```

Key points

-   **Transaction**:

    A transaction represents a batch of work that needs to be performed together. A transaction has to be started with `BEGIN`, and is either saved with `COMMIT` or all the work done so far reversed with `ROLLBACK`.

    We don't use `ROLLBACK` in our migration because our simple creation of tables and index is very unlikely to fail, and it is safeguarded by the `IF NOT EXISTS` clause, which does nothing if the table or index already exists.

-   **Index**:

    A database index is, in the layman's terms, a trade-off that improves the retrieval of information (if done right) by giving a little more every time data is added.

    Here are more formal definitions if the subject interests you: [Wikipedia](https://en.wikipedia.org/wiki/Database_index), [Use The Index, Luke](https://use-the-index-luke.com/sql/anatomy).

    In this case, the index is redundant, as `PostgreSQL` creates an index on the `PRIMARY KEY` of a table by default.

-   **SQL Language**

    The `SQL` language is part of an ISO standard, and most database engines comform to it partially. This means that code written for one `RDBMS` (say, `PostgreSQL`), will cannot be interpreted as-is by a different one (say, `SQLite3`). There are a lot of similarities, however, and the changes are often small.

    Keep in mind that the `SQL` you're seeing here is very `PostgreSQL` specific, and some, if not all of it, may not be executable in a different engine.

The test! Here is our first integration test.

```go
//integration_test.go
package main

import "testing"

func TestMigrations(t *testing.T) {
	store, removeStore := NewStore()
	defer removeStore()

	t.Run("migrate up", func(t *testing.T){
		_, err := MigrateUp(dummyWriter, store, "migrations", -1)
		if err != nil {
			t.Errorf("migration up failed: %v", err)
		}
	})
	t.Run("migrate down", func(t *testing.T){
		_, err := MigrateDown(dummyWriter, store, "migrations", -1)
		if err != nil {
			t.Errorf("migration down failed: %v", err)
		}
	})
}
```

## Try to run the test

I have deliberately turned off my `PostgreSQL` server to get an error.

If you'd like to do the same, run `sudo systemctl stop postgresql.service` in a shell, or kill your `docker` container of `postgres`.

```sh
--- FAIL: TestMigrations (2.68s)
    --- FAIL: TestMigrations/migrate_up (1.37s)
        integration_test.go:14: migration up failed: dial tcp 127.0.0.1:5432: connectex: No connection could be made because the target machine actively refused it.
    --- FAIL: TestMigrations/migrate_down (1.31s)
        integration_test.go:20: migration down failed: dial tcp 127.0.0.1:5432: connectex: No connection could be made because the target machine actively refused it.
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v3        4.073s
```

This output is one of many possible errors we may get. While we cannot control this type of error in the real world, this is useful to know in case we find it "in the wild". This is why we do `integration tests` in the first place!

Restart your `postgres` `docker` instance or run `sudo systemctl start postgresql.service`, and try again.

```sh
--- FAIL: TestMigrations (6.21s)
    --- FAIL: TestMigrations/migrate_up (6.15s)
        integration_test.go:14: migration up failed: pq: SSL is not enabled on the server
    --- FAIL: TestMigrations/migrate_down (0.06s)
        integration_test.go:20: migration down failed: pq: SSL is not enabled on the server
FAIL
exit status 1
FAIL    github.com/quii/learn-go-with-tests/databases/v3        12.427s
```

I admit that I deliberately left out an important part of the connection string, the query parameter `sslmode=disable`. Partly to get this error, partly to explain `SSL`.

Databases are generally used over networks, and, like all network connections, they should be `secured` if they have sensitive data. One of the security measures `PostgreSQL` can implements is **S**ecure **S**ocket **L**ayer, or **SSL**. It allows encrypted connections to and from the database.

Our database lives locally, so it would be redundant to implement encryption here.

Change the `connStr` constant inside `NewStore` to include the query parameter `sslmode=disable`.

```go
// bookshelf-store.go
...
const connStr = "postgres://bookshelf_user:secret-password@localhost:5432/bookshelf_db?sslmode=disable"
...
```

## Try to run the tests

Now our tests pass

```sh
PASS
ok      github.com/quii/learn-go-with-tests/databases/v3        12.989s
```

You probably noticed that our test are much slower now (~`12s` vs ~`1s` before).

That said, such is the nature of integration tests: testing between services requires more computing power, and has to account for things like latency, message queues and other nuisances that add to the test time.

It's not entirely hopeless though, as a solution exists! It's along the lines of "run the tests on someone else's computer".

Actually, it's exactly like "run the tests on someone else's computer", and it's called contionuous integration (commonly referred to as CI).

We won't cover CI in this chapter, but we'll point to some resources at the end. For the moment, we'll have to bite the bullet and endure the slow tests.

## Where are the tables?

So far our tests pass, and we assume they do what they're supposed to. But we're modifying a database, you would think there is _something_ happening somewhere that makes said modifications. And you would be right.

Let's extend our tests to ensure that we are getting some output.

```go
// integration_tests.go
...
const queryTables = `
SELECT tablename, tableowner
FROM pg_catalog.pg_tables
WHERE
	schemaname != 'pg_catalog'
	AND
	schemaname != 'information_schema';`

type pgTable struct {
	tableOwner string `sql:"tableowner"`
	tableName string `sql:"tablename"`
}

func TestMigrations(t *testing.T) {
	store, removeStore := NewStore()
	defer removeStore()

	t.Run("migrate up", func(t *testing.T) {
		_, err := MigrateUp(dummyWriter, store, "migrations", -1)
		if err != nil {
			t.Errorf("migration up failed: %v", err)
		}

		rows, err := store.db.Query(queryTables)
		if err != nil {
			t.Errorf("received error querying rows: %v",  err)
			t.FailNow()
		}
		defer rows.Close()

		tables := make([]pgTable, 0)
		for rows.Next() {
			var table pgTable
			if err := rows.Scan(&table.tableName, &table.tableOwner); err != nil {
				t.Errorf("error scanning row: %v", err)
				continue
			}
			tables = append(tables, table)
		}
		if err := rows.Err(); err != nil {
			t.Errorf("rows error: %v", err)
		}

		set := make(map[string]bool)
		for _, table := range tables {
			set[table.tableName] = true
		}

		if _, ok := set["books"]; !ok {
			t.Error("table \"books\" not returned")
		}
	})
	t.Run("migrate down", func(t *testing.T) {
		_, err := MigrateDown(dummyWriter, store, "migrations", -1)
		if err != nil {
			t.Errorf("migration down failed: %v", err)
		}

		rows, err := store.db.Query(queryTables)
		if err != nil {
			t.Errorf("received error querying rows: %v",  err)
			t.FailNow()
		}
		defer rows.Close()

		got := 0
		for rows.Next() {
			got++
		}
		if err := rows.Err(); err != nil {
			t.Errorf("rows error: %v", err)
		}
		if got > 0 {
			t.Errorf("got %d want 0 rows", got)
		}
	})
}
```

And our tests still pass, but this time we are sure of what our code does.

```sh
PASS
ok      github.com/quii/learn-go-with-tests/databases/v3        0.860s
```

## What is going on?!

If you're not faimilar with `sql` and how `go` handles it, you are probably confused right now. Let's break down the code and explain it in smaller chunks.

```go
const queryTables = `
SELECT tablename, tableowner
FROM pg_catalog.pg_tables
WHERE
	schemaname != 'pg_catalog'
	AND
	schemaname != 'information_schema';`
```

Here we're creating a constant string to `query` the desired fields from one of `PostgreSQL`'s administrative databases: `pg_catalog`, from the table `pg_tables`. The `WHERE` clause is to filter out `PostgreSQL`'s own tables.

```go
type pgTable struct {
	tableOwner string `sql:"tableowner"`
	tableName string `sql:"tablename"`
}
```

This struct will hold table information. The `tags` `sql:"tableowner"` and `sql:"tablename"`, tell the `database/sql` package utilities to map those columns to those fields.

For the `migrate up` test, we'll add comments on each line instead

```go
	t.Run("migrate up", func(t *testing.T) {
		_, err := MigrateUp(dummyWriter, store, "migrations", -1)
		if err != nil {
			t.Errorf("migration up failed: %v", err)
		}

		// this executes the queryTables query defined above
		rows, err := store.db.Query(queryTables)
		if err != nil {
			t.Errorf("received error querying rows: %v",  err)
			t.FailNow()
		}
		// prevent memory leaks
		defer rows.Close()

		// create a slice to hold our testable information
		tables := make([]pgTable, 0)
		// iterates through the `rows`, one by one
		for rows.Next() {
			var table pgTable
			// this scan the columns on each row, mapping them to
			// the pgTable created above
			if err := rows.Scan(&table.tableName, &table.tableOwner); err != nil {
				t.Errorf("error scanning row: %v", err)
				continue
			}
			tables = append(tables, table)
		}
		if err := rows.Err(); err != nil {
			t.Errorf("rows error: %v", err)
		}

		// we create a set of unique values using the table names
		set := make(map[string]bool)
		for _, table := range tables {
			set[table.tableName] = true
		}

		// if the table we want (books) is not in the set, then it
		// was not returned by the query, thus, should fail
		if _, ok := set["books"]; !ok {
			t.Error("table \"books\" not returned")
		}
	})
```

The `migrate down` tests uses very similar logic to the `migrate up` test, but instead it counts the rows returned. The number of rows should be `0`, because we `down`-migrated all the tables.

## Best practices

Recall that earlier we mentioned that our migrations should follow best practices:

-   Migrations should be _ordered_. ...
-   Migrations should be _reversible_. ...
-   Migrations should be _idempotent_. ...

We tested that they were _ordered_ with our unit tests, and that they're _reversible_ with our integration tests. But are they _idempotent_?

Let's write a test for it! Simply run the same migrations twice, the database should raise an error if it doesn't allow the same migration twice.

## Write the test first

We will simply run a bunch of migrations, with no particular order. The database should raise an error if there is a problem.

```go
// integration_test.go
func TestMigrations(t *testing.T) {
...
	t.Run("idempotency", func(t *testing.T){
		_, err := MigrateDown(dummyWriter, store, "migrations", -1)
		if err != nil {
			t.Errorf("first migrate down failed: %v", err)
		}

		_, err = MigrateUp(dummyWriter, store, "migrations", -1)
		if err != nil {
			t.Errorf("first migrate up failed: %v", err)
		}

		_, err = MigrateUp(dummyWriter, store, "migrations", -1)
		if err != nil {
			t.Errorf("second migrate up failed: %v", err)
		}

		_, err = MigrateDown(dummyWriter, store, "migrations", -1)
		if err != nil {
			t.Errorf("second migrate down failed: %v", err)
		}

		_, err = MigrateDown(dummyWriter, store, "migrations", -1)
		if err != nil {
			t.Errorf("third migrate down failed: %v", err)
		}
	})
}
```

## Try to run the test

We pass!

```sh
PASS
ok      github.com/quii/learn-go-with-tests/databases/v3        1.384s
```

## What we've accomplished so far

Up to this point, we have managed to create our migration tool, and have tested it thoroughly.

This tool is, however, secondary to our actual goals. Now we'll be moving on to create a `CRUD` application to interact with our database.

In case you are wondering, `CRUD` is an acronym for **C**reate, **R**etrieve, **U**pdate, **D**elete, often used to describe the basic operations needed to run against a storage system.

But first...

## Housekeeping

<!-- TODO: add links start:v3, end:v4 -->

Our application is growing, and so is our codebase. Before we are drowing in `.go` files, let's use `go`'s package structure to our advantage.

Let's create a package (directory), aptly name it `bookshelf`, and put our `migrate` related code (`bookshelf-store.go`) inside it. Don't forget to change occurrences of `package main` to `package bookshelf`.

Also change all test files (`migrate_test`, `integration_test`) `package main` to `package bookshelf_test`, and put them inside the `bookshelf` directory. You will have to import the `bookshelf` package you created into these tests.

Inside bookshelf, the `migrate` function will have to be exported (change to `Migrate`).

Inside bookshelf, an utility package called `testutils`, that will hold all the testing utilities (duh!) we have created so far: put assertions into a file called `assertions.go`, `CreateTempDir` inside `helpers.go` and `SpyStore` and its related functions and methods inside `store.go`.

Inside `migrate_test.go`, all occurrences of `migrate` will need to be changed to `bookshelf.Migrate`. You will need to import the `bookshelf/testutils` package and prepend all `assertions` and occurrences of `CreateTempDir` and `NewSpyStore` with `testutils.*`.

Finally, move the `migrations` directory inside the `bookshelf` dir.

---

### Note on the package import string:

To me, the import string of my `bookshelf` package looks like this:

```go
import (
	...
	"github.com/djangulo/learn-go-with-tests/databases/v4/bookshelf"
)
```

It's likely you did not clone the repository to go through this chapter, and that's OK.

Keep in mind that you may need to change your import string to something like

```go
import (
	...
	"github.com/YOUR_GITHUB_HANDLE/bookshelf"
)
```

---

Once you're done, folder structure should look like this:

```sh
.
└── bookshelf
    ├── bookshelf-store.go
    ├── integration_test.go
    ├── migrate_test.go
    ├── migrations
    │   ├── 0001_create_books_table.down.sql
    │   └── 0001_create_books_table.up.sql
    └── testutils
        ├── assertions.go
        └── store.go

3 directories, 7 files
```

Try running the tests, the compiler will tell you what to do. Keep correcting the errors, eventually they will run out, I promise.

This exercise in patience may seem pointless now, but it's well worth the effort.

## CRUD

<!-- TODO: add links start:v4, end:v5 -->

We have `4` operations to write and test. A lot of the code already in place helps us, so the workload will be lighter than before (we hope so at least).

In the `SQL` migration, We created a table called `books`, with columns called `id`, `title` and `author`. Let`s create a struct to hold these objects before we get into testing.

```go
// bookshelf/bookshelf-store.go
...
type Book struct {
	ID     int64  `sql:"id"`
	Title  string `sql:"title"`
	Author string `sql:"author"`
}
...
```

## Write the test first

Before we can retrieve, update or delete an object, we need to create it first! Logically, it makes sense to start here.

Change the `SpyStore` struct in `bookshelf/testutils/store.go` to include a slice of books.

```go
// bookshelf/testutils/store.go
package testutils

import (
	...
	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf"
)
...
type SpyStore struct {
	Migrations map[string]migration
	Books []*bookshelf.Book
}
...
func NewSpyStore() *SpyStore {
	books := make([]*bookshelf.Book, 0)
	return &SpyStore{
		Migrations: map[string]migration{},
		Books: books,
	}
}
```

And our test

```go
// bookshelf/crud_test.go
package bookshelf_test

import (
	"testing"

	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf"
	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf/testutils"
)

func TestCreate(t *testing.T) {
	store := testutils.NewSpyStore()

	var book bookshelf.Book
	err := store.CreateBook(&book, "Moby Dick", "Herman Melville")
	testutils.AssertNoError(t, err)
	if book.ID == 0 {
		t.Error("book returned without an ID")
	}
}

```

Notice that the `CreateBook` method does not have an `ID` field provided, this is because `primary keys` are usually autoincremented and provided by the database.

This `ID` field is our criteria for a passing test.

This is not written in stone, however: just happens that when we created the database table, we designated the `id` field as `SERIAL`, and `PostgreSQL` handles the auto-incrementing for us. But we could have designated a `PRIMARY KEY` of whichever type we would've wanted. For example, had we assigned the `title` as the primary key, the `id` field would have not been necessary. Or, if we designated `id` as `PRIMARY KEY`, but as type `INT` instead of `SERIAL`, our application would have had to find the latest `id` and increment it.

## Try to run the test

```sh
~$ $go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf_test [github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf.test]
databases\v5\bookshelf\crud_test.go:14:14: store.CreateBook undefined (type *"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf/testutils".SpyStore has no field or method CreateBook)
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf [build failed]
FAIL
```

## Write the minimal amount of code for the test to run and check the failing test output

Write the `CreateBook` method for the `SpyStore`.

```go
// bookshelf/testutils/store.go
...
func (s *SpyStore) CreateBook(book *bookshelf.Book, title, author string) error {
	return nil
}
...
```

When we run the tests again

```sh
~$ go test ./bookshelf/
--- FAIL: TestCreate (0.00s)
    crud_test.go:17: book returned without an ID
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  4.635s
FAIL
```

## Write enough code to make it pass

We should check that the passed in book has a valid (i.e. non-zero) `id`, and that the `Title` and `Author` match what was passed in. We should ensure our `id`'s are unique as well (as will the database).

As a side note, `PostgreSQL`'s behavior regarding the `SERIAL` type is to keep autoincrementing it no matter what, even if previous `id`s are gone. For example, if you create `books` with `id`s 1, 2, 3 and 4, and delete the `book` with `id` 2, the next `id` created will be `5`. You could reassign it should you want to, but this leads to confusion and it's generally bad practice.

Let's use a helper to find the last ID before we assign it to the book.

```go
// bookshelf/testutils/store.go
func (s *SpyStore) CreateBook(book *bookshelf.Book, title, author string) error {
	book.ID = newID(s)
	book.Title = title
	book.Author = author
	s.Books = append(s.Books, book)
	return nil
}

func newID(store *SpyStore) int64 {
	if len(store.Books) == 0 {
		return 1
	}
	var last int64
	for _, b := range store.Books {
		if b.ID > last {
			last = b.ID
		}
	}
	return last + 1
}
```

## Try to run the test

```sh
~$ go test ./bookshelf
ok      github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  1.384s
```

## But wait, what if we create the same book twice?

Our tests pass, but what about duplicated data? We don't need two, three or twenty entries of the same book.

Turn's out we have made a mistake when creating the tables in `migrations/0001_create_books_table.up.sql`. While we designated the `author` and the `title` columns to be required (`NOT NULL`, we should test against this too!), we did not designate the `title` as `UNIQUE`. We should be careful, as there could be different `author`s with books that share a `title`; this may seem like an edge case, but edge cases is one of many reasons why we test!

While you might be tempted to just change the `0001_create_books_table.up.sql` file, you shouldn't! You might break the production database by altering the existing migrations!

The proper way to do this is to add a new migration (`up` and `down` as well), that modifies the database to the desired behavior.

Let's start there. Create two new files: `migrations/0002_books_unique_title.up.sql` and `migrations/0002_books_unique_title.down.sql`.

```sql
-- migrations/0002_books_unique_title.up.sql
ALTER TABLE IF EXISTS books ADD CONSTRAINT books_unique_author_title UNIQUE (author, title);

```

```sql
-- migrations/0002_books_unique_title.down.sql
ALTER TABLE IF EXISTS books DROP CONSTRAINT IF EXISTS books_unique_author_title;
```

We should test our migrations before we move on, our integration tests should tell us if our `SQL` is correct.

```sh
~$ go test ./bookshelf
--- FAIL: TestMigrations (0.46s)
    --- FAIL: TestMigrations/idempotency (0.23s)
        integration_test.go:105: second migrate up failed: pq: relation "books_unique_author_title" already exists
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  0.776s
FAIL
```

It's complaining that the constraint already exists. Unfortunately, there is no handy `IF NOT EXISTS` for constraint creation on `PostgreSQL`. We have two ways to go about this:

1. Create a `PostgreSQL` function using the scripting language provided by `PostgreSQL`.
2. Or, drop the constraint before creation, to ensure it runs without a hitch.

Being honest here, option 2 is bad. This exposes your system to exist without the constraint, even for a few milliseconds, bad things could happen. Not to mention the cost of the unnecessary write operation.

But, since this is not an `SQL` book, we're going to opt for the easier of the two, that is, option 2. If this were a real application, option 1 would be the choice without question. If you still want to go this way, search online for "postgresql add constraint if not exists", answers abound.

Modify `0002_books_unique_title.up.sql` to drop the constraint just before creating it. Let's wrap it in a transaction, so at least it runs as a unit.

```sql
-- migrations/0002_books_unique_title.up.sql
BEGIN;
ALTER TABLE IF EXISTS books DROP CONSTRAINT IF EXISTS books_unique_author_title;
ALTER TABLE IF EXISTS books ADD CONSTRAINT books_unique_author_title UNIQUE (author, title);
COMMIT;
```

And run the tests again

```sh
~$ go test ./bookshelf
ok      github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  0.502s
```

Now our tests pass.

Take a moment to appreciate how our `go` integration tests just helped us test our that our `SQL` is up to our standards.

We've been to focused on the unit tests. We need to create integration tests for our `Createbook` function as well.

We need to implement it in the `Store` as well.

As usual, let's start with the test.

## Write the test first

```go
// bookshelf/integration_test.go
...
func TestCreateBook(t *testing.T) {
	store, removeStore := bookshelf.NewStore()
	defer removeStore()

	t.Run("can create a book", func(t *testing.T) {
		var book bookshelf.Book
		err := store.CreateBook(&book, "test-title", "test-author")
		if err != nil {
			t.Errorf("received error on CreateBook: %v", err)
		}
		if book.ID == 0 {
			t.Error("invalid ID received")
		}
	})

	t.Run("cannot create a duplicate title-author", func(t *testing.T) {
		var b1, b2 bookshelf.Book
		err := store.CreateBook(&b1, "test-title", "test-author")
		if err != nil {
			t.Errorf("received error on CreateBook: %v", err)
		}

		err = store.CreateBook(&b2, "test-title", "test-author")
		if err == nil {
			t.Error("wanted an error but didn't get one")
		}
		
	})

}

```

## Try to run the test

```sh
~$ go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf_test [github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf.test]
bookshelf\integration_test.go:131:15: store.CreateBook undefined (type *bookshelf.Store has no field or method CreateBook)
bookshelf\integration_test.go:142:15: store.CreateBook undefined (type *bookshelf.Store has no field or method CreateBook)
bookshelf\integration_test.go:147:14: store.CreateBook undefined (type *bookshelf.Store has no field or method CreateBook)
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf [build failed]
FAIL
```

## Write the minimal amount of code for the test to run and check the failing test output

Fair enough, no method exists. Let's add the signature to the `Storer` interface, as well as the method to the `Store`.

```go
// bookshelf/bookshelf-store.go
...
type Storer interface {
	ApplyMigration(name, stmt string) error
	CreateBook(*Book, string, string) error
}
...
func (s *Store) CreateBook(book *Book, title, author string) error {
	return nil
}
...
```

## Try to run the test

We receive our expected failure.

```sh
~$ go test ./bookshelf
--- FAIL: TestCreateBook (0.05s)
    --- FAIL: TestCreateBook/can_create_a_book (0.00s)
        integration_test.go:136: invalid ID received
    --- FAIL: TestCreateBook/cannot_create_a_duplicate_title-author (0.00s)
        integration_test.go:149: wanted an error but didn't get one
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  0.730s
FAIL
```

## Write enough code to make it pass

We can now implement it.

```go
// bookshelf/bookshelf-store.go
...
// CreateBook inserts a new Book into the database.
func (s *Store) CreateBook(book *Book, title, author string) error {
	stmt := "INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id, title, author;"
	row := s.DB.QueryRow(stmt, title, author)

	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return err
	}
	return nil
}
```

Run the tests again:

```sh
~$ go test ./bookshelf
--- FAIL: TestCreateBook (0.11s)
    integration_test.go:127: received error on CreateBook: pq: relation "books" does not exist
    integration_test.go:130: invalid ID received
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  1.031s
FAIL
```

Huh, "relation `books` does not exist". By now, you probably know this means the table `books` does not exist in the database.

Well, we did not migrate before we tested, so that makes sense.

Add a call to `MigrateUp`, and check the error, after acquiring the `store`.

```go
// bookshelf/integration_test.go
func TestCreateBook(t *testing.T) {
	store, removeStore := bookshelf.NewStore()
	defer removeStore()

	_, err := bookshelf.MigrateUp(dummyWriter, store, "migrations", -1)
	if err != nil {
		t.Errorf("migration up failed: %v", err)
		t.FailNow()
	}
	...
}
```

## Try to run the tests

```sh
~$ go test ./bookshelf
--- FAIL: TestCreateBook (0.20s)
    --- FAIL: TestCreateBook/cannot_create_a_duplicate_title-author (0.03s)
        integration_test.go:144: received error on CreateBook: pq: duplicate key value violates unique constraint "books_unique_author_title"
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  1.197s
FAIL
```

Is it the error we expect? Yes it is. But an error should be _passing_, so what is going on?

The test raising the error is the first `CreateBook` call inside `cannot create a duplicate title-author`. This is not what we planned!

As it turns out, there is only 1 database connection string in our application; the one that lives inside `NewStore`. This means that all along we've been migrating `up` and `down`, inserting and deleting into a single database! This is a very risky practice, as our tests may modify or delete sensitive data once we're running in production.

So what do we do?

## Test database

So far, we have been operating in the `bookshelf_db` database, that we created at the start of the chapter. We need a secondary database that we can test to our heart's content.

Our options are:

1. Write some `go` code that creates a test database on the fly. Runs the tests and drops it once it's done.
2. Create a test database outside our application (using `psql`), and hardcode the address in our application.

Both approaches have their downsides:

1. The first approach requires more `go` code to write, and depends on the privileges that the DB user (`bookshelf_user`, in our case) has. When created the database, we gave our user `bookshelf_user` the capacity to create databases with `CREATEDB`.
2. The second approach is simpler to implement, but then our tests depend on the existence of said test database. It also implies that we need to track 2 connections inside our application, as opposed to just 

We will opt for the first approach, and take advantage that our `MigrateDown` function cleans the database tables, due to the `CASCADE` statement at the end of the `down` migrations.

## Refactor

Now we need to refactor our code. Stop and think about what's going on inside `NewStore`.

```go
// bookshelf/bookshelf-store.go
func NewStore() (*Store, func()) {
	const connStr = "postgres://bookshelf_user:secret-password@localhost:5432/bookshelf_db?sslmode=disable"
	...
}
```

We've hard-coded the connection string into our function that was soupposed to give us flexibility!

Let's create some utilities to help ourselves.

Inside `bookshelf-store.go`, insert the following:

```go
//bookshelf/bookshelf-store.go
...
type DBConf struct {
	User, Pass, Host, Port, DBName,	SSLMode string
}

func (d *DBConf) String() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Pass, d.Host, d.Port, d.DBName, d.SSLMode)
}

func getenv(key, defaultValue string) string {
	envvar := os.Getenv(key)
	if envvar == "" {
		return defaultValue
	}
	return envvar
}

var MainDBConf DBConf
func init() {
	MainDBConf.User = getenv("POSGRES_USER", "bookshelf_user")
	MainDBConf.Pass = getenv("POSTGRES_PASSWORD", "secret-password")
	MainDBConf.Host = getenv("POSTGRES_HOST", "localhost")
	MainDBConf.Port = getenv("POSTGRES_PORT", "5432")
	MainDBConf.DBName = getenv("POSTGRES_DB", "bookshelf_db")
	MainDBConf.SSLMode = getenv("POSTGRES_SSLMODE0", "disable")
}
```
With the code above, we can choose the database we want to connect to via environment variables. The `getenv` function is a simple extension of [`os.Getenv`](https://golang.org/pkg/os/#Getenv) that provides save defaults in case the variables are not set.

The code inside the `init` function will run every time the package is called, so the `MainDBConf` will be instantiated and ready when the `bookshelf` package is imported.

Our `NewStore` function now looks like this:

```go
// bookshelf/bookshelf-store.go
func NewStore(conf *DBConf) (*Store, func()) {

	db, err := sql.Open("postgres", conf.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database %q: %v\n", conf.DBName, err)
		os.Exit(1)
	}

	// exponential backoff
	remove := func() {
		deadline := time.Now().Add(removeTimeout)
		for tries := 0; time.Now().Before(deadline); tries++ {
			err := db.Close()
			retryIn := time.Second << uint(tries)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"error closing connection to database %q, retrying in %v: %v\n",
					conf.DBName,
					retryIn,
					err,
				)
				time.Sleep(retryIn)
				continue
			}
			return
		}
		log.Fatalf("timeout of %v exceeded", removeTimeout)
	}

	return &Store{DB: db}, remove
}
```

With this tooling in place, we can create a helper function to instantiate a new database just for our tests. Let's get to it!

Insert the function below inside `bookshelf/testutils/helpers.go`:

```go
// bookshelf/testutils/helpers.go
package testutils
import (
	...
	"database/sql"
	"io/ioutil"
	"time"

	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf"
	_ "github.com/lib/pq"
)
...
func NewTestStore(conf *bookshelf.DBConf) (*bookshelf.Store, func(), error) {
	main, removeMain := bookshelf.NewStore(&bookshelf.MainDBConf)

	_, err := main.DB.Exec(
		fmt.Sprintf("CREATE DATABASE %s OWNER %s;",
			conf.DBName,
			bookshelf.MainDBConf.User,
		),
	)
	if err != nil {
		return nil, nil, err
	}

	testDB, err := sql.Open("postgres", conf.String())
	if err != nil {
		return nil, nil, err
	}

	remove := func() {
		closeDeadline := time.Now().Add(5 * time.Second)
		dropDeadline := time.Now().Add(10 * time.Second)
		for tries := 0; time.Now().Before(closeDeadline); tries++ {
			retryIn := time.Second << uint(tries)
			err := testDB.Close()
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"error closing test database %q, retrying in %v: %v\n",
					conf.DBName,
					retryIn,
					err,
				)
				time.Sleep(retryIn)
				continue
			}
			break
		}
		for tries := 0; time.Now().Before(dropDeadline); tries++ {
			retryIn := time.Second << uint(tries)
			_, err := main.DB.Exec(fmt.Sprintf("DROP DATABASE %s;", conf.DBName))
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"error dropping test database %q, retrying in %v: %v\n",
					conf.DBName,
					retryIn,
					err,
				)
				time.Sleep(retryIn)
				continue
			}
			break
		}
		removeMain()
	}
	return &bookshelf.Store{DB: testDB}, remove, nil
}
```

Remember `exponential backoff`? This is where this pattern shines in our codebase. If the `mainDB` tries to drop the test database, but it's being written by a test, the operation will fail. With exponential backoff, it'll give the running operations a bit of time to finish, and then finally drop the test database.

We now can create the test database inside each test function, run all our tests with a predictable state, and drop it once we're done. We can use the fact that `MigrateDown` clears the database to our advantage and clean it after each test.

However, to avoid creating the test database multiple times, let's group our `TestCreateBook` and `TestMigrate` into a single function. We can still get meaningful reporting by nesting them with `t.Run`.

There is nothing stopping us from creating as many test databases as we want, but each database created will make our tests that much slower.

Let's create another test utility to reset the database on a whim.

```go
// bookshelf/testutils/helpers.go
func ResetStore(store *bookshelf.Store) error {
	var err error
	_, err = bookshelf.MigrateDown(dummyWriter, store, "migrations", -1)
	if err != nil {
		return err
	}

	_, err = bookshelf.MigrateUp(dummyWriter, store, "migrations", -1)
	if err != nil {
		return err
	}

	return nil
}
```

And finally, our integration tests.

```go
// bookshelf/integration_test.go
...
import (
	...
	"os"
	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf/testutils"
)
...
var (
	dbconf = bookshelf.DBConf{
		User:    bookshelf.MainDBConf.User,
		Pass:    bookshelf.MainDBConf.Pass,
		Host:    bookshelf.MainDBConf.Host,
		Port:    bookshelf.MainDBConf.Port,
		DBName:  "bookshelf_test_db",
		SSLMode: bookshelf.MainDBConf.SSLMode,
	}
)

func TestDBIntegration(t *testing.T) {

	store, removeStore, err := testutils.NewTestStore(&dbconf)
	if err != nil {
		panic(err)
	}
	defer removeStore()

	t.Run("Migrate", func(t *testing.T) {

		t.Run("migrate up", func(t *testing.T) {
			_, err := bookshelf.MigrateUp(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("migration up failed: %v", err)
			}

			rows, err := store.DB.Query(queryTables)
			if err != nil {
				t.Errorf("received error querying rows: %v", err)
				t.FailNow()
			}
			defer rows.Close()

			tables := make([]pgTable, 0)
			for rows.Next() {
				var table pgTable
				if err := rows.Scan(&table.tableName, &table.tableOwner); err != nil {
					t.Errorf("error scanning row: %v", err)
					continue
				}
				tables = append(tables, table)
			}
			if err := rows.Err(); err != nil {
				t.Errorf("rows error: %v", err)
			}

			set := make(map[string]bool)
			for _, table := range tables {
				set[table.tableName] = true
			}

			if _, ok := set["books"]; !ok {
				t.Error("table books not returned")
			}
		})
		t.Run("migrate down", func(t *testing.T) {
			_, err := bookshelf.MigrateDown(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("migration down failed: %v", err)
			}

			rows, err := store.DB.Query(queryTables)
			if err != nil {
				t.Errorf("received error querying rows: %v", err)
				t.FailNow()
			}
			defer rows.Close()

			got := 0
			for rows.Next() {
				var a, b string
				if err := rows.Scan(&a, &b); err != nil {
					t.Errorf("error scanning row: %v", err)
					continue
				}
				fmt.Println(a, b)
				got++
			}
			if err := rows.Err(); err != nil {
				t.Errorf("rows error: %v", err)
			}
			if got > 0 {
				t.Errorf("got %d want 0 rows", got)
			}
		})
		t.Run("idempotency", func(t *testing.T) {
			_, err := bookshelf.MigrateDown(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("first migrate down failed: %v", err)
			}

			_, err = bookshelf.MigrateUp(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("first migrate up failed: %v", err)
			}

			_, err = bookshelf.MigrateUp(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("second migrate up failed: %v", err)
			}

			_, err = bookshelf.MigrateDown(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("second migrate down failed: %v", err)
			}

			_, err = bookshelf.MigrateDown(dummyWriter, store, "migrations", -1)
			if err != nil {
				t.Errorf("third migrate down failed: %v", err)
			}
		})
	})

	t.Run("CreateBook", func(t *testing.T) {
		t.Run("can create a book", func(t *testing.T) {
			testutils.ResetStore(store)

			var book bookshelf.Book
			err := store.CreateBook(&book, "test-title", "test-author")
			if err != nil {
				t.Errorf("received error on CreateBook: %v", err)
			}
			if book.ID == 0 {
				t.Error("invalid ID received")
			}
		})

		t.Run("cannot create a duplicate title-author", func(t *testing.T) {
			testutils.ResetStore(store)

			var b1, b2 bookshelf.Book
			err := store.CreateBook(&b1, "test-title", "test-author")
			if err != nil {
				t.Errorf("received error on CreateBook: %v", err)
			}

			err = store.CreateBook(&b2, "test-title", "test-author")
			if err == nil {
				t.Error("wanted an error but didn't get one")
			}
		})
	})
}

```

Our tests pass.

```sh
~$ go test ./bookshelf
ok      github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  4.644s
```

Before we move on with the rest of the `CRUD` tests (the `RUD` part), we need a few more tests for `CreateBook`.

