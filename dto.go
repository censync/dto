package dto

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	dtoFieldTag = "dto"
)

// Change default tag `dto` to custom
func SetFieldTag(tag string) {
	if tag == "" || strings.Trim(tag, " ") == "" {
		panic("empty tag name")
	}

	dtoFieldTag = tag
}

func RequestToDTO(dst interface{}, src ...interface{}) error {
	dstRv := reflect.ValueOf(dst)
	if dstRv.Kind() != reflect.Ptr || dstRv.IsNil() {
		return errors.New(fmt.Sprintf("cannot convert to dto: %s", reflect.TypeOf(dst)))
	}

	if len(src) > 0 {
		for index, val := range src {
			srcRv := reflect.ValueOf(val)
			switch srcRv.Kind() {
			case reflect.Struct:
				err := parseStruct(dstRv, srcRv)
				if err != nil {
					return err
				}
			case reflect.Ptr, reflect.UnsafePointer:
				if srcRv.IsNil() {
					return errors.New(fmt.Sprintf("cannot convert source to dto: %d argument is nil", index))
				}
				err := parseStruct(dstRv, srcRv.Elem())
				if err != nil {
					return err
				}

			default:
				return errors.New(fmt.Sprintf("cannot convert source to dto: %s", srcRv.Type()))
			}
		}
	}
	return nil
}
func parseStruct(dst reflect.Value, srcVal reflect.Value) error {
	if srcVal.Type().Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("unsupported type for dest argument: %s", srcVal.Type().Key()))
	}
	srcType := srcVal.Type()
	fieldsCount := srcVal.NumField()
	for i := 0; i < fieldsCount; i++ {
		fieldVal := srcVal.Field(i)

		for fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() {
			fieldVal = fieldVal.Elem()
		}

		tag := srcType.Field(i).Tag.Get(dtoFieldTag)

		if tag == "" {
			return errors.New(fmt.Sprintf("empty dto tag \"%s\" value for field: %s", dtoFieldTag, srcType.Field(i).Name))
		}

		// Looking "tag_field_name", converted to "TagFieldName"
		// for public field in target structure
		targetField := dst.Elem().FieldByName(tagValueToFieldName(tag))

		// Looking for simple value
		if !targetField.IsValid() {
			targetField = dst.Elem().FieldByName(tag)
		}

		// Looking for tag value in target structure
		if !targetField.IsValid() {
			targetField = findFieldByTagValue(dst.Elem(), tag)
		}

		if !targetField.IsValid() {
			return errors.New(fmt.Sprintf("not found field: %s", tag))
		}

		if !targetField.CanSet() {
			return errors.New(fmt.Sprintf("not writable field: %s", tag))
		}

		if targetField.Type() != srcType.Field(i).Type {
			return errors.New(fmt.Sprintf("incompatible types source %s and target %s", srcType.Field(i).Type, targetField.Type()))
		}

		targetField.Set(fieldVal)
	}
	return nil
}

//	Converts tag field value to capitalized field name
//  "fieldname" => "Fieldname"
//  "field_name" => "FieldName"
func tagValueToFieldName(src string) string {
	src = strings.Replace(src, "_", " ", -1)
	return strings.Replace(strings.Title(src), " ", "", -1)
}

func findFieldByTagValue(dst reflect.Value, tag string) reflect.Value {
	if dst.Kind() == reflect.Ptr {
		panic("dst cannot be pointer")
	}
	fieldsCount := dst.NumField()
	fieldType := dst.Type()
	for i := 0; i < fieldsCount; i++ {
		if fieldType.Field(i).Tag.Get(dtoFieldTag) == tag {
			return dst.Field(i)
		}
	}
	return reflect.Value{}
}
