package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// Tries to login with an incorrect email / password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// Tries to signup with an email address that is alread in use
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
