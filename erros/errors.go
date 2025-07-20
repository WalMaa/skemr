package errors

import "errors"

var (
	ErrDatabaseAlreadyExists = errors.New("Database already exists")
	ErrDatabaseNotFound      = errors.New("database not found")
)
