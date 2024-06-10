// Copyright © 2021-2022 Jonas Muehlmann
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// optional.go is a simple wrapper around a value and a presence flag
package optional

import (
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"github.com/JonasMuehlmann/optional.go/v3/result"

	"github.com/gocarina/gocsv"
)

// Optional holds a Wrappe and a flag indicating the existence or abscence of the Wrapee.
type Optional[T any] struct {
	wrappee  T    `json:"wrapee" db:"wrapee"`
	hasValue bool `json:"has_value" db:"has_value"`
}

// Some creates an optional holding the specified wrappee and sets the hasValue flag to true.
func Some[T any](wrappee T) Optional[T] {
	return Optional[T]{wrappee: wrappee, hasValue: true}
}

// None creates an optional holding the specified wrappee and sets the hasValue flag to false.
func None[T any]() Optional[T] {
	return Optional[T]{hasValue: false}
}

func FromT[T any](t T, ok bool) Optional[T] {
	if ok {
		return Some(t)
	}

	return None[T]()
}

func (optional Optional[T]) MustGet() T {
	if optional.IsNone() {
		panic("optional is empty")
	}
	return optional.wrappee

}

// GetElseAlt returns the Wrapee if the hasValue flag is true, otherwise alternative is returned.
func (optional Optional[T]) GetElseAlt(alternative T) T {
	if optional.IsSome() {
		return optional.wrappee
	}

	return alternative
}

// GetOrDefualt returns the Wrapee if the hasValue flag is true, otherwise T's default value is returned.
func (optional Optional[T]) GetElseDefault() T {
	if optional.IsSome() {
		return optional.wrappee
	}

	var t T

	return t
}

// GetElseFrom returns the Wrapee if the hasValue flag is true, otherwise the result of alternative() is returned.
func (optional Optional[T]) GetElseFrom(alternative func() T) Optional[T] {
	if optional.IsSome() {
		return Some(alternative())
	}
	return None[T]()
}

func (optional Optional[T]) TryGetElseFromT(alternative func() (T, error)) result.Result[T] {
	if optional.IsSome() {
		return result.Ok(optional.MustGet())
	}

	return result.FromTuple(alternative())

}

func (optional Optional[T]) TryGetElseFrom(alternative func() result.Result[T]) result.Result[T] {
	if optional.IsSome() {
		return result.Ok(optional.MustGet())
	}

	return alternative()

}

// GetTransformedOrEmpty returns transformer(wrappee) if hasValue == true, otherwise returns self.
func (optional Optional[T]) GetTransformedElseNone(transformer func(T) T) Optional[T] {
	if optional.IsSome() {
		return Some(transformer(optional.wrappee))
	}

	return None[T]()
}

func (optional Optional[T]) TryDo(f func()) Optional[T] {
	if optional.IsNone() {
		return optional
	}

	f()

	return optional
}

// Match executes someHandler if hasValue is true and noneHandler otherwise.
func (optional Optional[T]) Match(someHandler func(T), noneHandler func(T)) {
	if optional.IsSome() {
		someHandler(optional.wrappee)
	} else {
		noneHandler(optional.wrappee)
	}
}

// Set sets a wrappee to val value and the hasValue flag to true.
func (optional *Optional[T]) Set(val T) {
	optional.wrappee = val
	optional.hasValue = true
}

// Unset sets the hasValue flag to false and returns the wrappee.
func (optional *Optional[T]) Unset() T {
	optional.hasValue = false

	return optional.wrappee
}

func (optional Optional[T]) IsSome() bool {
	return optional.hasValue
}

func (optional Optional[T]) IsNone() bool {
	return !optional.hasValue
}

func (optional Optional[T]) AssertSome() {
	if optional.IsNone() {
		panic("optional is empty")
	}
}

func (optional Optional[T]) AssertNone() {
	if optional.IsSome() {
		panic("optional is not empty")
	}
}

func (optional Optional[T]) MarshalJSON() ([]byte, error) {
	if !optional.hasValue {
		return []byte("null"), nil
	}

	return json.Marshal(optional.wrappee)
}

