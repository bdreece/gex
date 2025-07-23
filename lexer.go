package gex

import (
	"fmt"
	"iter"
	"strings"
	"unicode/utf8"
)

type (
	// A Token is a basic, abstract unit of lexical meaning.
	//
	// Token is generic over type T (typically an enum), representing the
	// token type.
	Token[T any] struct {
		// The token type.
		Type T
		// The token value.
		Value string
	}

	// Lexer provides utility methods for emitting tokens and directing the state
	// of iteration.
	Lexer[T any] struct {
		name, input       string
		start, pos, width int
		eof               rune
		error             T
		tokens            chan Token[T]
		unread            bool
	}

	// A State is a node in the lexer state machine.
	//
	// Implementors of this function contract can take advantage of utility
	// functions provided through the [Lexer] object to emit tokens and
	// direct iteration over the input text.
	State[T any] func(l *Lexer[T]) (next State[T])
)

// String implements [fmt.Stringer]
func (t Token[T]) String() string { return t.Value }

// Name returns the name of the lexer.
func (l Lexer[T]) Name() string { return l.name }

// Input returns the original input text.
func (l Lexer[T]) Input() string { return l.input }

// EOF returns the sentinel EOF rune associated with the lexer.
func (l Lexer[T]) EOF() rune { return l.eof }

// Start returns the start position of the current token.
func (l Lexer[T]) Start() int { return l.start }

// Pos returns the position of the next rune.
func (l Lexer[T]) Pos() int { return l.pos }

// Current returns a string slice representing the current token.
func (l Lexer[T]) Current() string { return l.input[l.start:l.pos] }

// Width returns the width of the last scanned rune.
func (l Lexer[T]) Width() int { return l.width }

// Accept advances the iterator position if the next rune is in the valid set.
func (l *Lexer[T]) Accept(valid string) bool {
	if !strings.ContainsRune(valid, l.Next()) {
		l.Back()
		return false
	}

	return true
}

// AcceptRun advances the iterator position until it encounters a rune not found
// in the valid set.
func (l *Lexer[T]) AcceptRun(valid string) {
	for {
		if !l.Accept(valid) {
			break
		}
	}
}

// Back decrements the current iterator position. Can only be run once per call
// to Next.
func (l *Lexer[T]) Back() (ok bool) {
	if l.unread {
		return false
	}

	l.pos -= l.width
	l.unread = true
	return true
}

// Next returns the next rune in the input text, advancing the iterator
// position.
func (l *Lexer[T]) Next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return l.eof
	}

	rune, n := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = n
	l.pos += n
	l.unread = false

	return rune
}

// Runes returns an iterator over the runes in the sequence, until EOF is reached.
func (l *Lexer[T]) Runes() iter.Seq[rune] {
	return func(yield func(rune) bool) {
		for r := l.Next(); r != l.eof; {
			if !yield(r) {
				return
			}

			r = l.Next()
		}
	}
}

// Skip advances the start cursor to the current iterator position.
func (l *Lexer[T]) Skip() {
	l.start = l.pos
}

// Peek returns the next rune in the input text without advancing the iterator
// position.
func (l *Lexer[T]) Peek() rune {
	rune := l.Next()
	l.Back()

	return rune
}

// Emit emits a token with the provided kind to the output stream.
func (l *Lexer[T]) Emit(kind T) {
	l.tokens <- Token[T]{kind, l.input[l.start:l.pos]}
	l.start = l.pos
}

// Errorf formats and emits a synthetic error token using the configured error
// kind (defaults to the zero value).
func (l *Lexer[T]) Errorf(format string, args ...any) State[T] {
	l.tokens <- Token[T]{
		Type:  l.error,
		Value: fmt.Sprintf(l.name+": "+format, args...),
	}

	return nil
}

func (l *Lexer[T]) run(init State[T]) {
	for state := init; state != nil; {
		state = state(l)
	}

	close(l.tokens)
}
