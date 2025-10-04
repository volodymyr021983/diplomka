package auth

import "regexp"

// needed to validate Username while sign UP
func usernameValidator(value interface{}, tenantId string) *string {
	if len(value.(string)) < 3 {
		msg := "Usernames must be at least 3 characters long."
		return &msg
	}
	userNameCheck, err := regexp.Match(`^[A-Za-z0-9_-]+$`, []byte(value.(string)))
	if err != nil || !userNameCheck {
		msg := "Username must contain only alphanumeric, underscore or hyphen characters."
		return &msg
	}
	return nil
}
