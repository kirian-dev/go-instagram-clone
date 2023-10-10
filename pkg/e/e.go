package e

type ErrorResponse struct {
	Error string `json:"error"`
}

const (
	ErrBadRequest                = "bad request"
	ErrInvalidCredentials        = "invalid credentials"
	ErrPhoneNotExists            = "this phone already exists"
	ErrEmailExists               = "this email already exists"
	ErrEmailNotExists            = "this email not exists"
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
	ErrChatExists                = "chat with the same participants already exists"
	ErrUserNotCreatedByHimself   = "user not created a chat by himself"
	ErrCreateChatNoRights        = "create a chat without rights to members"
	ErrNoRights                  = "no rights"
	ErrChatNotFound              = "chat not found"
	ErrNotGroupChat              = "not a group chat"
	ErrParticipantNotFound       = "participants not found"
	ErrNotCorrectSender          = "not correct sender"
	ErrNotValidChatType          = "chat type not valid"
	ErrMessageNotFound           = "message not found"
	ErrPasswordDoesNotMatch      = "password does not match"
	ErrEmailMustBeVerified       = "email must be verified"
	ErrFailedToGetFile           = "Failed to receive file"
	ErrBigFileSize               = "File size exceeds 2MB"
	ErrFileMustBeCSV             = "The file must be in CSV format"
	ErrCreateFile                = "Create file failed"
)
