package entity

type NotifierMessage struct {
	Message   interface{}
	Consumers []string
}
