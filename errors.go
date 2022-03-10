package accounts

import "fmt"

type ConnectionError struct {
	Cause error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("connection error: %v", e.Cause.Error())
}

type ClientError struct {
	StatusCode int
	Message    string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("client side error %v: %v", e.StatusCode, e.Message)
}

type ServerError struct {
	StatusCode int
	Message    string
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("server side error %v: %v", e.StatusCode, e.Message)
}
