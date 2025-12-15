package dto

// SuccessResponse wraps successful API responses
type SuccessResponse[T any] struct {
	IsSuccess bool   `json:"is_success" example:"true"`
	Error     string `json:"error" example:""`
	Data      T      `json:"data"`
}

type EmptyData struct{}

// ErrorResponse wraps error API responses
type ErrorResponse struct {
	IsSuccess bool   `json:"is_success" example:"false"`
	Error     string `json:"error" example:"some error message"`
	Data      any    `json:"data" swaggertype:"object"`
}

// Endpoint response types
type (
	SuccessUserAndTokenResponse = SuccessResponse[UserAndTokenResponse]
	SuccessTokenResponse        = SuccessResponse[TokenResponse]
	SuccessEmptyResponse        = SuccessResponse[EmptyData]
	SuccessUserResponse         = SuccessResponse[UserResponse]
	SuccessUpdateUserResponse   = SuccessResponse[UpdateUserResponse]
)
