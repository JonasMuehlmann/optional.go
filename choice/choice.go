package choice

type Choice[T, U any] struct {
	either   T
	or       U
	isEither bool
}

func Either[T, U any](value T) Choice[T, U] {
	return Choice[T, U]{either: value, isEither: true}
}

func Or[T, U any](value U) Choice[T, U] {
	return Choice[T, U]{or: value, isEither: false}
}

func (choice Choice[T, U]) IsEither() bool {
	return choice.isEither
}

func (choice Choice[T, U]) IsOr() bool {
	return !choice.isEither
}

func (choice Choice[T, U]) AssertEither() {
	if choice.IsOr() {
		panic("Either is not set")
	}
}

func (choice Choice[T, U]) AssertOr() {
	if choice.IsEither() {
		panic("Or is not set")
	}

}

func (choice Choice[T, U]) MustGetEither() T {
	choice.AssertEither()

	return choice.either
}

func (choice Choice[T, U]) GetEitherElseOr(or U) Choice[T, U] {
	if choice.IsEither() {
		return choice
	}

	return Or[T, U](or)
}
func (choice Choice[T, U]) GetEitherElseDefault() Choice[T, U] {
	if choice.IsEither() {
		return choice
	}

	var t T

	return Either[T, U](t)
}

func (choice Choice[T, U]) GetEitherElseFrom(from func() T) Choice[T, U] {
	if choice.IsEither() {
		return choice
	}

	return Either[T, U](from())
}

func (choice Choice[T, U]) MustGetOr() U {
	choice.AssertOr()

	return choice.or
}

func (choice Choice[T, U]) GetOrElseEither(either T) Choice[T, U] {
	if choice.IsOr() {
		return choice
	}

	return Either[T, U](either)
}

func (choice Choice[T, U]) GetOrElseDefault() Choice[T, U] {
	if choice.IsOr() {
		return choice
	}

	var u U

	return Or[T, U](u)
}

func (choice Choice[T, U]) GetOrElseFrom(from func() U) Choice[T, U] {
	if choice.IsOr() {
		return choice
	}

	return Or[T, U](from())
}

func (choice Choice[T, U]) Match(either func(T), or func(U)) {
	if choice.IsEither() {
		either(choice.MustGetEither())
	} else {
		or(choice.MustGetOr())
	}
}
