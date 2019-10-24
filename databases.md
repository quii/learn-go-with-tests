# Databases

<!-- TODO: add links to start:v1, end:v2 -->

Oftentimes when creating software, it's necessary to save (or, more precisely, _persist_) some application state.

As an example, when you log into your online banking system, the system has to:

1. Check that it's really you accessing the system (this is called _authentication_, and is beyond the scope of this chapter).
2. Retrieve some information from _somewhere_ and show it to the user (you).

Information that is stored and meant to be long-lived is said to be [_persisted_](<https://en.wikipedia.org/wiki/Persistence_(computer_science)>), usually on a medium that can reliably reproduce the data stored.

Some storage systems, like the filesystem, can be effective for one-off or small amounts of storage, but they fall short for a larger application, for a number of reasons.

This is why most software applications, large and small, opt for storage systems that can provide:

-   Reliability: The data you want is there when you need it.
-   Concurrency: Imagine thousands of users accessing simultaneously.
-   Consistent: You expect the same inputs to produce the same results.
-   Durable: Data should remain there even in case of a system failure (power outage or system crash).

NOTE: The above bullet points are a rewording of the [_ACID principles_](https://en.wikipedia.org/wiki/ACID), it's a set of properties often expressed and used in database design.

_Databases_ are storage mediums that can provide these properties, and much much more.

Also note that, in general, there are two large branches of database types, [SQL](https://en.wikipedia.org/wiki/SQL) and [NOSQL](https://en.wikipedia.org/wiki/NoSQL). In this chapter we will be focusing on SQL databases, using the [database/sql](https://golang.org/pkg/database/sql) package and the `postgres` driver [pq](_ 'github.com/lib/pq').

There is a fair bit of CLI usage in this chapter (mainly setting up the database). For the sake of simplicity we will assume that you are running `ubuntu` on your machine, with `bash` installed. In the near future, look into the appendix for installation on other systems.

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

You can view users and databases with the commands `\du` and `\l` respectively, inside the `psql` shell.

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

We will initially write a (simple) tool to handle our (also simple) migrations, then we will be creating a library that will allow us to perform book-related `CRUD` operations on a `PostgreSQL` database.

## Write the test first

```go
// migrate_test.go
package main

import (
	"testing"
)

func TestMigrateUp(t *testing.T) {
	store, removeStore := NewPostgreSQLStore()
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

Below is some boilerplate that will assist in creating the database connection. We need an `*sql.DB` instance that we can pass on to our tests first. As we learned in the [Dependency Injection](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/dependency-injection) chapter, we should make a helper method so we can acquire an instance from anywhere in our application, while passing in the dependencies, in our case, the database connection string.

Before we continue down the happy-path, we need to make sure our `MigrateUp` function works as expected (unit-test) before we attempt to call it on the real database (integration-test).

The code directly below defines an interface, which will allow us to mock the database functionality, as well as a `NewPostgreSQLStore` function that will simplify its creation.

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

type PostgreSQLStore struct {
	db *sql.DB
}

const (
	removeTimeout = 10 * time.Second
)

func NewPostgreSQLStore() (*PostgreSQLStore, func()) {
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

	return &PostgreSQLStore{db: db}, remove
}

func (s *PostgreSQLStore) ApplyMigration(name, stmt string) error {
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

The `ApplyMigration` method is merely a wrapper around the `sql.DB.Exec` method, this allows us to abstract it when writing our unit tests, testing that the `MigrateUp`, function does what it's intended

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

## 1 - Existence of `dir`

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

## Write the test first

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

## Write the minimal amount of code for the test to run and check the failing test output

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

## Write enough code to make it pass

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

## 2 - Get files inside `dir`

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

## Try to run the test

```sh
# github.com/quii/learn-go-with-tests/databases/v2 [github.com/quii/learn-go-with-tests/databases/v2.test]
.\migrate_test.go:81:17: too many arguments in call to migrate
        have (*SpyStore, string, number, string)
        want (Storer, string, int)
FAIL    github.com/quii/learn-go-with-tests/databases/v2 [build failed]
```

## Write the minimal amount of code for the test to run and check the failing test output

The compiler is complaining, because migrate does not yet accept a direction. Let's `DRY` things a little preemtively this time. Add the following to `bookshelf-store.go`.

```go
// bookshelf-store.go
...
const (
	UP = "up"
	DOWN = "down
)
...
```

Now you can access the directions by a very explicit `UP` or `DOWN`.

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

## 2.1 Ordered migrations

Add the following assertions to `migrate-test.go`

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

	migrations, _ := migrate(store, tmpdir, -1, UP)
	AssertOrderAscending(t, store, migrations)
})

t.Run("down migrations should be ordered descending", func(t *testing.T){
	store := NewSpyStore()
	tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
	defer cleanup()

	migrations, _ := migrate(store, tmpdir, -1, DOWN)
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

## Write enough code to make it pass

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
	case DOWN:
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

On to point no. 4.

## 4. Run `N` migrations in either direction

4. It needs to run migrations `1` through `num`, or all of them if `num == -1`.

For `up` migrations, `1` through `num` makes sense: you will run the `num` first migrations in the dir. But what about `down` migrations? What does `num` represent in this case?

-   The migration number that it will go down to?
-   or the number of files to process before stopping?
-   What if the number of files change before applying `down` migrations?

These are all important questions, but since we are adhering to best practices, our migration files will be _idempotent_, so it does not matter if they are run repeatedly.

`num` should be the number of files to process on `down` migrations, reporting appropriately the state of the database when done.

## Write the test first

Our `CreateTempDir` function creates 3 `up` files, and 3 `down` files if the `empty` boolean param is `false`. Since our `SpyStore` is tracking how many calls each migration receives, we can check that only the first `num` have been called directly, so let's do that.

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

		migrations, _ := migrate(store, tmpdir, 2, UP)
		AssertSliceCalls(t, store, migrations, []int{1, 1, 0})
	})

	t.Run("runs as many migrations as the num param, down", func(t *testing.T) {
		store := NewSpyStore()
		tmpdir, _, cleanup := CreateTempDir(t, "test-migrations", false)
		defer cleanup()

		// `migrations` slice is reversed, so desired order is still (1,1,0)
		migrations, _ := migrate(store, tmpdir, 2, DOWN)
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

## Write enough code to make it pass

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

Our tests now pass.

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

		migrations, _ := migrate(store, tmpdir, -1, UP)
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

## 5.1 Report success

There's a lot of nuance to point number 5, so let's break it into simpler pieces:

1. Lastly, it needs to report on the success of each migration run, if a migration fails, the entire process should be halted.

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
		migrations, _ := migrate(gotBuf, store, tmpdir, -1, UP)
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
migrate(dummyWriter, store, "i-do-not-exist", -1, UP)
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

## 5.2 Reporting failure

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
		_, err = migrate(gotBuf, store, tmpdir, -1, UP)
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
	return migrate(out, os.Stdout, store, dir, num, UP)
}

func MigrateDown(out io.Writer, store Storer, dir string, num int) ([]string, error) {
	return migrate(out, os.Stdout, store, dir, num, DOWN)
}
```

## Write the test first

We're going to rely on the failure of `sql.DB`'s [`Exec`](https://golang.org/pkg/database/sql/#DB.Exec) method, which returns an error if the operation could not complete. We can also use some `PostgreSQL` utilities to list the tables from within our application, and examine the output to verify our tables are there.

Recall our migration written earlier to `migrations/0001_create_books_table.up.sql`. This time we will add some comments to explain what each line does.

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

    A database index is, in layman terms, a trade-off that improves the retrieval of information (if done right) by giving a little more every time data is added.

    Here are more formal definitions if the subject interests you: [Wikipedia](https://en.wikipedia.org/wiki/Database_index), [Use The Index, Luke](https://use-the-index-luke.com/sql/anatomy).

    In this case, the index is redundant, as `PostgreSQL` creates an index on the `PRIMARY KEY` of a table by default.

-   **SQL Language**

    The `SQL` language is part of an ISO standard, and most database engines comform to it partially. This means that code written for one `RDBMS` (say, `PostgreSQL`), will cannot be interpreted as-is by a different one (say, `SQLite3`). There are a lot of similarities, however, and the changes required are often small.

    Keep in mind that the `SQL` you're seeing here is very `PostgreSQL` specific, and some, if not all of it, may not be executable in a different engine.

The test! Here is our first integration test.

```go
//integration_test.go
package main

import "testing"

func TestMigrations(t *testing.T) {
	store, removeStore := NewPostgreSQLStore()
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

Change the `connStr` constant inside `NewPostgreSQLStore` to include the query parameter `sslmode=disable`.

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

You probably noticed that our test are much slower now. Such is the nature of integration tests: testing between services requires more computing power, and has to account for things like latency, message queues and other nuisances that add to the test time.

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
	store, removeStore := NewPostgreSQLStore()
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

Let's think for a minute about the direction we want to take. Just like we did with the `Migrate` function, we want to create package level functions whose behavior we want to control and validate, that in turn call the much simpler method implemented by the `Storer` interface. If you analyze the `ApplyMigration` method of both the `PostgreSQLStore` and the `SpyStore`, you'll notice that they are "dumb", in the sense that they only call their underlying storage engines. The _behavior_ that we wanted was enforced via the `Migrate` function.

This keeps our package _extendible_. With this structure, if a consumer of our `package`'s API wanted to use a different storage method, say, a `NoSQL` database, `AWS S3` file storage or simply a different database (like `MySQL`), they could do so by creating their storage object and having it implement our `Storer` interface, then they can simply plug it into our package level functions (like `Migrate`) and trust they will work (as the _behavior_ is still the same).

The `go` standard library is full useful interfaces like this. A couple of excellent examples are the `encoding/json` package [`Marshaler`](https://golang.org/pkg/encoding/json/#Marshaler) and [`Unmarshaler`](https://golang.org/pkg/encoding/json/#Unmarshaler) interfaces. `encoding/json` has ensured these interfaces are implemented by all standard types, but if you want a certain behavior, simply implement these interfaces in your custom type and encoding will work.

With this mindset, we should aim to write functions to `Create`, `Retrieve`, `Update`, `Delete` and `List` books, with those same names. These verbs will perform all kinds of validation, and, if they pass, then call the `Storer` methods by the same name.

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

While we're at it, we should test that the title or author aren't empty (remember the `NOT NULL` in our migrations?)

Create a new file `crud_test.go` and add the following code:

```go
// bookshelf/crud_test.go
package bookshelf_test

import (
	"testing"

	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf"
	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf/testutils"
)

func TestCreate(t *testing.T) {

	t.Run("creates accurately", func(t *testing.T) {
		store := testutils.NewSpyStore()

		title, author := "Moby Dick", "Herman Melville"

		book, err := bookshelf.Create(store, title, author)
		testutils.AssertNoError(t, err)
		if book.ID == 0 {
			t.Error("book returned without an ID")
		}
		if book.Title != title {
			t.Errorf("got %q want %q for title", book.Title, title)
		}
		if book.Author != author {
			t.Errorf("got %q want %q for author", book.Author, author)
		}
	})
	t.Run("error on empty title", func(t *testing.T) {
		store := testutils.NewSpyStore()

		title, author := "", "Herman Melville"

		_, err := bookshelf.Create(store, title, author)
		testutils.AssertError(t, err, bookshelf.ErrEmptyTitleField)
	})
	t.Run("error on empty author", func(t *testing.T) {
		store := testutils.NewSpyStore()

		title, author := "Moby Dick", ""

		_, err := bookshelf.Create(store, title, author)
		testutils.AssertError(t, err, bookshelf.ErrEmptyAuthorField)
	})
}
```

## Try to run the test

It fails, as expected, because we haven't written anything yet.

```sh
~$ go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf_test [github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf.test]
bookshelf/crud_test.go:17:16: undefined: "github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf".Create
bookshelf/crud_test.go:34:13: undefined: "github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf".Create
bookshelf/crud_test.go:35:33: undefined: "github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf".ErrEmptyTitleField
bookshelf/crud_test.go:42:13: undefined: "github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf".Create
bookshelf/crud_test.go:43:33: undefined: "github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf".ErrEmptyAuthorField
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf [build failed]
FAIL
```

## Write the minimal amount of code for the test to run and check the failing test output

Add the following into `bookshelf-store.go`

```go
// bookshelf/bookshelf-store.go
...
var (
	ErrEmptyTitleField = errors.New("empty title field")
	ErrEmptyAuthorField = errors.New("empty author field")
)
...
func Create(store Storer, title, author string) (*Book, error) {
	var book Book
	return &book, nil
}
```

Run the tests again. We get our expected failures.

```sh
~$ go test ./bookshelf
--- FAIL: TestCreate (0.00s)
    --- FAIL: TestCreate/creates_accurately (0.00s)
        crud_test.go:20: book returned without an ID
        crud_test.go:23: got "" want "Moby Dick" for title
        crud_test.go:26: got "" want "Herman Melville" for author
    --- FAIL: TestCreate/error_on_empty_title (0.00s)
        crud_test.go:35: wanted an error but didn't get one
        crud_test.go:35: got <nil> want empty title field
    --- FAIL: TestCreate/error_on_empty_author (0.00s)
        crud_test.go:43: wanted an error but didn't get one
        crud_test.go:43: got <nil> want empty author field
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  0.079s
FAIL
```

## Write enough code to make it pass

```go
// bookshelf/bookshelf-store.go
func Create(store Storer, title, author string) (*Book, error) {
	if title == "" {
		return nil, ErrEmptyTitleField
	}
	if author == "" {
		return nil, ErrEmptyAuthorField
	}
	var book Book
	err := store.Create(&book, title, author)
	if err != nil {
		return nil, err
	}
	return &book, err
}
```

## Try to run the test

Once more, expected failure

```sh
~$ go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf
bookshelf/bookshelf-store.go:177:14: store.Create undefined (type Storer has no field or method Create)
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf [build failed]
FAIL
```

## Write enough code to make it pass

Add the `Create` method to the `Storer` interface, with the same signature.

```go
// bookshelf/bookshelf-store.go
...
type Storer interface {
	ApplyMigration(name, stmt string) error
	Create(book *Book, title string, author string) error
}
...
```

And in `testutils/store.go`. We'll use a helper method to acquire the latest `id`.

```go
// testutils/store.go
package testutils
import (
	...
	"github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf"
)
...
type SpyStore struct {
	Migrations map[string]migration
	Books      []*bookshelf.Book
}
...
func (s *SpyStore) Create(book *bookshelf.Book, title, author string) error {
	book.ID = newID(s)
	book.Title = title
	book.Author = author
	s.Books = append(s.Books, book)
	return nil
}

func NewSpyStore() *SpyStore {
	return &SpyStore{
		Migrations: map[string]migration{},
		Books:      []*bookshelf.Book{},
	}
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
# github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf_test [github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf.test]
bookshelf/integration_test.go:27:32: cannot use store (type *bookshelf.PostgreSQLStore) as type bookshelf.Storer in argument to bookshelf.MigrateUp:
        *bookshelf.PostgreSQLStore does not implement bookshelf.Storer (missing Create method)
...
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf [build failed]
FAIL
```

Seems like we broke our integration tests. Our `PostgreSQLStore` does not implement the `Create` method. Let's do that now.

```go
// bookshelf/bookshelf-store.go
// Create inserts a new book into the postgres store.
func (s *PostgreSQLStore) Create(book *Book, title string, author string) error {
	stmt := "INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id, title, author;"
	row := s.DB.QueryRow(stmt, title, author)
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return err
	}
	return nil
}
```

Our tests now pass.

```sh
~$ go test ./bookshelf
ok      github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf      0.074s
```

The create method of the `PostgreSQLStore` is not tested yet, but we will get to that. Let's finish testing the behavior of the package level `Create` function first.

Notice that the `Create` method does not have an `ID` field provided, this is because `primary keys` are usually autoincremented and provided by the database.

This is not written in stone, however: just happens that when we created the database table, we designated the `id` field as `SERIAL`, and `PostgreSQL` handles the auto-incrementing for us. But we could have designated a `PRIMARY KEY` of whichever type we would've wanted. For example, had we assigned the `title` as the primary key, the `id` field would have not been necessary. Or, if we designated `id` as `PRIMARY KEY`, but as type `INT` instead of `SERIAL`, our application would have had to find the latest `id` and increment it.

As a side note, `PostgreSQL`'s behavior regarding the `SERIAL` type is to keep autoincrementing it no matter what, even if previous `id`s are gone. For example, if you create `books` with `id`s 1, 2, 3 and 4, and delete the `book` with `id` 2, the next `id` created will be `5`. You could reassign it should you want to, but this leads to confusion and it's generally regarded as a bad practice.

## But wait, what if we create the same book twice?

Our tests pass, but what about duplicated data? We don't need two, three or twenty entries of the same book.

As it turns out, we made a mistake when creating the tables in `migrations/0001_create_books_table.up.sql`. While we designated the `author` and the `title` columns to be required (`NOT NULL`, we should test against this too!), we did not designate the `title` as `UNIQUE`. We should be careful, as there could be different `author`s with books that share a `title`; this may seem like an edge case, but edge cases is one of many reasons why we test!

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

We should test our migrations before we move on, our current integration tests should tell us if our `SQL` is correct.

```sh
~$ go test ./bookshelf
--- FAIL: TestMigrations (0.11s)
    --- FAIL: TestMigrations/idempotency (0.06s)
        integration_test.go:98: second migrate up failed: pq: relation "books_unique_author_title" already exists
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  0.117s
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

This behavior is enforced by the database, at the database level. We need to give our potential users a way to prevent duplicated data as well.

We could create an `Exists` function, but this would be terribly inefficient, as it would query the database, to then tell the user `said book exists`, like a mathematician.

We should instead write the `Retrieve` function instead. We can use this to check if the book exists, and return it (or rather, populate the `*Book*` object with it).

## Write the test first

We have a choice to make for the retrieve function, we could

-   a) Create a function with signature`Retrieve(store, param interface{}) (*Book, error)` , and do a type switch inside: `int` leads to type search, and `string` leads to title search
-   b) Create two functions, with signatures `ByID(store Storer, id int64) (*Book, error)` and `ByTitleAuthor(store Storer, author, title string) (*Book, error)`.

I believe option `b)` is the better one. First, our code is more explicit. Second, there's no `interface{}` dances inside our tests. Let's do that.

Modify `NewSpyStore` to include an `initialBooks` parameter, so we can populate the `SpyStore` on creation.

```go
// bookshelf/testutils/store.go
...
func NewSpyStore(books []*bookshelf.Book) *SpyStore {
	return &SpyStore{
		Migrations: map[string]migration{},
		Books:     books,
	}
}
...
```

Add a new assertion in `testutils/assertions.go`

```go
// bookshelf/testutils/assertions.go
package testutils
import (
	...
	"github.com/djangulo/learn-go-with-tests/database/v5/bookshelf"
)
...
func AssertBooksEqual(t *testing.T, got, want *bookshelf.Book) {
	t.Helper()
	if got == nil || got.ID == 0 {
		t.Errorf("nil or invalid ID: %v", got)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

Then create a `dummyBooks` variable and add the following test to `crud_test.go`. You can change `TestCreate` to use the `AssertBooksEqual` helper too.

You will have to modify all occurrences of `NewSpyStore` in `migrate_test.go` and in `crud_test.go`.

```go
// bookshelf/crud_test.go
...
var (
	dummyBooks = make([]*bookshelf.Book, 0)
	testBooks = []*bookshelf.Book{
		&bookshelf.Book{ID: 10, Author: "W. Shakespeare", Title: "The Tragedie of Hamlet"},
		&bookshelf.Book{ID: 22, Author: "W. Shakespeare", Title: "Romeo & Juliet"},
		&bookshelf.Book{ID: 24, Author: "Ernest Hemingway", Title: "The Old Man and The Sea"},
	}
)
...
func TestByID(t *testing.T) {

	t.Run("ByID success", func(t *testing.T) {
		store := testutils.NewSpyStore(testBooks)

		book, err := bookshelf.ByID(store, testBooks[0].ID)
		testutils.AssertNoError(t, err)
		testutils.AssertBooksEqual(t, book, testBooks[0])
	})

	for _, test := range []struct {
		name string
		in   int64
		want error
	}{
		{"ByID not found", int64(42), bookshelf.ErrBookDoesNotExist},
		{"ByID zero value", int64(0), bookshelf.ErrZeroValueID},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := testutils.NewSpyStore(testBooks)
			_, err := bookshelf.ByID(store, test.in)
			testutils.AssertError(t, err, test.want)
		})
	}

}

func TestByAuthorTitle(t *testing.T) {

	t.Run("ByTitleAuthor success", func(t *testing.T) {
		store := testutils.NewSpyStore(testBooks)

		book, err := bookshelf.ByTitleAuthor(store, testBooks[0].Title, testBooks[0].Author)
		testutils.AssertNoError(t, err)
		testutils.AssertBooksEqual(t, book, testBooks[0])
	})

	for _, test := range []struct {
		name, title, author string
		want                error
	}{
		{"ByTitleAuthor failure empty title", "", "Herman Melville", bookshelf.ErrEmptyTitleField},
		{"ByTitleAuthor failure empty author", "Moby Dick", "", bookshelf.ErrEmptyAuthorField},
		{"ByTitleAuthor failure not found", "Moby Dick", "Herman Melville", bookshelf.ErrBookDoesNotExist},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := testutils.NewSpyStore(testBooks)
			_, err := bookshelf.ByTitleAuthor(store, test.title, test.author)
			testutils.AssertError(t, err, test.want)
		})
	}

}
```

## Try to run the test

Expected failures:

```sh
~$ go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf_test [github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf.test]
bookshelf/crud_test.go:52:16: undefined: bookshelf.ByID
bookshelf/crud_test.go:62:33: undefined: bookshelf.ErrBookDoesNotExist
bookshelf/crud_test.go:63:33: undefined: bookshelf.ErrZeroValueID
bookshelf/crud_test.go:67:14: undefined: bookshelf.ByID
bookshelf/crud_test.go:79:16: undefined: bookshelf.ByTitleAuthor
bookshelf/crud_test.go:90:62: undefined: bookshelf.ErrBookDoesNotExist
bookshelf/crud_test.go:94:14: undefined: bookshelf.ByTitleAuthor
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf [build failed]
FAIL
```

## Write the minimal amount of code for the test to run and check the failing test output

```go
// bookshelf/bookshelf-store.go
...
var (
	...
	ErrZeroValueID = errors.New("zero value ID")
	ErrBookDoesNotExist = errors.New("book does not exist")
)
...
func ByID(store Storer, id int64) (*Book, error) {
	var book Book
	return &book, nil
}
func ByTitleAuthor(store Storer, title, author string) (*Book, error) {
	var book Book
	return &book, nil
}
```

## Try to run the tests

```sh
~$ go test ./bookshelf
--- FAIL: TestByID (0.00s)
    --- FAIL: TestByID/ByID_success (0.00s)
        crud_test.go:46: nil or invalid ID: <nil>
        crud_test.go:46: got <nil> want &{10 The Tragedie of Hamlet W. Shakespeare}
    --- FAIL: TestByID/ByID_not_found (0.00s)
        crud_test.go:60: wanted an error but didn't get one
        crud_test.go:60: got <nil> want book does not exist
    --- FAIL: TestByID/ByID_zero_value (0.00s)
        crud_test.go:60: wanted an error but didn't get one
        crud_test.go:60: got <nil> want zero value ID
