// Package domain provides a module for the domain.
package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Entity is a entity.
type Entity struct {
	id         uuid.UUID
	entityType string
	createdAt  time.Time
	updatedAt  time.Time
}

// EntityDataModel is a entity data model.
type EntityDataModel struct {
	ID         uuid.UUID `json:"id"          bson:"id,omitempty"`
	EntityType string    `json:"entity_type" bson:"entity_type,omitempty"`
	CreatedAt  time.Time `json:"created_at"  bson:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"  bson:"updated_at"`
}

// IEntity is a entity.
type IEntity interface {
	ID() uuid.UUID
	CreatedAt() time.Time
	UpdatedAt() time.Time
	SetUpdatedAt(updatedAt time.Time)
	SetEntityType(entityType string)
	SetID(id uuid.UUID)
}

// NewEntityWithID creates a new entity with an id.
func NewEntityWithID(id uuid.UUID, entityType string) *Entity {
	return &Entity{
		id:         id,
		createdAt:  time.Now(),
		entityType: entityType,
	}
}

// NewEntity creates a new entity.
func NewEntity(entityType string) *Entity {
	return &Entity{
		createdAt:  time.Now(),
		entityType: entityType,
	}
}

// ID gets the id.
func (e *Entity) ID() uuid.UUID {
	return e.id
}

// CreatedAt gets the created at.
func (e *Entity) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt gets the updated at.
func (e *Entity) UpdatedAt() time.Time {
	return e.updatedAt
}

// EntityType gets the entity type.
func (e *Entity) EntityType() string {
	return e.entityType
}

// SetUpdatedAt sets the updated at.
func (e *Entity) SetUpdatedAt(updatedAt time.Time) {
	e.updatedAt = updatedAt
}

// SetEntityType sets the entity type.
func (e *Entity) SetEntityType(entityType string) {
	e.entityType = entityType
}

// SetID sets the id.
func (e *Entity) SetID(id uuid.UUID) {
	e.id = id
}

// ToDataModel converts the entity to a data model.
func (e *Entity) ToDataModel() *EntityDataModel {
	return &EntityDataModel{
		ID:         e.id,
		EntityType: e.entityType,
		CreatedAt:  e.createdAt,
		UpdatedAt:  e.updatedAt,
	}
}
