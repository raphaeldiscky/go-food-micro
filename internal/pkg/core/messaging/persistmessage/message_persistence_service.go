// Package persistmessage provides a message persistence service.
package persistmessage

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

// MessagePersistenceService is a type that represents a message persistence service.
type MessagePersistenceService interface {
	Add(ctx context.Context, storeMessage *StoreMessage) error
	Update(ctx context.Context, storeMessage *StoreMessage) error
	ChangeState(
		ctx context.Context,
		messageID uuid.UUID,
		status MessageStatus,
	) error
	GetAllActive(ctx context.Context) ([]*StoreMessage, error)
	GetByFilter(
		ctx context.Context,
		predicate func(*StoreMessage) bool,
	) ([]*StoreMessage, error)
	GetByID(ctx context.Context, id uuid.UUID) (*StoreMessage, error)
	Remove(ctx context.Context, storeMessage *StoreMessage) (bool, error)
	CleanupMessages(ctx context.Context) error
	Process(messageID string, ctx context.Context) error
	ProcessAll(ctx context.Context) error
	AddPublishMessage(
		messageEnvelope types.MessageEnvelope,
		ctx context.Context,
	) error
	AddReceivedMessage(
		messageEnvelope types.MessageEnvelope,
		ctx context.Context,
	) error
	// AddInternalMessage(
	//	internalCommand InternalMessage,
	//	ctx context.Context,
	// ) error
}
