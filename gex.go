// Package gex provides a generic framework for lexical analysis of UTF-8
// text in Go, based on the implementation discussed in
// ["Lexical Scanning in Go"] by [robpike].
//
// ["Lexical Scanning in Go"]: https://go.dev/talks/2011/lex.slide#1
// [robpike]: https://github.com/robpike
package gex

// Run creates a new [Lexer] with the provided input text and init state, then
// launches a goroutine that advances the state machine before finally closing
// the token stream.
//
// Additional options may be provided as variadic arguments.
func Run[T any](input string, init State[T], opts ...Option[T]) <-chan Token[T] {
	config := DefaultConfig[T]()
	config.Input = input
	config.Init = init
	for _, opt := range opts {
		opt.apply(config)
	}

	return RunWithConfig(config)
}

// RunWithConfig creates and starts the [Lexer] using parameters from the
// provided [Config] object.
//
// See [Run] for more details.
func RunWithConfig[T any](config *Config[T]) <-chan Token[T] {
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
