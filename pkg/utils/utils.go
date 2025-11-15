package utils

// return json response for api responses
func ReturnJsonResponse(status string, message string, data map[string]interface{}, errors map[string]interface{}) ApiResponse {
	response := ApiResponse{
		Status:  status,
		Message: message,
		Data:    data,
		Errors:  errors,
	}
	return response
}
