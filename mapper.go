package recordmapper

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
)

var (
	ErrNumOfColumnGroupsNotMatch = errors.New("number of column groups not match")
	ErrNumOfColumnsNotMatch      = errors.New("number of columns not match")
)

type Rows interface {
	Columns() ([]string, error)
	Scan(dest ...interface{}) error
}

func Scan(rows Rows, delim string, dest ...any) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	var (
		columnGroups = make([][]string, 1)
		groupIndex   int
	)
	for _, column := range columns {
		if column == delim {
			groupIndex++
			columnGroups = append(columnGroups, nil)
			continue
		}
		columnGroups[groupIndex] = append(columnGroups[groupIndex], column)
	}
	if len(columnGroups) != len(dest) {
		return fmt.Errorf("%w: %d != %d", ErrNumOfColumnGroupsNotMatch, len(columnGroups), len(dest))
	}
	var values []any
	for i, group := range columnGroups {
		v, err := bind(group, dest[i])
		if err != nil {
			return err
		}
		values = append(values, v...)

		if len(columnGroups)-1 != i {
			values = append(values, new(any))
		}
	}
	if len(values) != len(columns) {
		return fmt.Errorf("%w: %d != %d", ErrNumOfColumnsNotMatch, len(values), len(columns))
	}

	return rows.Scan(values...)
}

type MapperFunc func(columns []string, v reflect.Value) []any

var cache = make(map[reflect.Type]MapperFunc)

func bind(columns []string, dest any) ([]any, error) {
	v := reflect.ValueOf(dest)
	mapper, ok := cache[v.Type()]
	if !ok {
		var err error
		mapper, err = createMapper(v)
		if err != nil {
			return nil, err
		}
		cache[v.Type()] = mapper
	}
	return mapper(columns, v), nil
}

func createMapper(v reflect.Value) (MapperFunc, error) {
	if v.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("dest must be a pointer")
	}

	st := v.Elem()
	if st.Kind() != reflect.Struct {
		return nil, fmt.Errorf("dest must be a pointer to struct")
	}

	t := st.Type()
	fields := reflect.VisibleFields(t)
	fieldIndexMap := make(map[string][]int)
	for _, f := range fields {
		dbTag := f.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			fieldIndexMap[dbTag] = f.Index
			continue
		}
		fieldIndexMap[strcase.ToSnake(f.Name)] = f.Index
	}

	return func(columns []string, v reflect.Value) []any {
		var values []any
		base := reflect.Indirect(v)
		for _, column := range columns {
			v := base.FieldByIndex(fieldIndexMap[column])
			values = append(values, v.Addr().Interface())
		}
		return values
	}, nil
}