func (optional *Optional[T]) UnmarshalJSON(input []byte) error {
	if string(input) == "null" {
		return nil
	}

	// Try to parse key with value of equal type as "Wrapee"
	err := json.Unmarshal(input, &optional.wrappee)
	if err == nil {
		optional.hasValue = true

		return nil
	}

	alias := struct {
		Wrappee  T    `json:"wrapee" db:"wrapee"`
		HasValue bool `json:"has_value" db:"has_value"`
	}{}
	// Try to parse object with keys "wrapee" and "has_value"
	err = json.Unmarshal(input, &alias)
	if err != nil {
		optional.hasValue = false

		return err
	}

	optional.hasValue = alias.HasValue
	optional.wrappee = alias.Wrappee

	return nil
}

// Scan implements the database.sql.Scanner interface.
func (optional *Optional[T]) Scan(value any) error {
	if value == nil {
		var zero T

		optional.hasValue = false
		optional.wrappee = zero

		return nil
	}

	switch t := value.(type) {
	case T:
		optional.hasValue = true
		optional.wrappee = t

		return nil
	}

	return fmt.Errorf("failed to scan value of type %T into optional of type %T", value, optional.wrappee)
}

// Value implements the database.sql.driver.Valuer interface.
func (optional Optional[T]) Value() (driver.Value, error) {
	if !optional.hasValue {
		return nil, nil
	}

	if driver.IsValue(optional.wrappee) {
		return optional.wrappee, nil
	}

	return nil, nil
}

// String implements the fmt.Stringer interface.
func (optional Optional[T]) String() string {
	if optional.hasValue {
		return fmt.Sprint(optional.wrappee)
	}

	return "empty optional"
}

// MarshalText implements the encoding.TextMarshaller interface.
func (optional *Optional[T]) MarshalText() (text []byte, err error) {
	if !optional.hasValue {
		return []byte{}, nil
	}

	textMarshaler, ok := any(optional.wrappee).(encoding.TextMarshaler)
	if !ok {
		return []byte{}, fmt.Errorf("failed to marshal value of type %T, TextMarshaler interface not supported", optional.wrappee)
	}

	return textMarshaler.MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaller interface.
func (optional *Optional[T]) UnmarshalText(text []byte) error {
	textUnarshaler, ok := any(optional.wrappee).(encoding.TextUnmarshaler)
	if !ok {
		return fmt.Errorf(`failed to unmarshal value "%v" into optional of type %T, wrappee does not support TextUnmarshaler interface`, text, optional.wrappee)
	}

	return textUnarshaler.UnmarshalText(text)
}

// UnmarshalCSV implements the gocarina/gocsv.TypeUnmarshaller interface.
func (optional *Optional[T]) UnmarshalCSV(val string) error {
	if val == "" {
		optional.hasValue = false

		return nil
	}

	// Dirty Hack
	var temp []struct{ Foo T }

	err := gocsv.UnmarshalString("Foo\n"+val, &temp)
	if err != nil {
		return err
	}

	optional.Set(temp[0].Foo)

	return nil
}

// MarshalCSV implements the gocarina/gocsv.TypeMarshaller interface.
func (optional Optional[T]) MarshalCSV() (string, error) {
	// csvMarshaler, ok := any(optional.wrappee).(gocsv.TypeMarshaller)
	// if !ok {
	// 	return "", fmt.Errorf(`Failed to marshal value of type %T, wrappee does not support TypeMarshaller interface`, optional.wrappee)
	// }

	// return csvMarshaler.MarshalCSV()

	// Dirty Hack
	if !optional.hasValue {
		return "", nil
	}

	temp := []struct{ Foo T }{{optional.wrappee}}

	out, err := gocsv.MarshalString(temp)
	if err != nil {
		return "", err
	}

	// If it's ghetto and it works, it is not ghetto.
	out = out[4:]
	out = out[:len(out)-1]

	return out, err
}
