package bind

import "errors"

// ErrNotImplemented is returned until the generator is implemented.
var ErrNotImplemented = errors.New("bind: generator not implemented")

// Generate produces typed Go bindings for one document template.
func Generate(req Request) (Output, error) {
	_ = req
	return Output{}, ErrNotImplemented
}
