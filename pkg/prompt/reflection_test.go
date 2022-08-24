// Copyright (c) 2022 The Swedish Internet Foundation
//
// Distributed under the MIT License. (See accompanying LICENSE file or copy at
// <https://opensource.org/licenses/MIT>.)

package prompt

import (
	"reflect"
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

type randomStruct struct {
	SomeStr       string        `xml:"string"`
	SomeInt       int           `xml:"int"`
	SomeBool      bool          `xml:"bool"`
	SomeStruct    listStruct    `xml:"struct"`
	SomeBoolPtr   *bool         `xml:"bool_ptr"`
	SomeStringPtr *string       `xml:"string_ptr"`
	SomeIntPtr    *int          `xml:"int_ptr"`
	SomeStructPtr *randomStruct `xml:"struct_ptr"`
}

type listStruct struct {
	Strings    []string        `xml:"strings"`
	Ints       []int           `xml:"ints"`
	StructsPtr []*randomStruct `xml:"structs_ptr"`
	Structs    []randomStruct  `xml:"structs"`
}

func TestFieldValue_ValueAsString(t *testing.T) {
	var nilPtr *randomStruct

	for _, tc := range []struct {
		name           string
		fv             *FieldValue
		expectedString string
	}{
		{
			name: "int",
			fv: &FieldValue{
				Value: reflect.ValueOf(12345),
				Kind:  reflect.Int,
			},
			expectedString: "12345",
		},
		{
			name: "string",
			fv: &FieldValue{
				Value: reflect.ValueOf("hello"),
				Kind:  reflect.String,
			},
			expectedString: "hello",
		},
		{
			name: "bool",
			fv: &FieldValue{
				Value: reflect.ValueOf(true),
				Kind:  reflect.Bool,
			},
			expectedString: "true",
		},
		{
			name: "slice",
			fv: &FieldValue{
				Value:    reflect.ValueOf([]string{"hello", "there"}),
				Kind:     reflect.Slice,
				ListKind: reflect.String,
			},
			expectedString: "[hello there]",
		},
		{
			name: "struct",
			fv: &FieldValue{
				Value: reflect.ValueOf(randomStruct{SomeStr: "hello", SomeInt: 123}),
				Kind:  reflect.Struct,
			},
			expectedString: "{SomeStr:hello SomeInt:123 SomeBool:false SomeStruct:{Strings:[] Ints:[] StructsPtr:[] Structs:[]} SomeBoolPtr:<nil> SomeStringPtr:<nil> SomeIntPtr:<nil> SomeStructPtr:<nil>}",
		},
		{
			name: "pointer",
			fv: &FieldValue{
				Value: reflect.ValueOf(&randomStruct{SomeStr: "hello", SomeInt: 123}),
				Kind:  reflect.Pointer,
			},
			expectedString: "{SomeStr:hello SomeInt:123 SomeBool:false SomeStruct:{Strings:[] Ints:[] StructsPtr:[] Structs:[]} SomeBoolPtr:<nil> SomeStringPtr:<nil> SomeIntPtr:<nil> SomeStructPtr:<nil>}",
		},
		{
			name: "struct slice",
			fv: &FieldValue{
				Value: reflect.ValueOf([]randomStruct{
					{SomeStr: "flower", SomeInt: 321},
					{SomeStr: "hello", SomeInt: 123},
				}),
				Kind:     reflect.Slice,
				ListKind: reflect.Struct,
			},
			expectedString: "[{SomeStr:flower SomeInt:321 SomeBool:false SomeStruct:{Strings:[] Ints:[] StructsPtr:[] Structs:[]} SomeBoolPtr:<nil> SomeStringPtr:<nil> SomeIntPtr:<nil> SomeStructPtr:<nil>} {SomeStr:hello SomeInt:123 SomeBool:false SomeStruct:{Strings:[] Ints:[] StructsPtr:[] Structs:[]} SomeBoolPtr:<nil> SomeStringPtr:<nil> SomeIntPtr:<nil> SomeStructPtr:<nil>}]",
		},
		{
			name: "struct pointer slice",
			fv: &FieldValue{
				Value: reflect.ValueOf([]*randomStruct{
					{SomeStr: "flower", SomeInt: 321},
					{SomeStr: "hello", SomeInt: 123},
				}),
				Kind:     reflect.Slice,
				ListKind: reflect.Pointer,
			},
			expectedString: "[{SomeStr:flower SomeInt:321 SomeBool:false SomeStruct:{Strings:[] Ints:[] StructsPtr:[] Structs:[]} SomeBoolPtr:<nil> SomeStringPtr:<nil> SomeIntPtr:<nil> SomeStructPtr:<nil>} {SomeStr:hello SomeInt:123 SomeBool:false SomeStruct:{Strings:[] Ints:[] StructsPtr:[] Structs:[]} SomeBoolPtr:<nil> SomeStringPtr:<nil> SomeIntPtr:<nil> SomeStructPtr:<nil>}]",
		},
		{
			name: "nil pointer",
			fv: &FieldValue{
				Value: reflect.ValueOf(nilPtr),
				Kind:  reflect.Pointer,
			},
			expectedString: "<empty>",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.fv.ValueAsString())
		})
	}
}

