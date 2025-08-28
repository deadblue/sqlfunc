package sqlfunc

import (
	"reflect"
	"strings"
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

const (
	_TagKey = "sql"
)

var (
	_StructMappingCache = map[reflect.Type]map[string]string{}
)

func getStructMapping(rt reflect.Type) map[string]string {
	mapping, found := _StructMappingCache[rt]
	if !found {
		mapping = make(map[string]string)
		for i := range rt.NumField() {
			ft := rt.Field(i)
			if desc, found := ft.Tag.Lookup(_TagKey); found {
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
		_StructMappingCache[rt] = mapping
	}
	return mapping
}

// preprocessResult initializes mapping cache to result type
func preprocessResult[R any]() {
	if pt := reflect.TypeFor[*R](); pt.AssignableTo(_TypeScannable) {
		return
	}
	if rt := reflect.TypeFor[R](); rt.Kind() == reflect.Struct {
		getStructMapping(rt)
	}
}
