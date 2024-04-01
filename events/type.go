package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

// преимущество такого объявления пользовательского типа Type:
// сложнее ошибиться, нельзя передать что-то другое в функцию, получающую событие - ошибка на этапе компиляции
type Type int

const (
	Unknown Type = iota // Если мы не смогли определить тип события
	Message             // Если тип события корректный
)

type Event struct {
	Type Type
	Text string
	Meta interface{} // тут будет доп информация, например chatID или username
}
