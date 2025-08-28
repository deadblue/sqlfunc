package sqlfunc

import (
	"database/sql"
	"errors"
	"reflect"
)

type (
	destFunc[R any] func(result *R) (dest []any)
)

var (
	_TypeBoolPtr        = reflect.TypeFor[*bool]()
	_TypeStringPtr      = reflect.TypeFor[*string]()
	_TypeByteSlicePtr   = reflect.TypeFor[*[]byte]()
	_TypeIntPtr         = reflect.TypeFor[*int]()
	_TypeInt8Ptr        = reflect.TypeFor[*int8]()
	_TypeInt16Ptr       = reflect.TypeFor[*int16]()
	_TypeInt32Ptr       = reflect.TypeFor[*int32]()
	_TypeInt64Ptr       = reflect.TypeFor[*int64]()
	_TypeUintPtr        = reflect.TypeFor[*uint]()
	_TypeUint8Ptr       = reflect.TypeFor[*uint8]()
	_TypeUint16Ptr      = reflect.TypeFor[*uint16]()
	_TypeUint32Ptr      = reflect.TypeFor[*uint32]()
	_TypeUint64Ptr      = reflect.TypeFor[*uint64]()
	_TypeFloat32Ptr     = reflect.TypeFor[*float32]()
	_TypeFloat64Ptr     = reflect.TypeFor[*float64]()
	_TypeSqlRawBytesPtr = reflect.TypeFor[*sql.RawBytes]()

	_TypeScannable  = reflect.TypeFor[Scannable]()
	_TypeSqlScanner = reflect.TypeFor[sql.Scanner]()

	errUnscannableResult = errors.New("can not get dest from result")
)

func isDestType(pt reflect.Type) bool {
	if pt.AssignableTo(_TypeSqlScanner) {
		return true
	}
	switch pt {
	case _TypeBoolPtr, _TypeStringPtr, _TypeByteSlicePtr,
		_TypeIntPtr, _TypeInt8Ptr, _TypeInt16Ptr,
		_TypeInt32Ptr, _TypeInt64Ptr,
		_TypeUintPtr, _TypeUint8Ptr, _TypeUint16Ptr,
		_TypeUint32Ptr, _TypeUint64Ptr,
		_TypeFloat32Ptr, _TypeFloat64Ptr,
		_TypeSqlRawBytesPtr:
		return true
	default:
		return false
	}
}

func makeDestFunc[R any](columns []string) (f destFunc[R], err error) {
	pt := reflect.TypeFor[*R]()
	if pt.AssignableTo(_TypeScannable) {
		f = func(result *R) []any {
			return any(result).(Scannable).Dest(columns)
		}
		return
	}
	if len(columns) == 1 && isDestType(pt) {
		f = func(result *R) (dest []any) {
			return []any{result}
		}
		return
	}
	if rt := reflect.TypeFor[R](); rt.Kind() == reflect.Struct {
		mapping := getStructMapping(rt)
		f = func(result *R) (dest []any) {
			rv := reflect.ValueOf(*result)
			for _, column := range columns {
				field, found := mapping[column]
				if found {
					fv := rv.FieldByName(field)
					// Put field point to dest
					dest = append(dest, fv.Addr().Interface())
				} else {
					// Put void to dest
					dest = append(dest, Void{})
				}
			}
			return
		}
	} else {
		err = errUnscannableResult
	}
	return
}
