// Package errorz provides utility functions for errors.
package errorz

import "errors"

type (
	// SingleError is an error that unwraps into one error.
	SingleError interface{ Unwrap() error }
	// JoinedError is an error that unwraps into multiple errors.
	JoinedError interface{ Unwrap() []error }
)

// IsSingle checks if the error is a SingleError.
//
// Can return true even if the error is a JoinedError.
func IsSingle(err error) bool {
	_, ok := err.(SingleError)

	return ok
}

// IsJoined checks if the error is a JoinedError.
//
// Can return true even if the error is a SingleError.
func IsJoined(err error) bool {
	_, ok := err.(JoinedError)

	return ok
}

// IsUnwrappable checks if the error can be unwrapped.
func IsUnwrappable(err error) bool {
	return IsSingle(err) || IsJoined(err)
}

// UnwrapAll unwraps JoinedError or SingleError into a slice of errors.
//
// Returns nil if the error is not unwrappable or equal to nil.
func UnwrapAll(err error) []error {
	if err == nil {
		return nil
	}

	if e, ok := err.(JoinedError); ok {
		return e.Unwrap()
	}

	if e, ok := err.(SingleError); ok {
		return []error{
			e.Unwrap(),
		}
	}

	return nil
}

// IsMatching checks if the error is one of the targets.
//
// Returns false if there are no targets or error is nil.
func IsMatching(err error, targets ...error) bool {
	if err == nil || len(targets) == 0 {
		return false
	}

	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

// Matching returns a slice of all the errors matching the targets.
//
// Returns nil if there are no targets or errors.
func Matching(errs []error, targets ...error) []error {
	if len(errs) == 0 || len(targets) == 0 {
		return nil
	}

	matching := make([]error, 0)

	for i := range errs {
		if IsMatching(errs[i], targets...) {
			matching = append(matching, errs[i])
		}
	}

	return matching
}

// Allow returns an error if it matches any of the targets.
//
// Returns nil if there are no targets or error is nil.
func Allow(err error, targets ...error) error {
	if err == nil || len(targets) == 0 {
		return nil
	}

	if IsMatching(err, targets...) {
		return err
	}

	return nil
}

// Ignore returns an error if it does not match any of the targets.
//
// Returns nil if there are no targets or error is nil.
func Ignore(err error, targets ...error) error {
	if err == nil || len(targets) == 0 {
		return nil
	}

	if IsMatching(err, targets...) {
		return nil
	}

	return err
}
