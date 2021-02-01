package response

const (
	SuccessCode = "200"
	SuccessDesc = "Success."

	ErrBadRequestCode = "400"
	ErrBadRequestDesc = "Cannot bind request."

	ErrInvalidRequestCode = "400001"
	ErrInvalidRequestDesc = "Cannot parse public key."

	ErrInternalServerCode = "500"
	ErrInternalServerDesc = "Internal service error."
)

const (
	ValidateFieldError string = "Invalid parameters"
)
