package result

import "github.com/JonasMuehlmann/optional.go/choice"

type Result[T any] struct {
	choice choice.Choice[T, error]
}

func Ok[T any](value T) Result[T] {
	return Result[T]{choice: choice.Either[T, error](value)}
}

func Err[T any](value error) Result[T] {
	return Result[T]{choice: choice.Or[T, error](value)}
}

func ToResult[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}

	return Ok[T](value)
}

func ToTuple[T any](result Result[T]) (T, error) {
	if result.IsOk() {
		return result.MustGetOk(), nil
	}

	var t T
	return t, result.MustGetErr()
}

func (result Result[T]) ToEither() choice.Choice[T, error] {
	return result.choice
}

func (result Result[T]) IsOk() bool {
	return result.choice.IsEither()
}

func (result Result[T]) IsErr() bool {
	return result.choice.IsOr()
}

func (result Result[T]) AssertOk() {
	result.choice.AssertEither()
}

func (result Result[T]) AssertErr() {
	result.choice.AssertOr()
}

func (result Result[T]) MustGetOk() T {
	result.choice.AssertEither()

	return result.choice.MustGetEither()
}

func (result Result[T]) MustGetErr() error {
	result.choice.AssertOr()

	return result.choice.MustGetOr()
}

func (result Result[T]) Match(okHandler func(T), errHandler func(error)) {
	result.choice.Match(okHandler, errHandler)
}

func (result Result[T]) Try(f func() (T, error)) Result[T] {
	if result.IsErr() {
		return result
	}

	value, err := f()
	return ToResult(value, err)
}

func (result Result[T]) TryT(f func() (T, error)) Result[T] {
	if result.IsErr() {
		return result
	}

	value, err := f()
	return ToResult(value, err)
}
