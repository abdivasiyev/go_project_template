package validator

func ValidatePassword(oldPassword, newPassword string) error {
	if oldPassword != "" && newPassword == "" {
		return NewValidationError("new_password", "new_password is a required field")
	}

	if oldPassword == "" && newPassword != "" {
		return NewValidationError("old_password", "old_password is a required field")
	}

	return nil
}
