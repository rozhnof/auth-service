package kafka

import "github.com/google/uuid"

type Message struct {
	Key   uuid.UUID `json:"key"   db:"key"`
	Value []byte    `json:"value" db:"value"`
	Topic string    `json:"topic" db:"topic"`
}
