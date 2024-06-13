package api

type APIResponse struct {
	StatusCode int         `json:"statusCode"`
	Result     interface{} `json:"result,omitempty"`
}

func NewAPIResponse(statusCode int, result interface{}) APIResponse {
	return APIResponse{
		StatusCode: statusCode,
		Result:     result,
	}
}
