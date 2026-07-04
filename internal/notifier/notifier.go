package notifier

type Message struct {
	Title   string
	Content string
	Level   string
	Fields  map[string]interface{}
}

type Notifier interface {
	Send(msg Message) error
}
