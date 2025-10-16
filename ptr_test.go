package spec

import (
	"reflect"
	"testing"
)

func TestPtr(t *testing.T) {
	testCases := []struct {
		name string
		val  any
	}{
		{name: "int", val: 42},
		{name: "string", val: "test"},
		{name: "bool", val: true},
		{name: "nil", val: nil},
		{name: "slice", val: []int{1, 2, 3}},
		{name: "map", val: map[string]int{"a": 1, "b": 2}},
		{name: "struct", val: struct{ Field string }{Field: "value"}},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ptr := Ptr(tt.val)
			if ptr == nil {
				t.Fatal("Ptr returned nil")
			}
			if !reflect.DeepEqual(*ptr, tt.val) {
				t.Fatalf("expected %v, got %v", tt.val, *ptr)
			}
		})
	}
}
