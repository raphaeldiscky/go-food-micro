// Package persistmessage provides a store message.
package persistmessage

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// MessageDeliveryType is a type that represents the delivery type of a message.
type MessageDeliveryType int

const (
	Outbox   MessageDeliveryType = 1
	Inbox    MessageDeliveryType = 2
	Internal MessageDeliveryType = 4
)

// MessageStatus is a type that represents the status of a message.
type MessageStatus int

const (
	Stored    MessageStatus = 1
	Processed MessageStatus = 2
)

// StoreMessage is a struct that represents a message in the store.
type StoreMessage struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	DataType      string
	Data          string
	CreatedAt     time.Time `gorm:"default:current_timestamp"`
	RetryCount    int
	MessageStatus MessageStatus
	DeliveryType  MessageDeliveryType
}

// NewStoreMessage is a function that creates a new store message.
func NewStoreMessage(
	id uuid.UUID,
	dataType string,
	data string,
	deliveryType MessageDeliveryType,
) *StoreMessage {
	return &StoreMessage{
		ID:            id,
		DataType:      dataType,
		Data:          data,
		CreatedAt:     time.Now(),
		MessageStatus: Stored,
		RetryCount:    0,
		DeliveryType:  deliveryType,
	}
}

// ChangeState is a function that changes the state of a message.
func (sm *StoreMessage) ChangeState(messageStatus MessageStatus) {
	sm.MessageStatus = messageStatus
}

// IncreaseRetry is a function that increases the retry count of a message.
func (sm *StoreMessage) IncreaseRetry() {
	sm.RetryCount++
}

// TableName is a function that returns the table name of a message.
func (sm *StoreMessage) TableName() string {
	return "store_messages"
}
