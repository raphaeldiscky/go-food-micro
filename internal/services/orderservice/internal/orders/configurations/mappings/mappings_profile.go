// Package mappings contains the mappings for the orderservice.
package mappings

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"

	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/readmodels"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/valueobject"
	grpcOrderService "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/grpc/genproto"
)

// configureOrderToDtoMappings configures mappings from Order to DTOs.
func configureOrderToDtoMappings() error {
	// Order -> OrderDto
	if err := mapper.CreateMap[*aggregate.Order, *dtosV1.OrderDto](); err != nil {
		return err
	}

	// readmodels.OrderReadModel -> dtos.OrderReadDto
	if err := mapper.CreateMap[*readmodels.OrderReadModel, *dtosV1.OrderReadDto](); err != nil {
		return err
	}

	return nil
}

// configureDtoToOrderMappings configures mappings from DTOs to Order.
func configureDtoToOrderMappings() error {
	// OrderDto -> Order
	if err := mapper.CreateCustomMap[*dtosV1.OrderDto, *aggregate.Order](
		func(orderDto *dtosV1.OrderDto) (*aggregate.Order, error) {
			items, err := mapper.Map[[]*valueobject.ShopItem](orderDto.ShopItems)
			if err != nil {
				return nil, err
			}

			order, err := aggregate.NewOrder(
				orderDto.ID,
				items,
				orderDto.AccountEmail,
				orderDto.DeliveryAddress,
				orderDto.DeliveredTime,
				orderDto.CreatedAt,
			)
			if err != nil {
				return nil, err
			}

			return order, nil
		},
	); err != nil {
		return err
	}

	return nil
}

// configureGrpcMappings configures mappings to/from gRPC models.
func configureGrpcMappings() error {
	// dtos.OrderReadDto -> grpcOrderService.OrderReadModel
	if err := mapper.CreateCustomMap[*dtosV1.OrderReadDto, *grpcOrderService.OrderReadModel](
		func(orderReadDto *dtosV1.OrderReadDto) (*grpcOrderService.OrderReadModel, error) {
			if orderReadDto == nil {
				return nil, nil
			}
			items, err := mapper.Map[[]*grpcOrderService.ShopItemReadModel](orderReadDto.ShopItems)
			if err != nil {
				return nil, err
			}

			return &grpcOrderService.OrderReadModel{
				ID:              orderReadDto.ID,
				OrderID:         orderReadDto.OrderID,
				PaymentID:       orderReadDto.PaymentID,
				DeliveredTime:   timestamppb.New(orderReadDto.DeliveredTime),
				TotalPrice:      orderReadDto.TotalPrice,
				DeliveryAddress: orderReadDto.DeliveryAddress,
				AccountEmail:    orderReadDto.AccountEmail,
				Canceled:        orderReadDto.Canceled,
				Completed:       orderReadDto.Completed,
				Paid:            orderReadDto.Paid,
				Submitted:       orderReadDto.Submitted,
				CancelReason:    orderReadDto.CancelReason,
				ShopItems:       items,
				CreatedAt:       timestamppb.New(orderReadDto.CreatedAt),
				UpdatedAt:       timestamppb.New(orderReadDto.UpdatedAt),
			}, nil
		},
	); err != nil {
		return err
	}

	// aggregate.Order -> grpcOrderService.Order
	if err := mapper.CreateCustomMap[*aggregate.Order, *grpcOrderService.Order](
		func(order *aggregate.Order) (*grpcOrderService.Order, error) {
			items, err := mapper.Map[[]*grpcOrderService.ShopItem](order.ShopItems())
			if err != nil {
				return nil, err
			}

			return &grpcOrderService.Order{
				OrderID:         order.ID().String(),
				DeliveryAddress: order.DeliveryAddress(),
				DeliveredTime:   timestamppb.New(order.DeliveredTime()),
				AccountEmail:    order.AccountEmail(),
				Canceled:        order.Canceled(),
				Completed:       order.Completed(),
				Paid:            order.Paid(),
				CancelReason:    order.CancelReason(),
				Submitted:       order.Submitted(),
				TotalPrice:      order.TotalPrice(),
				CreatedAt:       timestamppb.New(order.CreatedAt()),
				UpdatedAt:       timestamppb.New(order.UpdatedAt()),
				ShopItems:       items,
				PaymentID:       order.PaymentID().String(),
			}, nil
		},
	); err != nil {
		return err
	}

	return nil
}

