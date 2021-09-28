package errs

import "errors"

//ErrEmailTaken is returned when trying to create a user the a taken email address
var ErrEmailTaken = errors.New("Email address already taken")


// ErrNoSuchEntity is thrown when no entity is found for the required key
var ErrNoSuchEntity = errors.New("no such entity")

var ErrFileAlreadyExist  = errors.New("File already exists!")
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

