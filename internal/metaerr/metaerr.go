package metaerr

import "errors"

type errMetadata struct {
	// err is the original error.
	err error
	// metadata is the container for structured error context.
	metadata []any
}

// Error returns the original error message,
// ensuring compatibility with the standard error interface.
func (e *errMetadata) Error() string {
	return e.err.Error()
}

// Unwrap allows errors wrapped by errMetadata to be compatible with
// standard error unwrapping mechanism.
func (e *errMetadata) Unwrap() error {
	return e.err
}

// WithMetadata creates a new error with metadata.
func WithMetadata(err error, pairs ...any) error {
	return &errMetadata{
		err:      err,
		metadata: pairs,
	}
}

// GetMetadata returns metadata from the error chain
// if there is no metadata in the chain, it will return an empty slice.
func GetMetadata(err error) []any {
	data := make([]any, 0)

	// we will iterate over all errors in the chain
	// and merge metadata from all of them
	for err != nil {
		// if current error is a metadata error
		// we will add its metadata to the existing
		// already collected metadata
		if e, ok := err.(*errMetadata); ok {
			data = append(data, e.metadata...)
		}

		err = errors.Unwrap(err)
	}

	return data
}
