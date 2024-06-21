package service

import (
	"context"

	"github.com/newRational/soc/internal/model"
)

type repository interface {
	Add(ctx context.Context, order *model.Order) error
}

type OrdersService struct {
	repo repository
}

func NewOrdersService(repository repository) *OrdersService {
	return &OrdersService{repo: repository}
}

func (s *OrdersService) Add(ctx context.Context, order *model.AddOrderRequest) (*model.AddOrderResponse, error) {
	orderRepo := &model.Order{
		ID:          order.ID,
		Title:       order.Title,
		Level:       order.Level,
		Description: order.Description,
		UpdatedAt:   order.UpdatedAt,
		CreatedAt:   order.CreatedAt,
	}

	err := s.repo.Add(ctx, orderRepo)
	if err != nil {
		return nil, err
	}

	resp := &model.AddOrderResponse{
		ID:          order.ID,
		Title:       order.Title,
		Level:       order.Level,
		Description: order.Description,
		UpdatedAt:   order.UpdatedAt,
		CreatedAt:   order.CreatedAt,
	}

	return resp, nil
}
