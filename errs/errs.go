package errs

import "errors"


// ErrNoSuchEntity is thrown when no entity is found for the required key
var ErrNoSuchEntity = errors.New("no such entity")

// IsNoSuchEntityError return if the error is a wrapped (or not) not such entity error
func IsNoSuchEntityError(e error) bool {
	return errors.Is(e, ErrNoSuchEntity)
}

// ErrPreConditionFailed is thrown when the precondition fails
var ErrPreConditionFailed = errors.New("Precondition Failed")

// IsPreConditionFailedError return if the error is a wrapped (or not) precondition fails
func IsPreConditionFailed(e error) bool {
	return errors.Is(e, ErrPreConditionFailed)
}

