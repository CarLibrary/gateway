package serializer

const (
	SQL = "sql error"
	BIND="binding error"
	MIDDLEWARE = "middleware error"
	OTHER = "other error"
)

const (
	OK = 40000
	BIND_ERROR = 40001
	GRPC_ERROR =40002
	TOKEN_ERROR = 40003
)

type ErrorResponse struct {
	ErrType string
	Msg string

}