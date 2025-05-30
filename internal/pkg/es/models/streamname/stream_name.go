// Package streamname provides stream name.
package streamname

import (
	"fmt"
	"strings"

	reflect "github.com/goccy/go-reflect"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
)

// StreamName is a stream name.
type StreamName string

// GetId gets the id from the stream name.
func (n StreamName) GetId() uuid.UUID {
	name := n.String()
	index := strings.Index(name, "-") + 1
	id := name[index:]

	return uuid.FromStringOrNil(id)
}

// String gets the string from the stream name.
func (n StreamName) String() string {
	return string(n)
}

// For gets stream name for aggregate.
func For[T models.IHaveEventSourcedAggregate](aggregate T) StreamName {
	var aggregateName string
	if t := reflect.TypeOf(aggregate); t.Kind() == reflect.Ptr {
		aggregateName = reflect.TypeOf(aggregate).Elem().Name()
	} else {
		aggregateName = reflect.TypeOf(aggregate).Name()
	}

	return StreamName(fmt.Sprintf("%s-%s", strings.ToLower(aggregateName), aggregate.ID().String()))
}

// ForID gets stream name for aggregate id.
func ForID[T models.IHaveEventSourcedAggregate](aggregateID uuid.UUID) StreamName {
	var aggregate T
	var aggregateName string
	if t := reflect.TypeOf(aggregate); t.Kind() == reflect.Ptr {
		aggregateName = reflect.TypeOf(aggregate).Elem().Name()
	} else {
		aggregateName = reflect.TypeOf(aggregate).Name()
	}

	return StreamName(fmt.Sprintf("%s-%s", strings.ToLower(aggregateName), aggregateID.String()))
}
