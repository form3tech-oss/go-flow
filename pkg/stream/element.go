package stream

type Element struct {
	Value interface{}
	Error error
}

func Value(value interface{}) Element {
	return Element{Value: value}
}

func Error(err error) Element {
	return Element{Error: err}
}
