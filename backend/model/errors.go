package model

type WebServiceException struct {
	Status  int
	Code    string
	Message string
}

func (e WebServiceException) Error() string {
	return e.Message
}

func (e WebServiceException) Is(err error) bool {
	a, ok := err.(WebServiceException)
	if !ok {
		return false
	}
	return e.Code == a.Code
}

func NewWebServiceException(message string, code string, status int) error {
	e := WebServiceException{
		Status:  status,
		Code:    code,
		Message: message,
	}
	return &e
}

type ErrorIs interface {
	Is(error) bool
}
