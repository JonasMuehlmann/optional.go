# optional.go
optional.go is a simple and generic wrapper around a value and a presence flag inspired by C++'s `std::optional`.

## Installation

```go get github.com/JonasMuehlmann/optional.go```

## How to use

This project exports the package `optional`, whose only exported entity is the generic `Optional` struct.

The `Optional` type implements several convenience methods and standard library interfaces( e.g json and SQL conversions).

To see the implemented methods and interfaces, refer to the documentation available at https://pkg.go.dev/github.com/JonasMuehlmann/optional.go.
