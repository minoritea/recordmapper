package recordmapper

import (
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
