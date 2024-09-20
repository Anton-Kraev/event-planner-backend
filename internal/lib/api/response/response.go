package response

type Status string

const (
	StatusOK    Status = "OK"
	StatusError Status = "Error"
)

type Response struct {
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
