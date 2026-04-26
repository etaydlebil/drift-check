package output

import "errors"

// errWriter is a shared io.Writer that always returns an error, used across
// multiple test files in this package.
type errWriter struct{}

func (errWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write error")
}