// configureOrderMappings configures the order-related mappings.
func configureOrderMappings() error {
	if err := configureOrderToDtoMappings(); err != nil {
		return err
	}

	if err := configureDtoToOrderMappings(); err != nil {
		return err
	}

	if err := configureGrpcMappings(); err != nil {
		return err
	}

	return nil
}

// configureShopItemMappings configures the shop item-related mappings.
func configureShopItemMappings() error {
	// ShopItem -> ShopItemDto
	if err := mapper.CreateMap[*valueobject.ShopItem, *dtosV1.ShopItemDto](); err != nil {
		return err
	}

	// ShopItemDto -> ShopItem
	if err := mapper.CreateCustomMap[*dtosV1.ShopItemDto, *valueobject.ShopItem](
		func(src *dtosV1.ShopItemDto) (*valueobject.ShopItem, error) {
			return valueobject.CreateNewShopItem(
				src.Title,
				src.Description,
				src.Quantity,
				src.Price,
			), nil
		},
	); err != nil {
		return err
	}

	// dtos.ShopItemDto -> readmodels.ShopItemReadModel
	if err := mapper.CreateMap[*dtosV1.ShopItemDto, *readmodels.ShopItemReadModel](); err != nil {
		return err
	}

	// readmodels.ShopItemReadModel -> dtos.ShopItemReadDto
	if err := mapper.CreateMap[*readmodels.ShopItemReadModel, *dtosV1.ShopItemReadDto](); err != nil {
		return err
	}

	// dtos.ShopItemReadDto -> grpcOrderService.ShopItemReadModel
	if err := mapper.CreateMap[*dtosV1.ShopItemReadDto, *grpcOrderService.ShopItemReadModel](); err != nil {
		return err
	}

	// valueobject.ShopItem -> grpcOrderService.ShopItem
	if err := mapper.CreateCustomMap[*valueobject.ShopItem, *grpcOrderService.ShopItem](
		func(src *valueobject.ShopItem) (*grpcOrderService.ShopItem, error) {
			return &grpcOrderService.ShopItem{
				Title:       src.Title(),
				Description: src.Description(),
				Quantity:    src.Quantity(),
				Price:       src.Price(),
			}, nil
		},
	); err != nil {
		return err
	}

	// grpcOrderService.ShopItem -> valueobject.ShopItem
	if err := mapper.CreateCustomMap[*grpcOrderService.ShopItem, *valueobject.ShopItem](
		func(src *grpcOrderService.ShopItem) (*valueobject.ShopItem, error) {
			return valueobject.CreateNewShopItem(
				src.Title,
				src.Description,
				src.Quantity,
				src.Price,
			), nil
		},
	); err != nil {
		return err
	}

	// grpcOrderService.ShopItem -> dtos.ShopItemDto
	if err := mapper.CreateMap[*grpcOrderService.ShopItem, *dtosV1.ShopItemDto](); err != nil {
		return err
	}

	return nil
}

// configureListResultMappings configures the list result-related mappings.
func configureListResultMappings() error {
	// ListResult[OrderReadDto] -> GetOrdersRes
	if err := mapper.CreateCustomMap[*utils.ListResult[*dtosV1.OrderReadDto], *grpcOrderService.GetOrdersRes](
		func(orders *utils.ListResult[*dtosV1.OrderReadDto]) (*grpcOrderService.GetOrdersRes, error) {
			o, err := mapper.Map[[]*grpcOrderService.OrderReadModel](orders.Items)
			if err != nil {
				return nil, err
			}

			return &grpcOrderService.GetOrdersRes{
				Pagination: &grpcOrderService.Pagination{
					Size:       int32(orders.Size),
					Page:       int32(orders.Page),
					TotalItems: orders.TotalItems,
					TotalPages: int32(orders.TotalPage),
				},
				Orders: o,
			}, nil
		},
	); err != nil {
		return err
	}

	return nil
}

// ConfigureOrdersMappings configures all the orders mappings.
func ConfigureOrdersMappings() error {
	if err := configureOrderMappings(); err != nil {
		return err
	}

	if err := configureShopItemMappings(); err != nil {
		return err
	}

	if err := configureListResultMappings(); err != nil {
		return err
	}

	return nil
}
