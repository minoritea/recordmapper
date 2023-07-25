package sqlite

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/minoritea/recordmapper"
)

func TestScanSqlite(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	_, err = db.ExecContext(ctx, "CREATE TABLE test (id INTEGER, name TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx, "CREATE TABLE test2 (id INTEGER, age INTEGER)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx, "INSERT INTO test VALUES (1, 'Alice')")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx, "INSERT INTO test2 VALUES (1, 20)")
	if err != nil {
		t.Fatal(err)
	}
	rows, err := db.QueryContext(ctx, "SELECT test.*, NULL, test2.* FROM test JOIN test2 ON test.id = test2.id")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	type testStruct struct {
		NameStruct struct {
			ID   int
			Name string
		}
		AgeStruct struct {
			ID  int
			Age int
		}
	}
	for rows.Next() {
		var ts testStruct
		err = recordmapper.Scan(rows, "NULL", &ts.NameStruct, &ts.AgeStruct)
		if err != nil {
			t.Fatal(err)
		}
		if ts.NameStruct.ID != 1 || ts.NameStruct.Name != "Alice" {
			t.Fatalf("scan error: %+v", ts)
		}
		if ts.AgeStruct.ID != 1 || ts.AgeStruct.Age != 20 {
			t.Fatalf("scan error: %+v", ts)
		}
	}
	err = rows.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}
}
