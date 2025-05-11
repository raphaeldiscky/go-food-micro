package mediator

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"

func RegisterMediatorHandlers(handlers []cqrs.HandlerRegisterer) error {
	for _, handler := range handlers {
		err := handler.RegisterHandler()
		if err != nil {
			return err
		}
	}

	return nil
}
