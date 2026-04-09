package models

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
	Status  int               `json:"status"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}
