package http

import "encoding/json"


type RequestError struct {
	StatusCode int
	Msg   string
}

func (e *RequestError) Error() string {
	err, _ := json.Marshal(e)
	return string(err)
}

func Error(statusCode int, msg string) *RequestError {
	return &RequestError{
		StatusCode: statusCode,
		Msg:   		msg,
	}
}
