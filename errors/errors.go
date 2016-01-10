package errors

import "fmt"

// Err is the string representation of an error. It may include formating directives.
type Err string

// Class defines where an error orginated from.
type Class int

const (
	Unknown Class = iota
	DB
	Serializer
	Codec
	Thunder
)

type Error interface {
	Class() Class
	Err() Err
	Error() string
}

type err struct {
	class   Class
	err     Err
	details []interface{}
}

func (e err) Class() Class {
	return e.class
}

func (e err) Err() Err {
	return e.err
}
func (e err) Error() string {
	return fmt.Sprintf(string(e.err), e.details...)
}

func New(c Class, e Err, details ...interface{}) Error {
	return err{class: c, err: e, details: details}
}
