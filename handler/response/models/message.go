package models

type Message struct {
	Message interface{} `json:"message"`
}

func NewMessage(data interface{}) *Message {
	return &Message{data}
}
