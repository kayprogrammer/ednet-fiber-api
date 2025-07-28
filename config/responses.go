package config

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Status  string             `json:"status"`
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Data    *map[string]string `json:"data,omitempty"`
}

// Error codes
var ERR_UNAUTHORIZED_USER = "unauthorized_user"
var ERR_INSTRUCTORS_ONLY = "instructors_only"
var ERR_ADMINS_ONLY = "admins_only"
var ERR_NETWORK_FAILURE = "network_failure"
var ERR_SERVER_ERROR = "server_error"
var ERR_INVALID_REQUEST = "invalid_request"
var ERR_INVALID_PARAM = "invalid_param"
var ERR_INVALID_ENTRY = "invalid_entry"
var ERR_INCORRECT_EMAIL = "incorrect_email"
var ERR_INCORRECT_OTP = "incorrect_otp"
var ERR_EXPIRED_OTP = "expired_otp"
var ERR_INCORRECT_TOKEN = "incorrect_token"
var ERR_EXPIRED_TOKEN = "expired_token"
var ERR_INVALID_AUTH = "invalid_auth"
var ERR_INVALID_TOKEN = "invalid_token"
var ERR_INVALID_USER = "invalid_user"
var ERR_INVALID_PAYLOAD = "invalid_payload"
var ERR_INVALID_CREDENTIALS = "invalid_credentials"
var ERR_UNVERIFIED_USER = "unverified_user"
var ERR_NON_EXISTENT = "non_existent"
var ERR_INVALID_OWNER = "invalid_owner"
var ERR_INVALID_PAGE = "invalid_page"
var ERR_INVALID_VALUE = "invalid_value"
var ERR_NOT_ALLOWED = "not_allowed"
var ERR_INVALID_DATA_TYPE = "invalid_data_type"
var ERR_PASSWORD_MISMATCH = "password_does_not_match"
var ERR_PASSWORD_SAME = "same_password"
var ERR_NOT_FOUND = "not_found"
var ERR_LIMITS_REACHED = "limits_reached"
var ERR_FORBIDDEN = "forbidden"
var ERR_TOO_MANY_REQUESTS = "too_many_requests"

func RequestErr(code string, message string, opts ...map[string]string) ErrorResponse {
	var data *map[string]string
	// Check if data is provided as an argument
	if len(opts) > 0 {
		data = &opts[0]
	}
	resp := ErrorResponse{Status: "failure", Code: code, Message: message, Data: data}
	return resp
}

func NotFoundErr(message string) ErrorResponse {
	return RequestErr(ERR_NON_EXISTENT, message)
}

func InvalidParamErr(message string) ErrorResponse {
	return RequestErr(ERR_INVALID_PARAM, message)
}

func RateLimitError(message string) ErrorResponse {
	return RequestErr(ERR_LIMITS_REACHED, message)
}

func ForbiddenErr(message string) ErrorResponse {
	return RequestErr(ERR_FORBIDDEN, message)
}

func ServerErr(message string) ErrorResponse {
	return RequestErr(ERR_SERVER_ERROR, message)
}

func ValidationErr(field string, message string) ErrorResponse {
	data := map[string]string{field: message}
	return RequestErr(ERR_INVALID_ENTRY, "Invalid Entry", data)
}

func APIError(c *fiber.Ctx, code int, data ErrorResponse) error {
	return c.Status(code).JSON(data)
}
