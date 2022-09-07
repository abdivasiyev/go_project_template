package validator

import "fmt"

type CustomError interface {
	Field() string
	Message() string
}

type ValidationError struct {
	fieldName string
	message   string
}

func NewValidationError(fieldName, message string) ValidationError {
	return ValidationError{
		fieldName: fieldName,
		message:   message,
	}
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.fieldName, v.message)
}

func (v ValidationError) Field() string {
	return v.fieldName
}

func (v ValidationError) Message() string {
	return v.message
}
