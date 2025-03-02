package paths

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ExampleStruct struct {
	PtrField1   *string
	PtrField2   *int
	NonPtrField string
}

func TestCheckPointersNotNil(t *testing.T) {
	t.Run("AllPointersSet", func(t *testing.T) {
		s1 := "hello"
		i1 := 42
		ex := ExampleStruct{
			PtrField1:   &s1,
			PtrField2:   &i1,
			NonPtrField: "non-pointer",
		}
		err := CheckPointersNotNil(ex)
		assert.NoError(t, err)
	})

	t.Run("OnePointerIsNil", func(t *testing.T) {
		s1 := "not nil"
		var i1 *int
		ex := ExampleStruct{
			PtrField1:   &s1,
			PtrField2:   i1,
			NonPtrField: "non-pointer",
		}
		err := CheckPointersNotNil(ex)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PtrField2 is nil")
	})

	t.Run("PointerToStruct_AllSet", func(t *testing.T) {
		s1 := "some string"
		i1 := 123
		ex := &ExampleStruct{
			PtrField1:   &s1,
			PtrField2:   &i1,
			NonPtrField: "another non-pointer",
		}
		err := CheckPointersNotNil(ex)
		assert.NoError(t, err)
	})

	t.Run("PointerToStruct_NilField", func(t *testing.T) {
		s1 := "some string"
		ex := &ExampleStruct{
			PtrField1: &s1,
			PtrField2: nil,
		}
		err := CheckPointersNotNil(ex)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "PtrField2 is nil")
	})

	t.Run("NonStructInput", func(t *testing.T) {
		notAStruct := 123
		err := CheckPointersNotNil(notAStruct)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrNotAStruct))
	})

	t.Run("NoPointerFields", func(t *testing.T) {
		type NoPointers struct {
			A int
			B string
		}
		ex := NoPointers{A: 1, B: "ok"}
		err := CheckPointersNotNil(ex)
		assert.NoError(t, err)
	})

	t.Run("OptionalFieldIsNil", func(t *testing.T) {
		type OptionalStruct struct {
			PtrFieldOptional *string `svmkit:"optional"`
			PtrFieldRequired *string
		}
		s1 := "I am required"
		ex := OptionalStruct{
			PtrFieldOptional: nil,
			PtrFieldRequired: &s1,
		}
		err := CheckPointersNotNil(ex)
		assert.NoError(t, err)
	})
}
