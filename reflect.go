package rencode

import (
	"fmt"
	"reflect"
)

func init() {
	byteSliceType = reflect.TypeOf([]byte{})
	stringType = reflect.TypeOf("")
}

var (
	byteSliceType reflect.Type
	stringType    reflect.Type
)

func isInt(kind reflect.Kind) bool {
	return kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64
}

func isUint(kind reflect.Kind) bool {
	return kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64
}

func isFloat(kind reflect.Kind) bool {
	return kind == reflect.Float32 || kind == reflect.Float64
}

func isCompatibleSlice(src, dest reflect.Type) bool {
	// log.Printf("%v = %v: %b", dest, src, src.Ass)

	return (src.Kind() == reflect.Array || src.Kind() == reflect.Slice) &&
		dest.Kind() == reflect.Slice &&
		src.Elem().AssignableTo(dest.Elem())
}

func convertAssign(src, dest interface{}) error {
	srcType := reflect.TypeOf(src)
	srcVal := reflect.ValueOf(src)
	destType := reflect.TypeOf(dest)
	destVal := reflect.ValueOf(dest)

	// Defer the pointer so we can work with the actual type
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
		destVal = destVal.Elem()
	}

	// Check if we can at least set the value of the destination
	if !destVal.CanSet() {
		return fmt.Errorf("cannot set value of type %s", destType)
	}

	if srcType.AssignableTo(destType) {
		destVal.Set(srcVal)
		return nil
	}

	destKind := destType.Kind()
	srcKind := srcType.Kind()

	if isInt(destKind) && isInt(srcKind) {
		destVal.SetInt(srcVal.Int())
		return nil
	} else if isUint(destKind) && isUint(srcKind) {
		destVal.SetUint(srcVal.Uint())
		return nil
	} else if isFloat(srcKind) && isFloat(destKind) {
		destVal.SetFloat(srcVal.Float())
		return nil
	} else if isCompatibleSlice(srcType, destType) {

		n := srcVal.Len()
		slice := reflect.MakeSlice(destVal.Type(), srcVal.Len(), srcVal.Cap())

		for i := 0; i < n; i++ {
			slice.Index(i).Set(srcVal.Index(i))
		}

		destVal.Set(slice)
		return nil
	} else if srcKind == reflect.String && destType.AssignableTo(byteSliceType) {
		// Special case when the source is a string and the destination is a byte slice
		destVal.Set(srcVal.Convert(byteSliceType))
		return nil
	} else if srcType.AssignableTo(byteSliceType) && destKind == reflect.String {
		// Special case when the source is a byte slice and the destination is a string
		destVal.Set(srcVal.Convert(stringType))
		return nil
	} else if isInt(srcKind) && isFloat(destKind) {
		destVal.SetFloat(srcVal.Convert(destType).Float())
		return nil
	}

	// TODO: refactor the below function to use reflection
	return convertAssignInteger(src, dest)
}