--- FAIL: TestByAuthorTitle (0.00s)
    --- FAIL: TestByAuthorTitle/ByAuthorTitle_success (0.00s)
        crud_test.go:73: nil or invalid ID: <nil>
        crud_test.go:73: got <nil> want &{10 The Tragedie of Hamlet W. Shakespeare}
    --- FAIL: TestByAuthorTitle/ByAuthorTitle_failure_empty_title (0.00s)
        crud_test.go:87: wanted an error but didn't get one
        crud_test.go:87: got <nil> want empty title field
    --- FAIL: TestByAuthorTitle/ByAuthorTitle_failure_empty_author (0.00s)
        crud_test.go:87: wanted an error but didn't get one
        crud_test.go:87: got <nil> want empty author field
    --- FAIL: TestByAuthorTitle/ByAuthorTitle_failure_not_found (0.00s)
        crud_test.go:87: wanted an error but didn't get one
        crud_test.go:87: got <nil> want book does not exist
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  0.122s
FAIL
```

We got our work cut out for us, let's get to it!

## Write enough code to make it pass

```go
// bookshelf/bookshelf-store.go
...
func ByID(store Storer, id int64) (*Book, error) {
	if id == 0 {
		return nil, ErrZeroValueID
	}
	var book Book
	err := store.ByID(&book, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBookDoesNotExist
		}
		return nil, err
	}
	return &book, nil
}

