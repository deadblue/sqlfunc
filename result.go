package sqlfunc

import (
	"reflect"
	"strings"
)

const (
	_TagName = "db"
)

/*
Scannable is an interface that can be implemented by query result struct, to
avoid the performance overhead of reflection.

When result struct implements [Scannable] interface, sqlfunc calls Dest to get
corresponding dests for columns.

Example:

	type UserResult struct{
		Foo string
		Bar int
	}

	func (r *UserResult) Dest(columns []string) (dest []any) {
		for _, column := range columns {
			switch column {
			case "foo":
				dest = append(dest, &r.Foo)
			case "bar":
				dest = append(dest, &r.Bar)
			default:
				dest = append(dest, sqlfunc.Void{})
			}
		}
		return
	}
*/
type Scannable interface {
	// Dest returns scanning dest for columns.
	Dest(columns []string) (dest []any)
}

var (
	_ColumnMappingCache = map[reflect.Type]map[string]string{}

	_ScannaleType = reflect.TypeFor[Scannable]()
)

func getColumnMapping(rt reflect.Type) map[string]string {
	mapping, ok := _ColumnMappingCache[rt]
	if !ok {
		mapping = make(map[string]string)
		for i := range rt.NumField() {
			ft := rt.Field(i)
			if desc, found := ft.Tag.Lookup(_TagName); found {
				if strings.ContainsRune(desc, ',') {
					props := strings.Split(desc, ",")
					mapping[props[0]] = ft.Name
				} else {
					mapping[desc] = ft.Name
				}
			} else {
				column := pascalToSnake(ft.Name)
				mapping[column] = ft.Name
			}
		}
		// Put mapping to cache
		_ColumnMappingCache[rt] = mapping
	}
	return mapping
}

func getResultDest(result any, columns []string) (dest []any) {
	if sr, ok := result.(Scannable); ok {
		return sr.Dest(columns)
	}
	rv := reflect.ValueOf(result)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Struct {
		columnMapping := getColumnMapping(rv.Type())
		for _, column := range columns {
			field, found := columnMapping[column]
			if found {
				fv := rv.FieldByName(field)
				// Put field point to dest
				dest = append(dest, fv.Addr().Interface())
			} else {
				// Put void to dest
				dest = append(dest, Void{})
			}
		}
	} else {
		// TODO: Support other type
	}
	return
}

// preprocessResult initializes mapping cache to result type
func preprocessResult[R any]() {
	rt := reflect.TypeFor[R]()
	if pt := reflect.PointerTo(rt); pt.AssignableTo(_ScannaleType) {
		return
	}
	_ = getColumnMapping(rt)
}
