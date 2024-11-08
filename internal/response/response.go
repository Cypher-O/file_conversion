package response

// APIResponse is the structure of the response for all API calls
type APIResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewSuccessResponse creates a successful response
func NewSuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Code:    0,
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(code int, message string) APIResponse {
	return APIResponse{
		Code:    code,
		Status:  "error",
		Message: message,
		Data:    nil,
	}
}
