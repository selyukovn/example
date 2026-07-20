package kafapi

import "fmt"

// ---------------------------------------------------------------------------------------------------------------------

// ErrorDecoding
//
// Данные не могут быть декодированы.
type ErrorDecoding struct {
	err error
}

func newErrorDecoding(err error) ErrorDecoding {
	return ErrorDecoding{err: err}
}

func (e ErrorDecoding) Error() string {
	return fmt.Sprintf("%T: %s", e, e.err.Error())
}

func (e ErrorDecoding) Unwrap() error {
	return e.err
}

// ---------------------------------------------------------------------------------------------------------------------

// ErrorMapping
//
// Данные декодированы, но не могут быть преобразованы в высокоуровневые типы/версии сообщений
type ErrorMapping struct {
	err     error
	Meta    *Meta
	RawData *RawData
}

func newErrorMapping(err error, meta *Meta, rawData *RawData) ErrorMapping {
	return ErrorMapping{
		err:     err,
		Meta:    meta,
		RawData: rawData,
	}
}

func (e ErrorMapping) Error() string {
	return fmt.Sprintf("%T: %s", e, e.err.Error())
}

func (e ErrorMapping) Unwrap() error {
	return e.err
}

// ---------------------------------------------------------------------------------------------------------------------

// ErrorUnsupported
//
// Данные не соответствуют ни одному известному типу/версии сообщения.
type ErrorUnsupported struct {
	Meta    *Meta
	RawData *RawData
}

func newErrorUnsupported(meta *Meta, rawData *RawData) ErrorUnsupported {
	return ErrorUnsupported{
		Meta:    meta,
		RawData: rawData,
	}
}

func (e ErrorUnsupported) Error() string {
	return fmt.Sprintf("неизвестный тип/версия сообщения #%v", e.Meta.Id)
}

// ---------------------------------------------------------------------------------------------------------------------

// ErrorHandling
//
// Ошибка обработки сообщения -- см. ServiceInterface.
type ErrorHandling struct {
	err error
}

func newErrorHandling(err error) ErrorHandling {
	return ErrorHandling{
		err: err,
	}
}

func (e ErrorHandling) Error() string {
	return e.err.Error()
}

func (e ErrorHandling) Unwrap() error {
	return e.err
}

// ---------------------------------------------------------------------------------------------------------------------
