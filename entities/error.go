package entities

type CustomError struct {
	Error error
	MessageReponse string
	Code uint16
}