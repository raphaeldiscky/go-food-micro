// Package messagepersistence provides a set of functions for the message persistence.
package messagepersistence

import (
	"context"
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/persistmessage"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// postgresMessagePersistenceService is a struct that contains the postgres message persistence service.
type postgresMessagePersistenceService struct {
	messagingDBContext *PostgresMessagePersistenceDBContext
	messageSerializer  serializer.MessageSerializer
	logger             logger.Logger
}

// Process processes a single message by ID.
func (m *postgresMessagePersistenceService) Process(messageID string, ctx context.Context) error {
	id, err := uuid.FromString(messageID)
	if err != nil {
		return customErrors.NewBadRequestErrorWrap(err, "invalid message ID format")
	}

	storeMessage, err := m.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if storeMessage.MessageStatus != persistmessage.Stored {
		return customErrors.NewConflictErrorWrap(
			errors.New("message is not in stored state"),
			fmt.Sprintf("message %s is in state %d", messageID, storeMessage.MessageStatus),
		)
	}

	return m.ChangeState(ctx, id, persistmessage.Processing)
}

// ProcessAll processes all stored messages.
func (m *postgresMessagePersistenceService) ProcessAll(ctx context.Context) error {
	storeMessages, err := m.GetAllActive(ctx)
	if err != nil {
		return err
	}

	for _, msg := range storeMessages {
		if err := m.ChangeState(ctx, msg.ID, persistmessage.Processing); err != nil {
			m.logger.Errorf("Failed to process message %s: %v", msg.ID, err)
		}
	}

	return nil
}

// AddPublishMessage adds a message to be published.
func (m *postgresMessagePersistenceService) AddPublishMessage(
	messageEnvelope types.MessageEnvelope,
	ctx context.Context,
) error {
	return m.AddMessageCore(ctx, messageEnvelope, persistmessage.Publish)
}

// AddReceivedMessage adds a received message.
func (m *postgresMessagePersistenceService) AddReceivedMessage(
	messageEnvelope types.MessageEnvelope,
	ctx context.Context,
) error {
	return m.AddMessageCore(ctx, messageEnvelope, persistmessage.Received)
}

// AddMessageCore adds a message core.
func (m *postgresMessagePersistenceService) AddMessageCore(
	ctx context.Context,
	messageEnvelope types.MessageEnvelope,
	deliveryType persistmessage.MessageDeliveryType,
) error {
	if messageEnvelope.Message == nil {
		return errors.New("messageEnvelope.Message is nil")
	}

	var id string
	switch message := messageEnvelope.Message.(type) {
	case types.IMessage:
		id = message.GeMessageId()
	// case IInternalCommand:
	//	id = message.InternalCommandId
	default:
		id = uuid.NewV4().String()
	}

	data, err := m.messageSerializer.SerializeEnvelop(messageEnvelope)
	if err != nil {
		return err
	}

	uuidId, err := uuid.FromString(id)
	if err != nil {
		return err
	}

	storeMessage := persistmessage.NewStoreMessage(
		uuidId,
		messageEnvelope.Message.GetMessageFullTypeName(),
		string(data.Data),
		deliveryType,
	)

	err = m.Add(ctx, storeMessage)
	if err != nil {
		return err
	}

	m.logger.Infof(
		"Message with id: %v and delivery type: %v saved in persistence message store",
		id,
		deliveryType,
	)

	return nil
}

// NewPostgresMessageService creates a new postgres message service.
func NewPostgresMessageService(
	postgresMessagePersistenceDBContext *PostgresMessagePersistenceDBContext,
	l logger.Logger,
) persistmessage.MessagePersistenceService {
	return &postgresMessagePersistenceService{
		messagingDBContext: postgresMessagePersistenceDBContext,
		logger:             l,
	}
}

// Add adds a message.
func (m *postgresMessagePersistenceService) Add(
	ctx context.Context,
	storeMessage *persistmessage.StoreMessage,
) error {
	dbContext := m.messagingDBContext.WithTxIfExists(ctx)

	// https://gorm.io/docs/create.html
	result := dbContext.DB().Create(storeMessage)
	if result.Error != nil {
		return customErrors.NewConflictErrorWrap(
			result.Error,
			"storeMessage already exists",
		)
	}

	m.logger.Infof("Number of affected rows are: %d", result.RowsAffected)

	return nil
}

// Update updates a message.
func (m *postgresMessagePersistenceService) Update(
	ctx context.Context,
	storeMessage *persistmessage.StoreMessage,
) error {
	dbContext := m.messagingDBContext.WithTxIfExists(ctx)

	// https://gorm.io/docs/update.html
	result := dbContext.DB().Updates(storeMessage)
	if result.Error != nil {
		return customErrors.NewInternalServerErrorWrap(
			result.Error,
			"error in updating the storeMessage",
		)
	}

	m.logger.Infof("Number of affected rows are: %d", result.RowsAffected)

	return nil
}

// ChangeState changes the state of a message.
func (m *postgresMessagePersistenceService) ChangeState(
	ctx context.Context,
	messageID uuid.UUID,
	status persistmessage.MessageStatus,
) error {
	storeMessage, err := m.GetByID(ctx, messageID)
	if err != nil {
		return customErrors.NewNotFoundErrorWrap(
			err,
			fmt.Sprintf(
				"storeMessage with id `%s` not found in the database",
				messageID.String(),
			),
		)
	}

	storeMessage.MessageStatus = status
	err = m.Update(ctx, storeMessage)

	return err
}

func (m *postgresMessagePersistenceService) GetAllActive(
	ctx context.Context,
) ([]*persistmessage.StoreMessage, error) {
	var storeMessages []*persistmessage.StoreMessage

	predicate := func(sm *persistmessage.StoreMessage) bool {
		return sm.MessageStatus == persistmessage.Stored
	}

	dbContext := m.messagingDBContext.WithTxIfExists(ctx)
	result := dbContext.DB().Where(predicate).Find(&storeMessages)
	if result.Error != nil {
		return nil, result.Error
	}

	return storeMessages, nil
}

func (m *postgresMessagePersistenceService) GetByFilter(
	ctx context.Context,
	predicate func(*persistmessage.StoreMessage) bool,
) ([]*persistmessage.StoreMessage, error) {
	var storeMessages []*persistmessage.StoreMessage

	dbContext := m.messagingDBContext.WithTxIfExists(ctx)
	result := dbContext.DB().Where(predicate).Find(&storeMessages)

	if result.Error != nil {
		return nil, result.Error
	}

	return storeMessages, nil
}

// GetByID gets a message by id.
func (m *postgresMessagePersistenceService) GetByID(
	_ context.Context,
	id uuid.UUID,
) (*persistmessage.StoreMessage, error) {
	var storeMessage *persistmessage.StoreMessage

	// https://gorm.io/docs/query.html#Retrieving-objects-with-primary-key
	// https://gorm.io/docs/query.html#Struct-amp-Map-Conditions
	// https://gorm.io/docs/query.html#Inline-Condition
	// https://gorm.io/docs/advanced_query.html
	result := m.messagingDBContext.DB().Find(&storeMessage, id)
	if result.Error != nil {
		return nil, customErrors.NewNotFoundErrorWrap(
			result.Error,
			fmt.Sprintf(
				"storeMessage with id `%s` not found in the database",
				id.String(),
			),
		)
	}

	m.logger.Infof("Number of affected rows are: %d", result.RowsAffected)

	return storeMessage, nil
}

// Remove removes a message.
func (m *postgresMessagePersistenceService) Remove(
	ctx context.Context,
	storeMessage *persistmessage.StoreMessage,
) (bool, error) {
	id := storeMessage.ID

	storeMessage, err := m.GetByID(ctx, id)
	if err != nil {
		return false, customErrors.NewNotFoundErrorWrap(
			err,
			fmt.Sprintf(
				"storeMessage with id `%s` not found in the database",
				id.String(),
			),
		)
	}

	dbContext := m.messagingDBContext.WithTxIfExists(ctx)

	result := dbContext.DB().Delete(storeMessage, id)
	if result.Error != nil {
		return false, customErrors.NewInternalServerErrorWrap(
			result.Error,
			fmt.Sprintf(
				"error in deleting storeMessage with id `%s` in the database",
				id.String(),
			),
		)
	}

	m.logger.Infof("Number of affected rows are: %d", result.RowsAffected)

	return true, nil
}

// CleanupMessages cleans up the messages.
func (m *postgresMessagePersistenceService) CleanupMessages(
	ctx context.Context,
) error {
	predicate := func(sm *persistmessage.StoreMessage) bool {
		return sm.MessageStatus == persistmessage.Processed
	}

	dbContext := m.messagingDBContext.WithTxIfExists(ctx)

	result := dbContext.DB().
		Where(predicate).
		Delete(&persistmessage.StoreMessage{})

	if result.Error != nil {
		return result.Error
	}

	m.logger.Infof("Number of affected rows are: %d", result.RowsAffected)

	return nil
}
