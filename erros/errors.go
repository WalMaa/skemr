package errors

import "errors"

var (
	ErrDatabaseAlreadyExists = errors.New("database already exists")
	ErrDatabaseNotFound      = errors.New("database not found")
	ErrProjectNotFound       = errors.New("project not found")
)
