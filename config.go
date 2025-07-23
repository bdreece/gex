package gex

type (
	// A Config provides additional configuration parameters for the lexer
	// state machine in the form of a struct.
	Config[T any] struct {
		// The input text (required).
		Input string
		// The initial state (required).
		Init State[T]

		// The name of the lexer (defaults to "gex").
		Name string
		// The buffer size of the lexer token stream (defaults to 2).
		Capacity int
		// The sentinel `rune` value used to demarcate the end of the
		// token stream (defaults to `*new(rune)`).
		EOF rune
		// The sentinel `T` value used when emitting a synthetic error
		// token (defaults to `*new(T)`).
		Error T
	}

	// An Option provides an internal mechanism for mutating a [Config].
	Option[T any] interface {
		apply(c *Config[T])
	}

	option[T any] func(c *Config[T])
)

func (fn option[T]) apply(config *Config[T]) { fn(config) }

func (c Config[T]) apply(config *Config[T]) {
	config.Name = c.Name
	config.Capacity = c.Capacity
	config.EOF = c.EOF
	config.Error = c.Error
}

// DefaultConfig returns a [Config] object populated with the default values.
func DefaultConfig[T any]() *Config[T] {
	return &Config[T]{
		Name:     "gex",
		Capacity: 2,
	}
}

// WithName returns an option configuring the name of the [Lexer].
func WithName[T any](name string) Option[T] {
	return option[T](func(config *Config[T]) {
		config.Name = name
	})
}

// WithCapacity returns an option configuring the capacity of the [Lexer] token
// stream.
func WithCapacity[T any](capacity int) Option[T] {
	return option[T](func(config *Config[T]) {
		config.Capacity = capacity
	})
}

// WithEOF returns an option configuring the sentinel EOF rune associated with
// the lexer. See [Config] for more details.
func WithEOF[T any](eof rune) Option[T] {
	return option[T](func(config *Config[T]) {
		config.EOF = eof
	})
}

// WithErrorType returns an option configuring the sentinel error value
// associated with the lexer. See [Config] for more details.
func WithErrorType[T any](typ T) Option[T] {
	return option[T](func(config *Config[T]) {
		config.Error = typ
	})
}
