package secrecy

import (
	"fmt"
	"reflect"
)

// Interface for securely erasing values from memory. You may want to implement
// it if you're struct contains unexported fields that can't be ignored using
// `Zeroizer:"ignore"` struct tag.
// This interface MUST be implemented using a struct receiver and not a pointer
// receiver.
type Zeroizer interface {
	Zeroize()
}

// Zeroize recursively changes to zeros memory pointed by the given value.
// Panics if data is or contains a struct with non ignored unexported fields and
// doesn't implements Zeroizer.
//
// This function zeroize all mutable memory, be sure to not share that memory
// as it may lead to unexpected behavior. All pointers (map, slice, pointers)
// type are emptied before being replaced by nil.
//
// Go strings are immutable and can't be zeroized (without using unsafe),
// if you want your secret to be fully zeroized, store them as []byte.
func Zeroize(v any) {
	vValue := reflect.ValueOf(v)
	if vValue.IsZero() {
		return
	}

	// Deref value.
	switch vValue.Kind() {
	case reflect.Pointer, reflect.Interface:
		vValue = vValue.Elem()
	default:
		break
	}

	vType := vValue.Type()

	switch vType.Kind() {
	case reflect.Pointer, reflect.Interface:
		if !vValue.IsNil() {
			Zeroize(vValue.Interface())
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < vValue.Len(); i++ {
			iValue := vValue.Index(i)
			Zeroize(iValue.Interface())
		}

	case reflect.Struct:
		fields := reflect.VisibleFields(vType)
		hasUnexportedField := false
		for _, f := range fields {
			if f.Tag.Get("zeroize") == "ignore" {
				continue
			}

			if !f.IsExported() {
				hasUnexportedField = true
				break
			}
		}
		if hasUnexportedField {
			var zeroizer Zeroizer
			if vType.Implements(reflect.TypeOf(&zeroizer).Elem()) {
				vValue.Interface().(Zeroizer).Zeroize()
			} else {
				panic(fmt.Errorf("type %v contains unexported field and doesn't implements Zeroizer interface", vType.Name()))
			}
		}

		for _, f := range fields {
			if f.IsExported() && f.Tag.Get("zeroize") != "ignore" {
				vField := vValue.FieldByName(f.Name)
				Zeroize(vField.Interface())
			}
		}

	case reflect.Map:
		if !vValue.IsNil() {
			iter := vValue.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()
				Zeroize(k.Interface())
				Zeroize(v.Interface())
				// Remove entry.
				vValue.SetMapIndex(k, reflect.Value{})
			}
		}

	case reflect.Bool, reflect.String,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.Func,
		reflect.UnsafePointer,
		reflect.Chan:
		// noop

	case reflect.Invalid:
		fallthrough
	default:
		panic("unexpected reflect.Kind")
	}

	// Finally set value to zero if possible.
	if vValue.CanSet() {
		vValue.SetZero()
	}
}
