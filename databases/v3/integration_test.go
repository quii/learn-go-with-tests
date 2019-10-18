package main

import "testing"

const queryTables = `
SELECT tablename, tableowner
FROM pg_catalog.pg_tables
WHERE
	schemaname != 'pg_catalog'
	AND
	schemaname != 'information_schema';`

type pgTable struct {
	tableOwner string `sql:"tableowner"`
	tableName  string `sql:"tablename"`
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
		_, err := MigrateDown(dummyWriter, store, "migrations", -1)
		if err != nil {
			t.Errorf("migration down failed: %v", err)
		}

		rows, err := store.db.Query(queryTables)
		if err != nil {
			t.Errorf("received error querying rows: %v", err)
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
	t.Run("idempotency", func(t *testing.T) {
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
