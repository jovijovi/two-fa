package twofa

import (
	"context"
)

// Options for func
type Options struct {
	// HashFunc interface
	HashFunc IHashFunc

	// With hash
	WithHash bool

	// Key name
	KeyName string

	// With QR
	WithQR bool

	// Options for implementations of the interface can be stored in a context
	Context context.Context
}

// OptionFunc used to initialise
type OptionFunc func(opts *Options)

// NewOptions new options
func NewOptions(optionFunc ...OptionFunc) Options {
	opts := Options{
		Context:  context.Background(),
		HashFunc: DefaultHashFunc(),
		WithHash: false,
	}

	for _, f := range optionFunc {
		f(&opts)
	}

	return opts
}

// WithDefaultHashFunc option to configure default hash function (sha1)
func WithDefaultHashFunc() OptionFunc {
	return func(o *Options) {
		o.WithHash = true
	}
}

// WithHashFunc option to configure hash function
func WithHashFunc(hashFunc IHashFunc) OptionFunc {
	return func(o *Options) {
		o.HashFunc = hashFunc
		o.WithHash = true
	}
}

// WithQR option to configure hash function
func WithQR(name string) OptionFunc {
	return func(o *Options) {
		o.KeyName = name
		o.WithQR = true
	}
}
