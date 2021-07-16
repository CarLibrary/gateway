package serializer

type CommonResponse struct {
	Statuscode int	`json:"statuscode"`
	Data interface{}	`json:"data"`
	Msg 	string	`json:"msg"`
	Error   ErrorResponse	`json:"error"`
}



