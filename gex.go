// Package gex provides a generic framework for lexical analysis of UTF-8
// text in Go, based on the implementation discussed in
// ["Lexical Scanning in Go"] by [robpike].
//
// ["Lexical Scanning in Go"]: https://go.dev/talks/2011/lex.slide#1
// [robpike]: https://github.com/robpike
package gex

// Starts analyzing the input text with the provided init state in a new goroutine,
// and returns a channel to receive the emitted [Token] objects.
//
// Additional options may be provided as variadic arguments.
func Start[T any](input string, init State[T], opts ...Option[T]) <-chan Token[T] {
	config := DefaultConfig[T]()
	config.Input = input
	config.Init = init
	for _, opt := range opts {
		opt.apply(config)
	}

	return StartWithConfig(config)
}

// Starts analyzing the input text using parameters provided through the [Config]
// object.
//
// See [Start] for more details.
func StartWithConfig[T any](config *Config[T]) <-chan Token[T] {
	lexer := Lexer[T]{
		name:   config.Name,
		input:  config.Input,
		eof:    config.EOF,
		error:  config.Error,
		tokens: make(chan Token[T], config.Capacity),
	}

	go lexer.run(config.Init)
	return lexer.tokens
}