func ByTitleAuthor(store Storer, title, author string) (*Book, error) {
	if title == "" {
		return nil, ErrEmptyTitleField
	}
	if author == "" {
		return nil, ErrEmptyAuthorField
	}
	var book Book
	err := store.ByTitleAuthor(&book, title, author)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBookDoesNotExist
		}
		return nil, err
	}
	return &book, nil
}
```

You know what's coming by now:

```sh
~$ go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf
bookshelf/bookshelf-store.go:195:13: store.ByID undefined (type Storer has no field or method ByID)
bookshelf/bookshelf-store.go:210:13: store.ByTitleAuthor undefined (type Storer has no field or method ByTitleAuthor)
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf [build failed]
FAIL
```

Add the `Storer` method signatures.

```go
// bookshelf/bookshelf-store.go
...
type Storer interface {
	ApplyMigration(name, stmt string) error
	Create(book *Book, title string, author string) error
	ByID(book *Book, id int64) error
	ByTitleAuthor(book *Book, title string, author string) error
}
...
```

## Try to run the tests

It'll be a looong list of the same error. `integration_test.go` and `crud_test.go` are not happy with our `Storer` changes.

```sh
~$ go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf_test [github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf.test]
bookshelf/crud_test.go:17:32: cannot use store (type *testutils.SpyStore) as type "github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf".Storer in argument to "github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf".Create:
        *testutils.SpyStore does not implement "github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf".Storer (missing ByID method)
...
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf [build failed]
FAIL
```

## Write enough code to make it pass

Add the `ByID` and the `ByTitleAuthor` to the `SpyStore`:

```go
// bookshelf/testutils/store.go
...
func (s *SpyStore) ByID(book *bookshelf.Book, id int64) error {
	for _, b := range s.Books {
		if b.ID == id {
			*book = *b
			return nil
		}
	}
	return bookshelf.ErrBookDoesNotExist
}

