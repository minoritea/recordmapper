package recordmapper

import (
	"fmt"
	"reflect"
	"testing"
)

type testRows []struct {
	name  string
	value any
}

func (t testRows) Columns() ([]string, error) {
	var names []string
	for _, c := range t {
		names = append(names, c.name)
	}
	return names, nil
}

func (t testRows) Scan(dest ...any) error {
	for i, c := range t {
		reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(c.value))
	}
	return nil
}

func TestScan(t *testing.T) {
	var testRows = testRows{
		{"id", 1},
		{"name", "test"},
		{"NULL", ""},
		{"id", 2},
		{"name", "test2"},
	}
	type testStruct struct {
		ID   int
		Name string
	}
	var ts, ts2 testStruct
	err := Scan(testRows, "NULL", &ts, &ts2)
	if err != nil {
		t.Fatal(err)
	}
	if ts.ID != 1 || ts.Name != "test" {
		t.Fatalf("scan error: %+v", ts)
	}
}

var rows = testRows{
	{"id", 1},
	{"name", "Alice"},
	{"NULL", ""},
	{"id", 1},
	{"age", 20},
}

func ExampleScan() {
	type T1 struct {
		ID   int
		Name string
	}
	type T2 struct {
		ID  int
		Age int
	}
	var (
		t1 T1
		t2 T2
	)
	// CREATE TABLE T1 (id INTEGER, name TEXT);
	// INSERT INTO T1 VALUES (1, 'Alice');
	// CREATE TABLE T2 (id INTEGER, age INTEGER);
	// INSERT INTO T2 VALUES (1, 20);
	//
	// rows := db.QueryContext(ctx, "SELECT t1.*, NULL, t2.* FROM test1 t1 JOIN test2 t2 ON t1.id = t2.id")
	err := Scan(rows, "NULL", &t1, &t2)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("name: %s, age: %d\n", t1.Name, t2.Age)
	// Output: name: Alice, age: 20
}
