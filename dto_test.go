package dto

import (
	"testing"
)

func TestSetFieldTag(t *testing.T) {
	newTag := "tag_value"
	SetFieldTag(newTag)
	if newTag != dtoFieldTag {
		t.Fatal("setting custom field tag error")
	}
}

func TestRequestToDTO(t *testing.T) {
	SetFieldTag("dto")

	struct1 := struct {
		Field1 int
	}{
		1,
	}
	struct2 := struct {
		Field2 int `dto:"field_t2"`
	}{
		2,
	}
	struct3 := struct {
		CustomField3 int `dto:"custom_field3"`
	}{
		3,
	}
	struct4 := struct {
		Field4 int `dto:"field4"`
		Field5 int `dto:"-"`
	}{
		4,
		5,
	}

	dtoStruct := struct {
		Field1  int
		FieldT2 int
		Field3  int `dto:"custom_field3"`
		Field4  int
		Field5  int
	}{}

	err := RequestToDTO(&dtoStruct, struct1, struct2, struct3, &struct4)

	if err != nil {
		t.Fatal("RequestToDTO assign error", err)
	}

	if dtoStruct.Field1 != 1 {
		t.Fatal("cannot assign tag capitalized to target without tag")
	}

	if dtoStruct.FieldT2 != 2 {
		t.Fatal("cannot assign tag capitalized snake case to target without tag")
	}

	if dtoStruct.Field3 != 3 {
		t.Fatal("cannot assign tag to target with tag")
	}

	if dtoStruct.Field4 != 4 {
		t.Fatal("cannot assign struct by pointer")
	}

	if dtoStruct.Field5 != 0 {
		t.Fatal("skip tag has no effect")
	}
}

func Test_tagToFieldName(t *testing.T) {
	if tagValueToFieldName("testfield") != "Testfield" {
		t.Fatal("simple value conversation error")
	}
	if tagValueToFieldName("test_field") != "TestField" {
		t.Fatal("snake case value conversation error")
	}
}
