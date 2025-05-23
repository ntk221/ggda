package ggda

import (
	"reflect"
)

type Builder[T any] struct {
	gen       *Generator[T]
	modifiers []func(v *T, index int)
}

func Build[T any]() *Builder[T] {
	return &Builder[T]{
		gen:       New[T](),
		modifiers: make([]func(v *T, index int), 0),
	}
}

// With sets values using a modifier function
func (b *Builder[T]) With(modifier func(v *T, index int)) *Builder[T] {
	b.modifiers = append(b.modifiers, modifier)
	return b
}

// WithDefaults sets default values using a struct
func (b *Builder[T]) WithDefaults(defaults T) *Builder[T] {
	v := reflect.ValueOf(defaults)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.IsZero() {
			b.gen.defaults[fieldType.Name] = field.Interface()
		}
	}

	return b
}

// Generate creates a slice of structs
func (b *Builder[T]) Generate(count int) []T {
	result := make([]T, count)
	for i := 0; i < count; i++ {
		result[i] = b.generateSingle(i)
	}
	return result
}

// GenerateOne creates a single struct
func (b *Builder[T]) GenerateOne() T {
	return b.generateSingle(0)
}

// generateSingle generates a single struct at the given index
func (b *Builder[T]) generateSingle(index int) T {
	var elem T

	// fill struct with defaults and auto-generation
	v := reflect.ValueOf(&elem).Elem()
	b.gen.fillStruct(v, index)

	// apply modifiers
	for _, modifier := range b.modifiers {
		modifier(&elem, index)
	}

	return elem
}
