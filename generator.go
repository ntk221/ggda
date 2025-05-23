// Package ggda provides utilities for generating test data from structs
// ggda - Go Generate Data Automatically
package ggda

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Generator[T any] struct {
	defaults map[string]interface{}
	customs  map[string]func(index int) interface{}
}

func New[T any]() *Generator[T] {
	return &Generator[T]{
		defaults: make(map[string]interface{}),
		customs:  make(map[string]func(index int) interface{}),
	}
}

// Generate creates a slice of structs with the specified count
func (g *Generator[T]) Generate(count int) []T {
	result := make([]T, count)
	for i := 0; i < count; i++ {
		var elem T
		v := reflect.ValueOf(&elem).Elem()
		g.fillStruct(v, i)
		result[i] = elem
	}
	return result
}

func (g *Generator[T]) GenerateOne() T {
	var elem T
	v := reflect.ValueOf(&elem).Elem()
	g.fillStruct(v, 0)
	return elem
}

// SetDefaults sets default values for specific fields
func (g *Generator[T]) SetDefaults(fieldName string, value interface{}) *Generator[T] {
	g.defaults[fieldName] = value
	return g
}

// SetCustom sets a custom generator function for a specific field
func (g *Generator[T]) SetCustom(fieldName string, fn func(index int) interface{}) *Generator[T] {
	g.customs[fieldName] = fn
	return g
}

// fillStruct fills a struct with test data
// This method is public so that it can be used by Builder
func (g *Generator[T]) fillStruct(v reflect.Value, index int) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Check for custom generator
		if customFn, ok := g.customs[fieldName]; ok {
			field.Set(reflect.ValueOf(customFn(index)))
			continue
		}

		// Check for default value
		if defaultVal, ok := g.defaults[fieldName]; ok {
			field.Set(reflect.ValueOf(defaultVal))
			continue
		}

		// Auto-generate based on type
		g.autoFill(field, fieldType, index)
	}
}

// autoFill automatically fills a field based on its type
func (g *Generator[T]) autoFill(field reflect.Value, fieldType reflect.StructField, index int) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(fmt.Sprintf("%s_%d", strings.ToLower(fieldType.Name), index+1))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(int64(index + 1))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.SetUint(uint64(index + 1))
	case reflect.Float32, reflect.Float64:
		field.SetFloat(float64(index+1) * 1.1)
	case reflect.Bool:
		field.SetBool(index%2 == 0)
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			field.Set(reflect.ValueOf(time.Now()))
		}
	}
}

// GenerateSlice creates a slice of structs with the specified count
func GenerateSlice[T any](count int) []T {
	result := make([]T, count)
	var zero T
	v := reflect.ValueOf(zero)

	// Check if T is a struct type
	if v.Kind() == reflect.Struct {
		// For struct types, use Generator
		gen := New[T]()
		return gen.Generate(count)
	}

	// For primitive types, generate directly
	for i := 0; i < count; i++ {
		result[i] = generatePrimitive[T](i)
	}
	return result
}

// GenerateSliceWith creates a slice of structs with custom modification
func GenerateSliceWith[T any](count int, modifier func(item *T, index int)) []T {
	result := make([]T, count)
	var zero T
	v := reflect.ValueOf(zero)

	// Check if T is a struct type
	if v.Kind() == reflect.Struct {
		// For struct types, use Generator
		gen := New[T]()
		for i := 0; i < count; i++ {
			var elem T
			v := reflect.ValueOf(&elem).Elem()
			gen.fillStruct(v, i)
			if modifier != nil {
				modifier(&elem, i)
			}
			result[i] = elem
		}
	} else {
		// For primitive types
		for i := 0; i < count; i++ {
			elem := generatePrimitive[T](i)
			if modifier != nil {
				modifier(&elem, i)
			}
			result[i] = elem
		}
	}
	return result
}

// generatePrimitive generates a primitive value
func generatePrimitive[T any](index int) T {
	var result T
	v := reflect.ValueOf(&result).Elem()

	switch v.Kind() {
	case reflect.String:
		v.SetString(fmt.Sprintf("text_%d", index+1))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(index + 1))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(index + 1))
	case reflect.Float32, reflect.Float64:
		v.SetFloat(float64(index+1) * 1.1)
	case reflect.Bool:
		v.SetBool(index%2 == 0)
	}

	return result
}
