# serr [![Go Reference](https://pkg.go.dev/badge/github.com/daaku/serr.svg)](https://pkg.go.dev/github.com/daaku/serr)

Package serr provides stackful errors and nothing else.

A traditional method to augment errors is stack traces. The `serr` package
allows for programmers to add stack traces to errors without destroying the
original error value.

It provides pretty formatting and structured access to the stack traces.
Structured access is in terms of the `[]uintptr` collected via
`runtime.Callers`. Pretty formatting is provided in terms of support for the
`%+v` style format modifier.

This is the entire extent of this library. Anything else is out of scope.
