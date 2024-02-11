package models

type ErrorResponse struct {
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

type AdvancedErrorResponse struct {
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (e AdvancedErrorResponse) Error() string {
	return e.Message
}
