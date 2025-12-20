package dto

type SuccessResponse[T any] struct {
	Data      T      `json:"data"`
	Error     string `json:"error"      example:""`
	IsSuccess bool   `json:"is_success" example:"true"`
}

type EmptyData struct{}

type ErrorResponse struct {
	Data      any    `json:"data"       swaggertype:"object"`
	Error     string `json:"error"                           example:"some error message"`
	IsSuccess bool   `json:"is_success"                      example:"false"`
}

type (
	SuccessUserAndTokenResponse = SuccessResponse[UserAndTokenResponse]
	SuccessTokenResponse        = SuccessResponse[TokenResponse]
	SuccessEmptyResponse        = SuccessResponse[EmptyData]
	SuccessUserResponse         = SuccessResponse[UserResponse]
	SuccessUpdateUserResponse   = SuccessResponse[UpdateUserResponse]
	SuccessUsersSearchResponse  = SuccessResponse[UsersSearchResponse]
)
