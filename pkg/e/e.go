package e

type ErrorResponse struct {
	Error string `json:"error"`
}

const (
	ErrBadRequest                = "bad request"
	ErrInvalidCredentials        = "invalid credentials"
	ErrPhoneNotExists            = "this phone already exists"
	ErrEmailNotExists            = "this email already exists"
	ErrNoSuchUser                = "user not found"
	ErrNotFound                  = "not found"
	ErrInvalidFormat             = "request body must be in JSON format"
	ErrInternalServer            = "internal server error"
	ErrUnauthorized              = "unauthorized"
	ErrForbidden                 = "forbidden"
	ErrInvalidToken              = "invalid token"
	ErrTokenExpired              = "token expired"
	ErrRefreshNotFound           = "refresh token not found"
	ErrInvalidRefreshToken       = "invalid refresh token"
	ErrRefreshTokenExpired       = "refresh token expired"
	ErrAuthorizationHeaderNotSet = "authorization header not set"
	ErrUserNotFound              = "user not found"
	ErrUserContextNotFound       = "user context not found or not authorized to access this account yet"
	ErrAllowedParticipants       = "allowed number of participants 2"
	ErrChatExists                = " chat with the same participants already exists"
)
