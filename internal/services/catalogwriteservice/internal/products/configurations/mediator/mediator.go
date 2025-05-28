package mediator

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"

// RegisterMediatorHandlers is a function that registers the mediator handlers.
func RegisterMediatorHandlers(handlers []cqrs.HandlerRegisterer) error {
	for _, handler := range handlers {
		err := handler.RegisterHandler()
		if err != nil {
			return err
		}
	}

	return nil
}