func (s *SpyStore) ByTitleAuthor(book *bookshelf.Book, title, author string) error {
	title, author = strings.ToLower(title), strings.ToLower(author)
	for _, b := range s.Books {
		if strings.ToLower(b.Title) == title && strings.ToLower(b.Author) == author {
			*book = *b
			return nil
		}
	}
	return bookshelf.ErrBookDoesNotExist
}
```

And to the `PostgreSQLStore`:

```go
// bookshelf/bookshelf-store.go
func (s *PostgreSQLStore) ByID(book *Book, id int64) error {
	stmt := "SELECT id, title, author FROM books WHERE id = $1 LIMIT 1;"
	row := s.DB.QueryRow(stmt, id)
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgreSQLStore) ByTitleAuthor(book *Book, title, author string) error {
	stmt := "SELECT id, title, author FROM books WHERE title = $1 AND author = $2 LIMIT 1;"
	row := s.DB.QueryRow(stmt, title, author)
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return err
	}
	return nil
}
```

And we're back to green!

```sh
~$ go test ./bookshelf
ok      github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf      0.144s
```

Now that we've written and tested our retrieve methods, we can implement them inside the `Create` to validate existence.

This woul render the `Create` method a tad inefficient: it would entail a `read+write` operation as opposed to just a `write`. Because of this, we want to leave it optional to the consumers of our API. Let's instead write a `GetOrCreate` function that does exactly that, and leave `Create` unchanged.

## Write the test first

```go
// bookshelf/crud_test.go
func TestGetOrCreate(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		for _, test := range []struct {
			name, title, author string
			want                *bookshelf.Book
		}{
			{"GetOrCreate retrieves from store", testBooks[1].Title, testBooks[1].Author, testBooks[1]},
			{
				"GetOrCreate insert into store",
				"Moby Dick",
				"Herman Melville",
				&bookshelf.Book{testBooks[len(testBooks)-1].ID + 1, "Moby Dick", "Herman Melville"},
			},
		} {
			t.Run(test.name, func(t *testing.T) {
				store := testutils.NewSpyStore(testBooks)
				book, err := bookshelf.GetOrCreate(store, test.title, test.author)
				testutils.AssertNoError(t, err)
				testutils.AssertBooksEqual(t, book, test.want)
			})
		}
	})
	t.Run("failure", func(t *testing.T) {
		for _, test := range []struct {
			name, title, author string
			want                error
		}{
			{"GetOrCreate failure empty title", "", "Herman Melville", bookshelf.ErrEmptyTitleField},
			{"GetOrCreate failure empty author", "Moby Dick", "", bookshelf.ErrEmptyAuthorField},
		} {
			t.Run(test.name, func(t *testing.T) {
				store := testutils.NewSpyStore(testBooks)
				_, err := bookshelf.GetOrCreate(store, test.title, test.author)
				testutils.AssertError(t, err, test.want)
			})
		}
	})
}
```

## Try to run the test

```sh
~$ go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf_test [github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf.test]
bookshelf/crud_test.go:126:18: undefined: bookshelf.GetOrCreate
bookshelf/crud_test.go:142:15: undefined: bookshelf.GetOrCreate
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf [build failed]
FAIL
djangulo@sif:~/go/src/github.com/djangulo/l
```

## Write the minimal amount of code for the test to run and check the failing test output

```go
// bookshelf/bookshelf-store.go
...
func GetOrCreate(store Storer, title, author string) (*Book, error){
	var book Book
	return &book, nil
}
```

## Try to run the test

Our expected failures:

```sh
~$ go test ./bookshelf
--- FAIL: TestGetOrCreate (0.00s)
    --- FAIL: TestGetOrCreate/success (0.00s)
        --- FAIL: TestGetOrCreate/success/GetOrCreate_retrieves_from_store (0.00s)
            crud_test.go:112: nil or invalid ID: &{0  }
            crud_test.go:112: got &{0  } want &{22 Romeo & Juliet W. Shakespeare}
        --- FAIL: TestGetOrCreate/success/GetOrCreate_insert_into_store (0.00s)
            crud_test.go:112: nil or invalid ID: &{0  }
            crud_test.go:112: got &{0  } want &{23 Moby Dick Herman Melville}
    --- FAIL: TestGetOrCreate/failure (0.00s)
        --- FAIL: TestGetOrCreate/failure/GetOrCreate_failure_empty_title (0.00s)
            crud_test.go:127: wanted an error but didn't get one
            crud_test.go:127: got <nil> want empty title field
        --- FAIL: TestGetOrCreate/failure/GetOrCreate_failure_empty_author (0.00s)
            crud_test.go:127: wanted an error but didn't get one
            crud_test.go:127: got <nil> want empty author field
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  0.120s
FAIL
```

## Write enough code to make it pass

Thankfully, all the hard work is already done (and tested). Just use the `ByTitleAuthor` function to validate existence, and if the error is `ErrBookDoesNotExist`, call `Create`.

```go
// bookshelf/bookshelf-store.go
...
func GetOrCreate(store Storer, title, author string) (*Book, error) {
	book, _ := ByTitleAuthor(store, title, author)
	if book != nil {
		return book, nil
	}
	book, err := Create(store, title, author)
	if err != nil {
		return nil, err
	}
	return book, nil

```

## Try to run the tests

And our tests now pass.

```sh
~$ go test ./bookshelf
ok      github.com/djangulo/learn-go-with-tests/databases/v5/bookshelf  0.153s
```

## The `C` and the `R`

<!-- TODO: start: v5, end: v6 -->

We've been to focused on the unit tests. We need to create integration tests for our `Create`, `ByID`, `ByTitleAuthor` and `GetOrCreate` functions.

As usual, let's start with the test.

## Write the test first

```go
// bookshelf/integration_test.go
...
func TestCreateBook(t *testing.T) {
	store, removeStore := bookshelf.NewPostgreSQLStore()
	defer removeStore()

	t.Run("can create a book", func(t *testing.T) {
		book, err := bookshelf.Create(store, "test-title", "test-author")
		if err != nil {
			t.Errorf("received error on CreateBook: %v", err)
		}
		testutils.AssertBooksEqual(t, book, &bookshelf.Book{1, "test-title", "test-author"})
	})

	t.Run("cannot create a duplicate title-author", func(t *testing.T) {
		_, err := bookshelf.Create(store, "test-title", "test-author")
		if err != nil {
			t.Errorf("received error on CreateBook: %v", err)
		}
		_, err = bookshelf.Create(store, "test-title", "test-author")
		if err == nil {
			t.Error("wanted an error but didn't get one")
		}

	})
}
```

## Try to run the test

```sh
~$ go test ./bookshelf
--- FAIL: TestCreateBook (0.00s)
    --- FAIL: TestCreateBook/can_create_a_book (0.00s)
        integration_test.go:121: received error on CreateBook: pq: relation "books" does not exist
        integration_test.go:123: nil or invalid ID: <nil>
        integration_test.go:123: got <nil> want &{1 test-title test-author}
    --- FAIL: TestCreateBook/cannot_create_a_duplicate_title-author (0.00s)
        integration_test.go:129: received error on CreateBook: pq: relation "books" does not exist
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf  0.155s
FAIL
```

Huh, "relation `books` does not exist". By now, you probably know this means the table `books` does not exist in the database.

Well, we did not migrate before we tested, so that makes sense.

Add a call to `MigrateUp`, and check the error, after acquiring the `store`.

```go
// bookshelf/integration_test.go
func TestCreateBook(t *testing.T) {
	store, removeStore := bookshelf.NewPostgreSQLStore()
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
~$ go test  ./bookshelf
--- FAIL: TestCreateBook (0.05s)
    --- FAIL: TestCreateBook/cannot_create_a_duplicate_title-author (0.00s)
        integration_test.go:135: received error on CreateBook: pq: duplicate key value violates unique constraint "books_unique_author_title"
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf  0.190s
FAIL
```

Is it the error we expect? Yes it is. But the tests should be _passing_, so what is going on?

The test raising the error is the first `CreateBook` call inside `cannot create a duplicate title-author`. This is not what we planned!

As it turns out, there is only 1 database connection string in our application; the one that lives inside `NewPostgreSQLStore`. This means that all along we've been migrating `up` and `down`, inserting and deleting into a single database! This is a very risky practice, as our tests may modify or delete sensitive data once we're running in production.

So what do we do?

## The test database

So far, we have been operating in the `bookshelf_db` database, that we created at the start of the chapter. We need a secondary database that we can test to our heart's content.

Our options are:

1. Write some `go` code that creates a test database on the fly. Runs the tests and drops it once it's done.
2. Create a test database outside our application (using `psql`), and hardcode the address in our application.

Both approaches have their downsides:

1. The first approach requires more `go` code to write, and depends on the privileges that the DB user (`bookshelf_user`, in our case) has. When created the database, we gave our user `bookshelf_user` the capacity to create databases with `CREATEDB`.
2. The second approach is simpler to implement, but then our tests depend on the existence of said test database. It also implies that we need to track 2 connections inside our application, as opposed to just

We will opt for the first approach, and take advantage that our `MigrateDown` function cleans the database tables, due to the `CASCADE` statement at the end of the `down` migrations.

## Refactor

Now we need to refactor our code. Stop and think about what's going on inside `NewPostgreSQLStore`.

```go
// bookshelf/bookshelf-store.go
func NewPostgreSQLStore() (*PostgreSQLStore, func()) {
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

Our `NewPostgreSQLStore` function now looks like this:

```go
// bookshelf/bookshelf-store.go
func NewPostgreSQLStore(conf *DBConf) (*PostgreSQLStore, func()) {

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

	return &PostgreSQLStore{DB: db}, remove
}
```

With this tooling in place, we can create a helper function to instantiate a new database just for our tests.

Insert the functions and variables below inside `bookshelf/testutils/helpers.go`:

```go
// bookshelf/testutils/helpers.go
package testutils
import (
	...
	"database/sql"
	"time"
	"math/rand"

	"github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf"
	_ "github.com/lib/pq"
)
...
type TestDBRegistry struct {
	Databases map[string]*bookshelf.DBConf
	Prefix    string
}

func (t *TestDBRegistry) Add(conf *bookshelf.DBConf) string {
	rand.Seed(time.Now().UnixNano())
	dbname := (*t).Prefix + "_" + randString(20)
	
	(*conf).DBName = dbname
	(*t).Databases[dbname] = conf
	
	return dbname
}

func (t *TestDBRegistry) Remove(dbname string) {
	if _, ok := (*t).Databases[dbname]; ok {
		delete((*t).Databases, dbname)
	}
}

func randString(n int) string {
	b := make([]rune, n)
	for i := 0; i < n; i++ {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

var (
	chars = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
		ActiveTestDBRegistry = &TestDBRegistry{
	Databases: map[string]*bookshelf.DBConf{},
	Prefix:    "bookshelf_test_db",
	}
)

func NewTestPostgreSQLStore() (*bookshelf.PostgreSQLStore, func(), error) {
	main, removeMain := bookshelf.NewPostgreSQLStore(&bookshelf.MainDBConf)
	dbconf := &bookshelf.DBConf{
		User:    bookshelf.MainDBConf.User,
		Pass:    bookshelf.MainDBConf.Pass,
		Host:    bookshelf.MainDBConf.Host,
		Port:    bookshelf.MainDBConf.Port,
		SSLMode: bookshelf.MainDBConf.SSLMode,
	}

	dbname := ActiveTestDBRegistry.Add(dbconf)

	_, err := main.DB.Exec(
		fmt.Sprintf("CREATE DATABASE %s OWNER %s;",
			dbname,
			bookshelf.MainDBConf.User,
		),
	)
	if err != nil {
		return nil, nil, err
	}

	testDB, err := sql.Open("postgres", dbconf.String())
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
					dbname,
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
			_, err := main.DB.Exec(fmt.Sprintf("DROP DATABASE %s;", dbname))
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"error dropping test database %q, retrying in %v: %v\n",
					dbname,
					retryIn,
					err,
				)
				time.Sleep(retryIn)
				continue
			}
			break
		}
		ActiveTestDBRegistry.Remove(dbname)
		removeMain()
	}
	return &bookshelf.PostgreSQLStore{DB: testDB}, remove, nil
}
```

Remember `exponential backoff`? This is where this pattern shines in our codebase. If the `mainDB` tries to drop the test database, but it's being written by a test, the operation will fail. With exponential backoff, it'll give the running operations a bit of time to finish, and then finally drop the test database.

We now can create the test database inside each test function, run all our tests with a predictable state, and drop it once we're done. We can use the fact that `MigrateDown` clears the database to our advantage and clean it after each test.

For the sake of demonstration, we will use this helper function... **a lot**. Each of the integration tests we write will run in its own database. This allows us to run with a predictable state every time. Even better, now that we have disposable databases, we can make put `test` migrations inside our `migrations` directory, ensuring that not only the test databases will be available, but they also will be populated with data.

If all this creating databases is slowing down your computer, feel free to group all integration tests from here on into a single one, using `t.Run` freely to nest different function names.

Add the following test to `migrate_test.go`, to ensure that sub-directories are ignored:

```go
// bookshelf/migrate_test.go
...
	t.Run("ignores subdirectiories", func(t *testing.T) {
		store := testutils.NewSpyStore(dummyBooks)

		tmpdir, _, cleanup := testutils.CreateTempDir(t, "test-migrations", false)
		defer cleanup()
		subdir, err := ioutil.TempDir(tmpdir, "subdirectory")
		if err != nil {
			t.FailNow()
		}
		f, err := ioutil.TempFile(subdir, "subfile.*.up.sql")
		if err != nil {
			t.FailNow()
		}
		if _, ok := store.Migrations[f.Name()]; ok {
			t.Errorf("%q is not supposed to exist in the store", f.Name())
		}
	})
	...
```

## Try to run the test

Seems like we didn't break anything.

```sh
~$ go test ./bookshelf
ok      github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf  5.362s
```

Now create a copy of of the `migrations` directory into itself, name it `test`, and add the following migrations into it:

```sql
-- migrations/test/0003_insert_test_data.up.sql
INSERT INTO books (title, author) VALUES 
('Alice''s Adventures in Wonderland', 'Lewis Carroll'),
('The Ball and The Cross', 'G.K. Chesterton'),
('The Man Who Was Thursday', 'G.K. Chesterton'),
('Moby Dick', 'Herman Melville'),
('Paradise Lost', 'John Milton'),
('The Tragedie of Julius Caesar', 'William Shakespeare'),
('The Tragedie of Hamlet', 'William Shakespeare'),
('The Tragedie of Macbeth', 'William Shakespeare'),
('Romeo and Juliet', 'William Shakespeare') ON CONFLICT DO NOTHING;
```
```sql
-- migrations/test/0003_insert_test_data.down.sql
DELETE  FROM books;
```

Let's include a boolean `migrate` parameter in `NewTestPostgreSQLStore`,  to optionally populate the new test database..

Let's create another test utility to reset the database on a whim.

```go
// bookshelf/testutils/helpers.go
...
func NewTestPostgreSQLStore(migrate bool) (*bookshelf.PostgreSQLStore, func(), error) {
	...
	store := bookshelf.PostgreSQLStore{DB: testDB}
	if migrate {
		bookshelf.MigrateUp(dummyWriter, &store, "migrations/test", -1)
	}

	return &store, remove, nil
}

func ResetStore(store *bookshelf.PostgreSQLStore) error {
	var err error
	_, err = bookshelf.MigrateDown(dummyWriter, store, "migrations/test", -1)
	if err != nil {
		return err
	}

	_, err = bookshelf.MigrateUp(dummyWriter, store, "migrations/test", -1)
	if err != nil {
		return err
	}

	return nil
}
```

Keep in mind that each database created will make our tests slower. This affects us now, but in reality we should be offloading this types of tests to CI. You want to be as thorough as possible if you have full control over the external service.

And finally, our integration tests. We've included the tests for `ByID`, `ByTitleAuthor`, and `GetOrCreate` that were omitted before for brevity.

```go
// bookshelf/integration_test.go
...
import (
	...
	"fmt"
	"os"
	"github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf/testutils"
)
...
func TestMigrateIntegration(t *testing.T) {
	store, removeStore, err := testutils.NewTestPostgreSQLStore(false)
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
}
```

Since `TestMigrateIntegration` takes so much space, let's start using table-driven tests for our successes and failures.

```go
// bookshelf/integration_test.go
func TestCreateIntegration(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		for _, test := range []struct {
			name, title, author string
			want                *bookshelf.Book
		}{
			{"can create", "test-title", "test-author", &bookshelf.Book{10, "test-title", "test-author"}},
		} {
			t.Run(test.name, func(t *testing.T) {
				store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
				if err != nil {
					fmt.Fprintf(os.Stdout, "db creation failed on test %q", test.name)
					t.FailNow()
				}
				defer removeStore()

				book, err := bookshelf.Create(store, test.title, test.author)

				testutils.AssertNoError(t, err)
				testutils.AssertBooksEqual(t, book, test.want)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		t.Run("cannot create a duplicate title-author", func(t *testing.T) {
			store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
			if err != nil {
				fmt.Fprintf(os.Stdout, "db creation failed on test for duplicate title-author")
				t.FailNow()
			}
			defer removeStore()

			_, err = bookshelf.Create(store, "test-title", "test-author")
			if err != nil {
				t.Errorf("received error on CreateBook: %v", err)
			}

			_, err = bookshelf.Create(store, "test-title", "test-author")
			if err == nil {
				t.Error("wanted an error but didn't get one")
			}
		})
	})

}
func TestByIDIntegration(t *testing.T) {
	store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
	if err != nil {
		fmt.Fprintf(os.Stdout, "db creation failed on TestByIDIntegration")
		t.FailNow()
	}
	defer removeStore()

	t.Run("success", func(t *testing.T) {
		got, err := bookshelf.ByID(store, 1)
		testutils.AssertNoError(t, err)
		want := &bookshelf.Book{1, "Alice's Adventures in Wonderland", "Lewis Carroll"}
		testutils.AssertBooksEqual(t, got, want)
	})

	t.Run("failure", func(t *testing.T) {
		for _, test := range []struct {
			name string
			in   int64
			want error
		}{
			{"not found", int64(42), bookshelf.ErrBookDoesNotExist},
			{"zero value", int64(0), bookshelf.ErrZeroValueID},
		} {
			t.Run(test.name, func(t *testing.T) {
				_, err := bookshelf.ByID(store, test.in)
				testutils.AssertError(t, err, test.want)
			})
		}
	})
}

func TestByTitleAuthorIntegration(t *testing.T) {
	store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
	if err != nil {
		fmt.Fprintf(os.Stdout, "db creation failed on TestByIDIntegration")
		t.FailNow()
	}
	defer removeStore()

	t.Run("success", func(t *testing.T) {
		got, err := bookshelf.ByTitleAuthor(store, "Alice's Adventures in Wonderland", "Lewis Carroll")
		want := &bookshelf.Book{1, "Alice's Adventures in Wonderland", "Lewis Carroll"}
		testutils.AssertNoError(t, err)
		testutils.AssertBooksEqual(t, got, want)
	})

	t.Run("failure", func(t *testing.T) {
		for _, test := range []struct {
			name, title, author string
			want                error
		}{
			{"empty title", "", "Herman Melville", bookshelf.ErrEmptyTitleField},
			{"empty author", "Moby Dick", "", bookshelf.ErrEmptyAuthorField},
			{"not found", "The DaVinci Code", "Dan Brown", bookshelf.ErrBookDoesNotExist},
		} {
			t.Run(test.name, func(t *testing.T) {
				_, err := bookshelf.ByTitleAuthor(store, test.title, test.author)
				testutils.AssertError(t, err, test.want)
			})
		}
	})
}

func TestGetOrCreateIntegration(t *testing.T) {
	store, removeStore, err := testutils.NewTestPostgreSQLStore(true)
	if err != nil {
		fmt.Fprintf(os.Stdout, "db creation failed on TestByIDIntegration")
		t.FailNow()
	}
	defer removeStore()

	t.Run("success", func(t *testing.T) {
		for _, test := range []struct {
			name, title, author string
			want                *bookshelf.Book
		}{
			{"retrieves", "Alice's Adventures in Wonderland", "Lewis Carroll", &bookshelf.Book{1, "Alice's Adventures in Wonderland", "Lewis Carroll"}},
			{"inserts", "DaVinci", "Dan Brown", &bookshelf.Book{10, "DaVinci", "Dan Brown"}},
		} {
			t.Run(test.name, func(t *testing.T) {
				book, err := bookshelf.GetOrCreate(store, test.title, test.author)
				testutils.AssertNoError(t, err)
				testutils.AssertBooksEqual(t, book, test.want)
			})
		}
	})
	t.Run("failure", func(t *testing.T) {
		for _, test := range []struct {
			name, title, author string
			want                error
		}{
			{"empty title", "", "Herman Melville", bookshelf.ErrEmptyTitleField},
			{"empty author", "Moby Dick", "", bookshelf.ErrEmptyAuthorField},
		} {
			t.Run(test.name, func(t *testing.T) {
				_, err := bookshelf.GetOrCreate(store, test.title, test.author)
				testutils.AssertError(t, err, test.want)
			})
		}
	})
}
```

Our tests pass.

```sh
~$ go test ./bookshelf
ok      github.com/djangulo/learn-go-with-tests/databases/v6/bookshelf  2.610s
```

## List books

<!-- TODO: start:v6, end: v7 -->

We've done a great job so far testing and implementing our `bookshelf` package. But we still need more functionality in order for it to be a fully fledged storage option.

The `List` function should be the easiest one to implement. Let's start there.

## Write the test first

There's actually not a lot to validate. The `List` function retrieves books if there are, and an empty slice if there aren't any results. We should still return an error in case the `Storer.List` method returns one.

Let's add an optional `search` parameter which would filter results in the `title` or `author` columns. Also, safeguard the tests against nil pointers (in case of errors) and index overflow (as it may not return the length we want).

```go
// bookshelf/crud_test.go
...
import (
	...
	"reflect"
)
...
func TestList(t *testing.T) {
	store := testutils.NewSpyStore(testBooks)

	for _, test := range []struct {
		name, query string
		want        []*bookshelf.Book
	}{
		{"List success: all", "", testBooks},
		{"List success: by author", "shake", testBooks[:2]},
		{"List success: by title", "old man", testBooks[2:2]},
		{"List empty if no match", "this query fails", testBooks[2:2]},
	} {
		t.Run(test.name, func(t *testing.T) {
			books, err := bookshelf.List(store, test.query)
			testutils.AssertNoError(err)
			if !reflect.DeepEqual(books, test.want) {
				t.Errorf("got %v want %v", books, test.want)
			}
		})
	}
}
```

## Try to run the test

Fails as expected, as `List` does not exist.

```sh
~$ go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v7/bookshelf_test [github.com/djangulo/learn-go-with-tests/databases/v7/bookshelf.test]
bookshelf\crud_test.go:147:18: undefined: bookshelf.List
FAIL    github.com/djangulo/learn-go-with-tests/databases/v7/bookshelf [build failed]
FAIL
```

## Write the minimal amount of code for the test to run and check the failing test output

Write the `List` function into `bookshelf-store.go`.

```go
// bookshelf/bookshelf-store.go
func List(store Storer, search string) ([]*Book, error) {
	return nil, nil
}
```

Run the tests:

```sh
~$ go test ./bookshelf
--- FAIL: TestList (0.00s)
    --- FAIL: TestList/List_success:_all (0.00s)
        crud_test.go:152: got [] want [0x904020 0x904060 0x9040a0]
    --- FAIL: TestList/List_success:_by_author (0.00s)
        crud_test.go:152: got [] want [0x904020 0x904060]
    --- FAIL: TestList/List_success:_by_title (0.00s)
        crud_test.go:152: got [] want [0x9040a0]
    --- FAIL: TestList/List_empty_if_no_match (0.00s)
        crud_test.go:152: got [] want []
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v7/bookshelf  5.916s
FAIL
```

## Write enough code to make it pass

We're going to have to add the `Storer` interface methods, and modify both our `PostgreSQLStore` as well as the `SpyStore`.

```go
// bookshelf/bookshelf-store.go
...
type Storer interface {
	...
	List(books *[]*Book, query string) error
}
...
func (s *PostgreSQLStore) List(books *[]*Book, query string) error {
	var rows *sql.Rows
	var err error
	stmt := "SELECT id, title, author FROM books"
	if query != "" {
		stmt += " WHERE LOWER(title) = $1 OR LOWER(author) = $2;"
		rows, err = s.DB.Query(stmt, query, query)
	} else {
		stmt += ";"
		rows, err = s.DB.Query(stmt)
	}
	if err != nil {
		return err
	}

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			fmt.Fprintf(os.Stdout, "error scanning row: %v\n", err)
			continue
		}
		*books = append(*books, &book)
	}
	if err := rows.Close(); err != nil {
		return err
	}
	return nil
}
...
func List(store Storer, search string) ([]*Book, error) {
	books := make([]*Book, 0)
	if search != "" {
		search = strings.ToLower(search)
		err := store.List(&books, search)
		if err != nil {
			return nil, err
		}
		return books, nil
	}

	err := store.List(&books, "")
	if err != nil {
		return nil, err
	}
	return books, nil
}
```

And our `SpyStore`.

```go
// bookshelf/testutils/store.go
func (s *SpyStore) List(books *[]*bookshelf.Book, query string) error {
	if query == "" {
		for _, b := range s.Books {
			*books = append(*books, b)
		}
		return nil
	}
	for _, b := range s.Books {
		if strings.Contains(strings.ToLower(b.Title), query) ||
			strings.Contains(strings.ToLower(b.Author), query) {
			*books = append(*books, b)
		}
	}
	return nil
}
```

Our tests pass:

```sh
~$ go test ./bookshelf
ok      github.com/djangulo/learn-go-with-tests/databases/v7/bookshelf  4.168s
```

## Integrate

Let's write some integration tests for our `List` function. As a matter of fact, we can just reuse the test cases written for `crud_test.go`. The difference here is that the `PostgreSQLStore` returns new objects, and the `ID`s will now be assigned by the store. Let's reuse the loop and test cases, but change the assertion logic.

```go
// bookshelf/integration_test.go
...
import (
	...
	"reflect"
)
...
	t.Run("List", func(t *testing.T) {
		testutils.ResetStore(store)
		for _, b := range testBooks {
			disposable := new(bookshelf.Book)
			store.Create(disposable, b.Title, b.Author)
		}

		for _, test := range []struct {
			name, query string
			want        []*bookshelf.Book
		}{
			{"List success: all", "", testBooks},
			{"List success: by author", "shake", testBooks[:2]},
			{"List success: by title", "old man", testBooks[2:2]},
			{"List empty if no match", "this query fails", []*bookshelf.Book{}},
		} {
			t.Run(test.name, func(t *testing.T) {
				books, err := bookshelf.List(store, test.query)
				testutils.AssertNoError(t, err)

				wantSet := make(map[string]*bookshelf.Book)
				for _, b := range test.want {
					wantSet[b.Title] = b
				}

				for _, b := range books {
					if _, ok := wantSet[b.Title]; !ok {
						t.Errorf("unwanted result %v", *b)
					}
				}
			})
		}
	})
...
```

With this, our tests are still green.

```sh
~$ go test ./bookshelf
ok      github.com/djangulo/learn-go-with-tests/databases/v7/bookshelf  4.023s
```

## Update and Delete

<!-- TODO: start: v7, end: v8 -->

The `Update` and `Delete` functions should be fairly straightforward implementation:

- For `Update`, provide an `ID`, new values and return the updated result. Return an error if both values are empty or if the `ID` is invalid.
- For `Delete`, all we need is the `ID`. It should return the values of the deleted book on success, and an error if the `ID` is invalid or if it does not exist in the store.

## Write the test first

```go
// bookshelf/crud_test.go
...
func TestUpdate(t *testing.T) {

	for _, test := range []struct {
		name     string
		fields   map[string]interface{}
		book     *bookshelf.Book
		wantBook *bookshelf.Book
	}{
		{
			"Update success: title only",
			map[string]interface{}{"title": "Romeo And Juliet"},
			testBooks[1],
			&bookshelf.Book{testBooks[1].ID, "Romeo And Juliet", testBooks[1].Author},
		},
		{
			"Update success: author only",
			map[string]interface{}{"author": "William Shakespeare"},
			testBooks[1],
			&bookshelf.Book{testBooks[1].ID, testBooks[1].Title, "William Shakespeare"},
		},
		{
			"Update success: title+author",
			map[string]interface{}{"title": "Romeo And Juliet", "author": "William Shakespeare"},
			testBooks[1],
			&bookshelf.Book{testBooks[1].ID, "Romeo And Juliet", "William Shakespeare"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := testutils.NewSpyStore(testBooks)
			updated, err := bookshelf.Update(store, test.book.ID, test.fields)
			testutils.AssertNoError(t, err)
			testutils.AssertBooksEqual(t, updated, test.wantBook)
		})
	}

	for _, test := range []struct {
		name   string
		id     int64
		fields map[string]interface{}
		want   error
	}{
		{"Update failure: empty fields", 10, nil, bookshelf.ErrEmptyFields},
		{"Update failure: invalid fields", 10, map[string]interface{}{"isbn": 10}, bookshelf.ErrInvalidFields},
		{"Update failure: zero value ID", 0, map[string]interface{}{"author": "doesn't matter"}, bookshelf.ErrZeroValueID},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := testutils.NewSpyStore(testBooks)
			_, err := bookshelf.Update(store, test.id, test.fields)
			testutils.AssertError(t, err, test.want)
		})
	}
}

func TestDelete(t *testing.T) {

	for _, test := range []struct {
		name string
		id   int64
		want *bookshelf.Book
	}{
		{"Delete success", 10, testBooks[0]},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := testutils.NewSpyStore(testBooks)
			_, err := bookshelf.Delete(store, test.id)
			testutils.AssertNoError(t, err)
			var dummy bookshelf.Book
			err = store.ByID(&dummy, test.id)
			testutils.AssertError(t, err, bookshelf.ErrBookDoesNotExist)
		})
	}

	for _, test := range []struct {
		name string
		id   int64
		want error
	}{
		{"Delete failure: does not exist", 42, bookshelf.ErrBookDoesNotExist},
		{"Delete failure: zero value ID", 0, bookshelf.ErrZeroValueID},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := testutils.NewSpyStore(testBooks)
			_, err := bookshelf.Delete(store, test.id)
			testutils.AssertError(t, err, test.want)
		})
	}
}
```

## Try to run the test

```sh
~$ go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf_test [github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf.test]
bookshelf\crud_test.go:168:20: undefined: bookshelf.Update
bookshelf\crud_test.go:179:48: undefined: bookshelf.ErrEmptyFields
bookshelf\crud_test.go:184:14: undefined: bookshelf.Update
bookshelf\crud_test.go:198:78: undefined: bookshelf.ErrInvalidFields
bookshelf\crud_test.go:201:20: undefined: bookshelf.Delete
bookshelf\crud_test.go:220:14: undefined: bookshelf.Delete
FAIL    github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf [build failed]
FAIL//
```

## Write the minimal amount of code for the test to run and check the failing test output

Insert the new error and the `Delete` and `Update` functions.

```go
// bookshelf/bookshelf-store.go
var (
	...
	ErrEmptyFields = errors.New("fields are empty")
	ErrInvalidFields = errors.New("invalid fields")
)
...
func Update(store Storer, id int64, fields map[string]interface{}) (*Book, error) {
	var book Book
	return &book, nil
}

func Delete(store Storer, id int64) (*Book, error) {
	var book Book
	return &book, nil
}
```

## About `map[string]interface{}`

As it turns out, the type `map[string]interface{}` is ideal for handling updates. First off, most `database/sql` functions and methods accept a variadic `...interface{}` argument, these methods then use introspection to cast the type of the given field to the necessary column. Second, the `map[string]interface{}` type conveniently holds the column name and the value to update it to.

You can pass a map with `1`, `2` or `N` fields, it doesn't matter, validation can be handled in a single helper function. If the tabels were to be extended, instead of adding a new parameter to the `Storer`, the functions and the methods, simply add them to the `fields` param and you'd be done.

We didn't do this for the `Create` methods because `Title` and `Author` are required, although we could have if we wanted to.

Here are some helpers that will assist in validating and cleaning the `fields` parameter.

```go
// bookshelf/bookshelf-store.go
...
type canonicalFieldsMap map[string]bool
var bookFields = map[string]bool{
	"id": true,
	"author": true,
	"title": true,
}
func (c *canonicalFieldsMap) Validate(m map[string]interface{}) error {
	fields := make([]string, 0)
	for f := range m {
		fields = append(fields, f)
	}
	errorFields := make([]string, 0)
	for _, field := range fields {
		if _, ok := (*c)[field]; !ok {
			errorFields = append(errorFields, field)
		}
	}
	if len(errorFields) > 0 {
		return ErrInvalidFields
	}
	return nil
}
func dropField(m map[string]interface{}, toDrop ...string) map[string]interface{} {
	for _, field := range toDrop {
		if _, ok := m[field]; ok {
			delete(m, field)
		}
	}
	return m
}
```

## Try to run the tests

There's something to be said about a thorough error report that tells you exactly what to do.

```sh
~$ go test ./bookshelf
--- FAIL: TestUpdate (0.00s)
    --- FAIL: TestUpdate/Update_success:_title_only (0.00s)
        crud_test.go:188: nil or invalid ID: &{0  }
        crud_test.go:188: got &{0  } want &{22 Romeo And Juliet W. Shakespeare}
    --- FAIL: TestUpdate/Update_success:_author_only (0.00s)
        crud_test.go:188: nil or invalid ID: &{0  }
        crud_test.go:188: got &{0  } want &{22 Romeo & Juliet William Shakespeare}
    --- FAIL: TestUpdate/Update_success:_title+author (0.00s)
        crud_test.go:188: nil or invalid ID: &{0  }
        crud_test.go:188: got &{0  } want &{22 Romeo And Juliet William Shakespeare}
    --- FAIL: TestUpdate/Update_failure:_empty_fields (0.00s)
        crud_test.go:205: wanted an error but didn't get one
        crud_test.go:205: got <nil> want fields are empty
    --- FAIL: TestUpdate/Update_failure:_invalid_fields (0.00s)
        crud_test.go:205: wanted an error but didn't get one
        crud_test.go:205: got <nil> want invalid fields: isbn
    --- FAIL: TestUpdate/Update_failure:_zero_value_ID (0.00s)
        crud_test.go:205: wanted an error but didn't get one
        crud_test.go:205: got <nil> want zero value ID
--- FAIL: TestDelete (0.00s)
    --- FAIL: TestDelete/Delete_success (0.00s)
        crud_test.go:225: wanted an error but didn't get one
        crud_test.go:225: got <nil> want book does not exist
    --- FAIL: TestDelete/Delete_failure:_does_not_exist (0.00s)
        crud_test.go:240: wanted an error but didn't get one
        crud_test.go:240: got <nil> want book does not exist
    --- FAIL: TestDelete/Delete_failure:_zero_value_ID (0.00s)
        crud_test.go:240: wanted an error but didn't get one
        crud_test.go:240: got <nil> want zero value ID
FAIL
FAIL    github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf  3.861s
FAIL
```

Let's get to work!

## Write the minimal amount of code for the test to run and check the failing test output

```go
// bookshelf/bookshelf-store.go
func Update(store Storer, id int64, fields map[string]interface{}) (*Book, error) {
	if id == 0 {
		return nil, ErrZeroValueID
	}
	err := bookFields.Validate(fields)
	if err != nil {
		return nil, err
	}

	fields = dropField(fields, "id") // cannot update id
	var book Book
	err = store.Update(&book, id, fields)
	if err != nil {
		return nil, err
	}
	return &book, nil
}
...
func Delete(store Storer, id int64) (*Book, error) {
	if id == 0 {
		return nil, ErrZeroValueID
	}
	var book Book
	err := store.Delete(&book, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrBookDoesNotExist
		}
		return nil, err
	}
	return &book, nil
}
```

Run the tests:

```sh
~$ go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf
bookshelf\bookshelf-store.go:380:13: store.Update undefined (type Storer has no field or method Update)
bookshelf\bookshelf-store.go:393:14: store.Delete undefined (type Storer has no field or method Delete)
FAIL    github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf [build failed]
FAIL
```

## Write the minimal amount of code for the test to run and check the failing test output

Write the method signatures into `Storer`.

```go
// bookshelf/bookshelf-store.go
...
type Storer interface {
	...
	Update(book *Book, id int64, fields map[string]interface{}) error
	Delete(book *Book, id int64) error
}
...
```

## Try to run the tests


```sh
~$ go test ./bookshelf
# github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf/testutils
bookshelf\testutils\helpers.go:133:32: cannot use store (type *bookshelf.PostgreSQLStore) as type bookshelf.Storer in argument to bookshelf.MigrateDown:
        *bookshelf.PostgreSQLStore does not implement bookshelf.Storer (missing Delete method)
bookshelf\testutils\helpers.go:138:30: cannot use store (type *bookshelf.PostgreSQLStore) as type bookshelf.Storer in argument to bookshelf.MigrateUp:
        *bookshelf.PostgreSQLStore does not implement bookshelf.Storer (missing Delete method)
FAIL    github.com/djangulo/learn-go-with-tests/dat
```

## Write enough code to make it pass

Implement the `Update` and `Delete` methods for the `PostgreSQLStore`.

```go
// bookshelf/bookshelf-store.go
func (s *PostgreSQLStore) Update(book *Book, id int64, fields map[string]interface{}) error {
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	for column, value := range fields {
		columns = append(columns, column)
		values = append(values, value)
	}
	values = append(values, id)

	stmt := fmt.Sprintf("UPDATE books SET (%s) = ROW(", strings.Join(columns, ", "))
	var i int
	for i = 1; i <= len(columns); i++ {
		stmt += fmt.Sprintf("$%d", i)
		if i != len(columns) {
			stmt += ","
		}
	}
	stmt += fmt.Sprintf(") WHERE id = $%d RETURNING id, title, author;", i)

	fmt.Println(stmt)
	row := s.DB.QueryRow(stmt, values...)
	err := row.Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostgreSQLStore) Delete(book *Book, id int64) error {
	stmt := "DELETE FROM books WHERE id = $1 RETURNING title, author;"
	row := s.DB.QueryRow(stmt, id)
	err := row.Scan(&book.Title, &book.Author)
	if err != nil {
		return err
	}
	book.ID = 0
	return nil
}
```

And in the `SpyStore`.
```go
// bookshelf/testutils/store.go
func (s *SpyStore) Update(book *bookshelf.Book, id int64, fields map[string]interface{}) error {
	for _, b := range s.Books {
		if (*b).ID == id {
			*book = *b
		}
	}
	if title, ok := fields["title"]; ok {
		title := title.(string)
		book.Title = title
	}
	if author, ok := fields["author"]; ok {
		author := author.(string)
		book.Author = author
	}
	return nil
}
func (s *SpyStore) Delete(book *bookshelf.Book, id int64) error {
	err := (*s).ByID(book, id)
	if err != nil {
		return err
	}
	var idx int
	for i, b := range s.Books {
		if (*b).ID == id {
			idx = i
			break
		}
	}
	book.ID = 0
	s.Books = append(s.Books[:idx], s.Books[idx+1:]...)
	return nil
}
...
```

And with all that, our tests now pass (phew).

```sh
~$ go test ./bookshelf
ok      github.com/djangulo/learn-go-with-tests/databases/v8/bookshelf  5.241s
```

For integration tests, we can do the same we did for `List`, and reuse the test cases, while swapping out the store and assertions.

```go
// bookshelf/integration_test.go
...
t.Run("Update", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			for _, test := range []struct {
				name     string
				fields   map[string]interface{}
				id       int64
				wantBook *bookshelf.Book
			}{
				{
					"title only",
					map[string]interface{}{"title": "Romeo And Juliet"},
					1,
					&bookshelf.Book{1, "Romeo And Juliet", "W. Shakespeare"},
				},
				{
					"author only",
					map[string]interface{}{"author": "William Shakespeare"},
					1,
					&bookshelf.Book{1, "The Tragedie of Hamlet", "William Shakespeare"},
				},
				{
					"title+author",
					map[string]interface{}{"title": "Romeo And Juliet", "author": "William Shakespeare"},
					1,
					&bookshelf.Book{1, "Romeo And Juliet", "William Shakespeare"},
				},
			} {
				t.Run(test.name, func(t *testing.T) {
					testutils.ResetStore(store)
					for _, b := range testBooks {
						disposable := new(bookshelf.Book)
						store.Create(disposable, b.Title, b.Author)
					}
					_, err := bookshelf.Update(store, test.id, test.fields)
					testutils.AssertNoError(t, err)
				})
			}

		})
		t.Run("failure", func(t *testing.T) {
			for _, test := range []struct {
				name   string
				id     int64
				fields map[string]interface{}
				want   error
			}{
				{"empty fields", 1, nil, bookshelf.ErrEmptyFields},
				{"invalid fields", 1, map[string]interface{}{"isbn": 10}, bookshelf.ErrInvalidFields},
				{"zero value ID", 0, map[string]interface{}{"author": "doesn't matter"}, bookshelf.ErrZeroValueID},
			} {
				t.Run(test.name, func(t *testing.T) {
					_, err := bookshelf.Update(store, test.id, test.fields)
					testutils.AssertError(t, err, test.want)
				})
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			for _, test := range []struct {
				name string
				id   int64
				want *bookshelf.Book
			}{
				{"deletes", 1, testBooks[0]},
			} {
				t.Run(test.name, func(t *testing.T) {
					testutils.ResetStore(store)
					for _, b := range testBooks {
						disposable := new(bookshelf.Book)
						store.Create(disposable, b.Title, b.Author)
					}
					_, err := bookshelf.Delete(store, test.id)
					testutils.AssertNoError(t, err)
					var dummy bookshelf.Book
					err = store.ByID(&dummy, test.id)
					testutils.AssertError(t, err, bookshelf.ErrBookDoesNotExist)
				})
			}
		})
		t.Run("failure", func(t *testing.T) {
			for _, test := range []struct {
				name string
				id   int64
				want error
			}{
				{"Delete failure: does not exist", 42, bookshelf.ErrBookDoesNotExist},
				{"Delete failure: zero value ID", 0, bookshelf.ErrZeroValueID},
			} {
				t.Run(test.name, func(t *testing.T) {
					store := testutils.NewSpyStore(testBooks)
					_, err := bookshelf.Delete(store, test.id)
					testutils.AssertError(t, err, test.want)
				})
			}
		})
	})
...
```

## Wrapping up

We managed to create a working `bookshelf` package, that interacts with a database and maintains certain *behaviors*.

Not only that, our package is extendible: it provides a default storage medium, the `PostgreSQL` database, and if that's not enough, you can create any other storage type you wanted as long as they implement the `Storer` interface.

Note that we only created a _library_ with the tools for implementing persistent storage of Books, we did not provide a program. But with the library there, creating, say, a CLI or a webserver is trivial.

We covered a lot of different topics related to databases and DevOps, here are some more resources to expand further in one of the many directions.

### PostgreSQL

As mentioned before, `PostgreSQL` is a powerful database. We only scratched the surface here of what it's capable of.

The [documentation](https://www.postgresql.org/docs/11/index.html) is excellent: its well organized, covers examples, use cases, concerns and workarounds for A LOT of scenarios. That should be your first place to look for answers related to this awesome database.

This is just one of many databases though, below is by-no-means-comprehensive list:

- [MySQL](https://www.mysql.com/).
- [SQLite](https://sqlite.org/index.html).
- [MariaDB](https://mariadb.org/)
- [SQL Server](https://www.microsoft.com/en-us/sql-server/)
- [Cassandra](https://cassandra.apache.org/)

In this guide we used the term _database_ freely to refer to a Relational Database. In the strictest sense, this is very inaccurate, as [there are many, many different types of databases](https://en.wikipedia.org/wiki/Outline_of_databases#Types_of_databases).

### Continuous integration

Continuous Integration, or CI, was mentioned briefly in this chapter. It is simply merging all (tested) code into the main codebase.

It's a natural segueway that follows TDD.

While `CI` itself is a concept, you need a framework to itegrate it into your coding flow.

Generally, you create a series of tasks that need to be executed in sequence, that can either pass or fail (sound familiar?). These tasks can range from building your application, to testing it in isolation, to publishing it online or merging with a `master` repository.

Here are some of the most recognized CI/CD frameworks out there:

| Project   | Homepage                 | Available as | Remarks                                                                                                              |
| :-------- | :----------------------- | :----------- | :------------------------------------------------------------------------------------------------------------------- |
| Travis CI | https://travis-ci.org/   | SaaS         | Free for open source projects.                                                                                       |
| Circle CI | https://travis-ci.org/   | SaaS         | Free for open source projects.                                                                                       |
| Concourse | https://concourse-ci.org | Self-hosted  | Written in `Go`, it's relatively simple to use, and provides a `CLI`.                                                |
| Jenkins   | https://jenkins.io/      | Self-hosted  | Probably the easiest to setup. Oldest in the market. Immense community. Ocean of plugins to extend. Written in Java. |
| GoCD      | https://www.gocd.org/    | Self-hosted  | Written in Java (I know, right?). Extendible via plugins.                                                            |

### Database normalization

If you're interested in databases and Database design, I suggest you familiarize yourself with **database normalization**. Below are some useful resources.

-   [Database Normalization](https://en.wikipedia.org/wiki/Database_normalization).
-   [A Simple Guide to Five Normal Forms in Relational Database Theory](http://www.bkent.net/Doc/simple5.htm).
-   [Relational Database Design/Normalization](https://en.wikibooks.org/wiki/Relational_Database_Design/Normalization).
