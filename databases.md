# Databases

Oftentimes when creating software, it's necessary to save (or, more precisely, _persist_) some application state.

As an example, when you log into your online banking system, the system has to

1. Check that it's really you accessing the system (this is called _authentication_, and is beyond the scope of this chapter)
2. Retrieve some information from _somewhere_ and show it to the user (you).

Information that is stored and meant to be long-lived is said to be [_persisted_](<https://en.wikipedia.org/wiki/Persistence_(computer_science)>), usually on a medium that can reliably reproduce the data stored.

Some storage systems, like the filesystem, can be effective for one-off or small amounts of storage, but they fall short for larger application, for a number of reasons.

This is why most software applications, large and small, opt for storage systems that can provide

-   Reliability: The data you want is there when you need it
-   Concurrency: Imagine thousands of users accessing simultaneously.
-   Consistent: You expect the same inputs to produce the same results
-   Durable: Data should remain there even in case of a system failure (power outage or system crash)

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

The easiest (and cleanest) way of getting `PostgreSQL` up and running is by using `docker`. This will create the database and user

-   [`Docker` installation instructions](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
-   See [https://hub.docker.com/\_/postgres](https://hub.docker.com/_/postgres) for more details on how to use this image.

```bash
~$ docker run \
    --name my-postgres \ # name of the instance
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

```bash
~$ sudo apt-get upgrade
~$ sudo apt-get install postgresql postgresql-contrib
```

PostgreSQL installs and initializes a database called `postgres`, and a user also called `postgres`. Since this is a system-wide install, we don't want to pollute this main database with this application's tables (`PostgreSQL` uses these to store administrative data), so we will have to create a user and a database.

Note that inside the `psql` shell, anything after a double hyphen (`--`) is considered a comment

```
~$ sudo -i -u postgres # this will switch you to the postgres user
~$ psql
psql (10.10 (Ubuntu 10.10-0ubuntu0.18.04.1))
Type "help" for help

postgres=# CREATE USER bookshelf_user WITH CREATEDB PASSWORD 'secret-password';
CREATE ROLE
postgres=# CREATE DATABASE bookshelf_db OWNER bookshelf_user;
CREATE DATABASE

postgres=# -- you can view users and databases with the commands \du and \l respectively
```

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

We will initiall write a (simple) tool to handle our (also simple) migrations, then we will be creating a CRUD program to interact with our spiffy, real database.

## Write the test first

```golang
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

```bash
# command-line-arguments [command-line-arguments.test]
.\migrate_test.go:9:10: undefined: MigrateUp
FAIL    command-line-arguments [build failed]
FAIL
```

## Write enough code to make it pass

We will take some liberties with writing code, keeping it in small batches and testing as we go.

We will write some boilerplate that will assist in creating the database connection. We need an `*sql.DB` instance that we can pass on to our tests first. As we learned in the `Dependency Injection` chapter, we should make a helper method so we can acquire an instance from anywhere in our application, while passing in the dependencies, in our case, the database connection string.

Before we get started in the happy-path, we need to make sure our `MigrateUp` function works as expected (unit-test) before we attempt to call it on the real database (integration-test).

Let's start by creating an interface, which will allow us to mock the database functionality, as well as a `NewStore` function that will simplify its creation.

```golang
// bookshelf-store.go
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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

const removeTimeout = 20 * time.Second

func NewStore(dbConnStr) (*Store, func()) {
	db, err := sql.Open("postgres", dbConnStr)
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
				fmt.Fprintf(
					os.Stderr,
					"error closing connection to database, retrying in %v: %v\n",
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
	return &Store{db: db}, remove
}

func (s *Store) ApplyMigration(stmt string) error {
	_, err := s.db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}
```

---

### Note on exponential backoff (`remove` anonymous function)

Generally when calling external services you want to account for the possibility of failure, yet the reasons why the service could fail are numerous. The database is an external service because (generally) databases run in a client-server paradigm. In our case, the database connection could fail to close for a number of reasons. This pattern is useful to give the service some time to correct itself or reset before attempting our call again. This snippet does so by waiting on powers of 2 incrementally: 1, 2, 4, 8, 16, ..., N seconds, with a hard limit set at `removeTimeout` seconds.

This is here as a demonstration mostly, where it will prove useful is during the integration tests. The test database that we create will have be destroyed after running tests, thus, it can't have read or write operations running.

Code explained

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

```golang
MigrateUp(store Storer, dir string, num int)
```

There's a lot to our `MigrateUp` function, so let's break it down and test each step.

1. It needs to check whether the `dir` passed in exists.
2. It needs to get all the filenames inside `dir`, and allow us to iterate over them.
3. It needs to run only `up` migrations.
4. It needs to run migrations `1` through `num`, or all of them if `num == -1`.
5. Lastly, it needs to report on the success of each migration run, if a migration fails, the entire process should be halted.

Seeing as we will have a `MigrateDown` as well, and so far they seem only to differ on step `3`, we can use a little foresight and create a utility function `migrate`, which will be used by both variants.

## Write the test first

Let's write our mock store (and test) first, so we can test the `migrate` function in isolation.

```golang
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
	mig := s.migrations[name]
	mig.name = name
	mig.stmt = stmt
	mig.called++
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

```golang
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