func TestFieldValue_ListValuesAsString(t *testing.T) {
	for _, tc := range []struct {
		name            string
		fv              *FieldValue
		expectedStrings []string
	}{
		{
			name: "ints",
			fv: &FieldValue{
				Value:    reflect.ValueOf([]int{1, 2, 3}),
				ListKind: reflect.Int,
			},
			expectedStrings: []string{"1", "2", "3"},
		},
		{
			name: "strings",
			fv: &FieldValue{
				Value:    reflect.ValueOf([]string{"1", "2", "3"}),
				ListKind: reflect.String,
			},
			expectedStrings: []string{"1", "2", "3"},
		},
		{
			name: "bools",
			fv: &FieldValue{
				Value:    reflect.ValueOf([]bool{true, true, false}),
				ListKind: reflect.Bool,
			},
			expectedStrings: []string{"true", "true", "false"},
		},
		{
			name: "structs",
			fv: &FieldValue{
				Value: reflect.ValueOf([]randomStruct{
					{SomeStr: "flower", SomeInt: 321},
					{SomeStr: "hello", SomeInt: 123},
				}),
				ListKind: reflect.Struct,
			},
			expectedStrings: []string{
				"{SomeStr:flower SomeInt:321 SomeBool:false SomeStruct:{Strings:[] Ints:[] StructsPtr:[] Structs:[]} SomeBoolPtr:<nil> SomeStringPtr:<nil> SomeIntPtr:<nil> SomeStructPtr:<nil>}",
				"{SomeStr:hello SomeInt:123 SomeBool:false SomeStruct:{Strings:[] Ints:[] StructsPtr:[] Structs:[]} SomeBoolPtr:<nil> SomeStringPtr:<nil> SomeIntPtr:<nil> SomeStructPtr:<nil>}",
			},
		},
		{
			name: "pointers",
			fv: &FieldValue{
				Value: reflect.ValueOf([]*randomStruct{
					{SomeStr: "flower", SomeInt: 321},
					{SomeStr: "hello", SomeInt: 123},
				}),
				ListKind: reflect.Pointer,
			},
			expectedStrings: []string{
				"{SomeStr:flower SomeInt:321 SomeBool:false SomeStruct:{Strings:[] Ints:[] StructsPtr:[] Structs:[]} SomeBoolPtr:<nil> SomeStringPtr:<nil> SomeIntPtr:<nil> SomeStructPtr:<nil>}",
				"{SomeStr:hello SomeInt:123 SomeBool:false SomeStruct:{Strings:[] Ints:[] StructsPtr:[] Structs:[]} SomeBoolPtr:<nil> SomeStringPtr:<nil> SomeIntPtr:<nil> SomeStructPtr:<nil>}",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.ElementsMatch(t, tc.expectedStrings, tc.fv.ListValuesAsString())
		})
	}
}

func TestFieldValue_RemoveAtIndex(t *testing.T) {
	parent := &listStruct{
		Strings:    []string{"hello", "there"},
		Ints:       []int{4},
		StructsPtr: []*randomStruct{{SomeStr: "hello"}, {SomeInt: 123}},
	}

	for _, tc := range []struct {
		name           string
		fv             *FieldValue
		idx            int
		expectedResult any
	}{
		{
			name: "list of strings",
			fv: &FieldValue{
				Value: reflect.ValueOf(parent).Elem().Field(0),
			},
			idx: 1,
			expectedResult: &listStruct{
				Strings:    []string{"hello"},
				Ints:       []int{4},
				StructsPtr: []*randomStruct{{SomeStr: "hello"}, {SomeInt: 123}},
			},
		},
		{
			name: "remove last in list of strings",
			fv: &FieldValue{
				Value: reflect.ValueOf(parent).Elem().Field(0),
			},
			idx: 0,
			expectedResult: &listStruct{
				Strings:    []string{},
				Ints:       []int{4},
				StructsPtr: []*randomStruct{{SomeStr: "hello"}, {SomeInt: 123}},
			},
		},
		{
			name: "list of ints",
			fv: &FieldValue{
				Value: reflect.ValueOf(parent).Elem().Field(1),
			},
			idx: 0,
			expectedResult: &listStruct{
				Strings:    []string{},
				Ints:       []int{},
				StructsPtr: []*randomStruct{{SomeStr: "hello"}, {SomeInt: 123}},
			},
		},
		{
			name: "list of structs",
			fv: &FieldValue{
				Value: reflect.ValueOf(parent).Elem().Field(2),
			},
			idx: 0,
			expectedResult: &listStruct{
				Strings:    []string{},
				Ints:       []int{},
				StructsPtr: []*randomStruct{{SomeInt: 123}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.fv.RemoveAtIndex(tc.idx)
			assert.Equal(t, tc.expectedResult, parent)
		})
	}
}

func TestFieldValue_AddToStructList(t *testing.T) {
	parent := &listStruct{
		StructsPtr: []*randomStruct{{SomeStr: "hello"}},
		Structs:    []randomStruct{{SomeInt: 123}},
	}

	for _, tc := range []struct {
		name           string
		fv             *FieldValue
		newValue       any
		expectedResult any
	}{
		{
			name: "pointer",
			fv: &FieldValue{
				Value:    reflect.ValueOf(parent).Elem().Field(2),
				ListKind: reflect.Pointer,
			},
			newValue: &randomStruct{SomeInt: 321},
			expectedResult: &listStruct{
				StructsPtr: []*randomStruct{{SomeStr: "hello"}, {SomeInt: 321}},
				Structs:    []randomStruct{{SomeInt: 123}},
			},
		},
		{
			name: "not pointer",
			fv: &FieldValue{
				Value:    reflect.ValueOf(parent).Elem().Field(3),
				ListKind: reflect.Struct,
			},
			newValue: randomStruct{SomeStr: "flower"},
			expectedResult: &listStruct{
				StructsPtr: []*randomStruct{{SomeStr: "hello"}, {SomeInt: 321}},
				Structs:    []randomStruct{{SomeInt: 123}, {SomeStr: "flower"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.fv.AddToStructList(tc.newValue)
			assert.Equal(t, tc.expectedResult, parent)
		})
	}
}

func Test_DereferencePointer(t *testing.T) {
	i := 123
	parent := &randomStruct{
		SomeBoolPtr:   null.BoolFrom(true).Ptr(),
		SomeStringPtr: null.StringFrom("hello").Ptr(),
		SomeIntPtr:    &i,
		SomeStructPtr: &randomStruct{
			SomeStr: "hello",
		},
	}

	for _, tc := range []struct {
		name           string
		inputParent    any
		inputFv        *FieldValue
		expectedOutput *FieldValue
		expectPanic    bool
	}{
		{
			name:        "panics if value is snot pointer",
			inputParent: parent,
			inputFv: &FieldValue{
				Name:  "bool",
				Value: reflect.ValueOf(parent.SomeBool),
			},
			expectPanic: true,
		},
		{
			name:        "panics if field not found",
			inputParent: parent,
			inputFv: &FieldValue{
				Name: "not found",
			},
			expectPanic: true,
		},
		{
			name:        "pointer to bool",
			inputParent: parent,
			inputFv: &FieldValue{
				Name:  "bool_ptr",
				Value: reflect.ValueOf(parent.SomeBoolPtr),
				Type:  reflect.TypeOf(parent.SomeBoolPtr),
				Kind:  reflect.Pointer,
			},
			expectedOutput: &FieldValue{
				Name:  "bool_ptr",
				Value: reflect.ValueOf(parent.SomeBoolPtr).Elem(),
				Type:  reflect.TypeOf(true),
				Kind:  reflect.Bool,
			},
		},
		{
			name:        "pointer to string",
			inputParent: parent,
			inputFv: &FieldValue{
				Name:  "string_ptr",
				Value: reflect.ValueOf(parent.SomeStringPtr),
				Type:  reflect.TypeOf(parent.SomeStringPtr),
				Kind:  reflect.Pointer,
			},
			expectedOutput: &FieldValue{
				Name:  "string_ptr",
				Value: reflect.ValueOf(parent.SomeStringPtr).Elem(),
				Type:  reflect.TypeOf(""),
				Kind:  reflect.String,
			},
		},
		{
			name:        "pointer to int",
			inputParent: parent,
			inputFv: &FieldValue{
				Name:  "int_ptr",
				Value: reflect.ValueOf(parent.SomeIntPtr),
				Type:  reflect.TypeOf(parent.SomeIntPtr),
				Kind:  reflect.Pointer,
			},
			expectedOutput: &FieldValue{
				Name:  "int_ptr",
				Value: reflect.ValueOf(parent.SomeIntPtr).Elem(),
				Type:  reflect.TypeOf(1),
				Kind:  reflect.Int,
			},
		},
		{
			name:        "pointer to struct",
			inputParent: parent,
			inputFv: &FieldValue{
				Name:  "struct_ptr",
				Value: reflect.ValueOf(parent.SomeStructPtr),
				Type:  reflect.TypeOf(parent.SomeStructPtr),
				Kind:  reflect.Pointer,
			},
			expectedOutput: &FieldValue{
				Name:  "struct_ptr",
				Value: reflect.ValueOf(parent.SomeStructPtr).Elem(),
				Type:  reflect.TypeOf(randomStruct{}),
				Kind:  reflect.Struct,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectPanic {
				assert.Panics(t, func() {
					DereferencePointer(tc.inputParent, tc.inputFv)
				})
				return
			}

			gotFv := DereferencePointer(tc.inputParent, tc.inputFv)
			assert.Equal(t, tc.expectedOutput, gotFv)
		})
	}

	// see that initialization works
	in := &randomStruct{}
	gotFv := DereferencePointer(in, &FieldValue{
		Name:  "struct_ptr",
		Value: reflect.ValueOf(in.SomeStructPtr),
		Type:  reflect.TypeOf(parent.SomeStructPtr),
		Kind:  reflect.Pointer,
	})

	// can't compare values since address will be different
	assert.Equal(t, reflect.Struct, gotFv.Kind)
	assert.Equal(t, "struct_ptr", gotFv.Name)
	assert.Equal(t, gotFv.Type, reflect.TypeOf(randomStruct{}))
}

func Test_GetFieldFromList(t *testing.T) {
	for _, tc := range []struct {
		name          string
		idx           int
		list          any
		expectedValue any
	}{
		{
			name:          "int list",
			idx:           2,
			list:          []int{1, 2, 3},
			expectedValue: 3,
		},
		{
			name:          "string list",
			idx:           1,
			list:          []string{"1", "2", "3"},
			expectedValue: "2",
		},
		{
			name:          "bool list",
			idx:           0,
			list:          []bool{true, true, false},
			expectedValue: true,
		},
		{
			name:          "struct list",
			idx:           1,
			list:          []randomStruct{{SomeStr: "hello"}, {SomeInt: 3}},
			expectedValue: randomStruct{SomeInt: 3},
		},
		{
			name:          "struct ponter list",
			idx:           0,
			list:          []*randomStruct{{SomeStr: "hello"}, {SomeInt: 3}},
			expectedValue: &randomStruct{SomeStr: "hello"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gotValue := GetFieldFromList(tc.idx, tc.list)
			assert.Equal(t, tc.expectedValue, gotValue)
		})
	}
}

func Test_InitializeIfNil(t *testing.T) {
	parent := &randomStruct{}
	i := 0

	for _, tc := range []struct {
		name            string
		parent          any
		fv              *FieldValue
		expectedOutcome any
	}{
		{
			name:   "bool pointer",
			parent: parent,
			fv: &FieldValue{
				Name:  "bool_ptr",
				Value: reflect.ValueOf(parent.SomeBoolPtr),
				Type:  reflect.TypeOf(parent.SomeBoolPtr),
				Kind:  reflect.Pointer,
			},
			expectedOutcome: &randomStruct{
				SomeBoolPtr: null.BoolFrom(false).Ptr(),
			},
		},
		{
			name:   "intpointer",
			parent: parent,
			fv: &FieldValue{
				Name:  "int_ptr",
				Value: reflect.ValueOf(parent.SomeIntPtr),
				Type:  reflect.TypeOf(parent.SomeIntPtr),
				Kind:  reflect.Pointer,
			},
			expectedOutcome: &randomStruct{
				SomeBoolPtr: null.BoolFrom(false).Ptr(),
				SomeIntPtr:  &i,
			},
		},
		{
			name:   "string pointer",
			parent: parent,
			fv: &FieldValue{
				Name:  "string_ptr",
				Value: reflect.ValueOf(parent.SomeStringPtr),
				Type:  reflect.TypeOf(parent.SomeStringPtr),
				Kind:  reflect.Pointer,
			},
			expectedOutcome: &randomStruct{
				SomeBoolPtr:   null.BoolFrom(false).Ptr(),
				SomeIntPtr:    &i,
				SomeStringPtr: null.StringFrom("").Ptr(),
			},
		},
		{
			name:   "struct pointer",
			parent: parent,
			fv: &FieldValue{
				Name:  "struct_ptr",
				Value: reflect.ValueOf(parent.SomeStructPtr),
				Type:  reflect.TypeOf(parent.SomeStructPtr),
				Kind:  reflect.Pointer,
			},
			expectedOutcome: &randomStruct{
				SomeBoolPtr:   null.BoolFrom(false).Ptr(),
				SomeIntPtr:    &i,
				SomeStringPtr: null.StringFrom("").Ptr(),
				SomeStructPtr: &randomStruct{},
			},
		},
		{
			name:   "not pointer is ok",
			parent: parent,
			fv: &FieldValue{
				Name:  "string",
				Value: reflect.ValueOf(parent.SomeStr),
				Type:  reflect.TypeOf(parent.SomeStr),
				Kind:  reflect.Pointer,
			},
			expectedOutcome: &randomStruct{
				SomeBoolPtr:   null.BoolFrom(false).Ptr(),
				SomeIntPtr:    &i,
				SomeStringPtr: null.StringFrom("").Ptr(),
				SomeStructPtr: &randomStruct{},
				SomeStr:       "",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			InitializeIfNil(tc.parent, tc.fv)
			assert.Equal(t, tc.parent, tc.expectedOutcome)
		})
	}
}

func Test_GetFromTag(t *testing.T) {
	for _, tc := range []struct {
		name          string
		parent        any
		tag           string
		expectedValue any
	}{
		{
			name:          "int",
			parent:        randomStruct{SomeInt: 3},
			tag:           "int",
			expectedValue: 3,
		},
		{
			name:          "string",
			parent:        &randomStruct{SomeStr: "hello"},
			tag:           "string",
			expectedValue: "hello",
		},
		{
			name:          "bool",
			parent:        &randomStruct{SomeBool: true},
			tag:           "bool",
			expectedValue: true,
		},
		{
			name:          "bool pointer",
			parent:        &randomStruct{SomeBoolPtr: null.BoolFrom(true).Ptr()},
			tag:           "bool_ptr",
			expectedValue: null.BoolFrom(true).Ptr(),
		},
		{
			name:          "struct pointer",
			parent:        &randomStruct{SomeStructPtr: &randomStruct{SomeStr: "hello"}},
			tag:           "struct_ptr",
			expectedValue: &randomStruct{SomeStr: "hello"},
		},
		{
			name:          "struct",
			parent:        &randomStruct{SomeStruct: listStruct{Strings: []string{"hello"}}},
			tag:           "struct",
			expectedValue: listStruct{Strings: []string{"hello"}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gotValue := GetFromTag(tc.parent, tc.tag)
			assert.Equal(t, tc.expectedValue, gotValue)
		})
	}
}

func Test_SetValue(t *testing.T) {
	for _, tc := range []struct {
		name            string
		parent          any
		tag             string
		value           any
		expectedOutcome any
	}{
		{
			name:            "string",
			parent:          &randomStruct{},
			tag:             "string",
			value:           "hello",
			expectedOutcome: &randomStruct{SomeStr: "hello"},
		},
		{
			name:            "int",
			parent:          &randomStruct{},
			tag:             "int",
			value:           3,
			expectedOutcome: &randomStruct{SomeInt: 3},
		},
		{
			name:            "bool",
			parent:          &randomStruct{},
			tag:             "bool",
			value:           true,
			expectedOutcome: &randomStruct{SomeBool: true},
		},
		{
			name:            "bool pointer",
			parent:          &randomStruct{},
			tag:             "bool_ptr",
			value:           true,
			expectedOutcome: &randomStruct{SomeBoolPtr: null.BoolFrom(true).Ptr()},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			SetValue(tc.parent, tc.tag, tc.value)
			assert.Equal(t, tc.expectedOutcome, tc.parent)
		})
	}
}

func Test_GetFieldValues(t *testing.T) {
	parent := &randomStruct{
		SomeStr:       "hello",
		SomeBool:      true,
		SomeStringPtr: null.StringFrom("hello").Ptr(),
	}
	fvs := GetFieldValues(parent)

	assert.ElementsMatch(t, fvs, []*FieldValue{
		{
			Name:  "string",
			Value: reflect.ValueOf(parent).Elem().Field(0),
			Type:  reflect.TypeOf(parent.SomeStr),
			Kind:  reflect.String,
		}, {
			Name:  "int",
			Value: reflect.ValueOf(parent).Elem().Field(1),
			Type:  reflect.TypeOf(parent.SomeInt),
			Kind:  reflect.Int,
		}, {
			Name:  "bool",
			Value: reflect.ValueOf(parent).Elem().Field(2),
			Type:  reflect.TypeOf(parent.SomeBool),
			Kind:  reflect.Bool,
		}, {
			Name:  "struct",
			Value: reflect.ValueOf(parent).Elem().Field(3),
			Type:  reflect.TypeOf(parent.SomeStruct),
			Kind:  reflect.Struct,
		}, {
			Name:  "bool_ptr",
			Value: reflect.ValueOf(parent).Elem().Field(4),
			Type:  reflect.TypeOf(parent.SomeBoolPtr),
			Kind:  reflect.Pointer,
		}, {
			Name:  "string_ptr",
			Value: reflect.ValueOf(parent).Elem().Field(5),
			Type:  reflect.TypeOf(parent.SomeStringPtr),
			Kind:  reflect.Pointer,
		}, {
			Name:  "int_ptr",
			Value: reflect.ValueOf(parent).Elem().Field(6),
			Type:  reflect.TypeOf(parent.SomeIntPtr),
			Kind:  reflect.Pointer,
		}, {
			Name:  "struct_ptr",
			Value: reflect.ValueOf(parent).Elem().Field(7),
			Type:  reflect.TypeOf(parent.SomeStructPtr),
			Kind:  reflect.Pointer,
		},
	})

	listParent := []string{"1", "2", "3"}
	fvs = GetFieldValues(listParent)

	assert.ElementsMatch(t, fvs, []*FieldValue{
		{
			Value: reflect.ValueOf(listParent).Index(0),
			Type:  reflect.TypeOf(listParent),
			Kind:  reflect.String,
		}, {
			Value: reflect.ValueOf(listParent).Index(1),
			Type:  reflect.TypeOf(listParent),
			Kind:  reflect.String,
		}, {
			Value: reflect.ValueOf(listParent).Index(2),
			Type:  reflect.TypeOf(listParent),
			Kind:  reflect.String,
		},
	})
}
