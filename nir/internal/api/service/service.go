package service

import (
	"context"
	"time"

	"github.com/newRational/soc/infrastructure/logger"
	"github.com/newRational/soc/internal/model"
)

type sender interface {
	SendMessage(message model.Message) error
}

type Service struct {
	sender sender
}

func NewService(sender sender) *Service {
	return &Service{sender: sender}
}

func (s *Service) Add(_ context.Context, order *model.AddOrderRequest) error {
	err := s.sender.SendMessage(
		model.Message{
			ID:          order.ID,
			Title:       order.Title,
			Level:       order.Level,
			Description: order.Description,
			UpdatedAt:   order.UpdatedAt,
			CreatedAt:   order.CreatedAt,
		},
	)

	if err != nil {
		//fmt.Println("Api send sync message error: ", err)
		logger.Errorf("Api send sync message error: ", err)
		return model.ErrSendReq
	}

	return nil
}

func (s *Service) Get(_ context.Context, id uint64) (*model.ReadOrderResponse, error) {
	time := time.Date(2024, time.May, 30, 0, 0, 0, 0, time.UTC)
	order := model.Order{
		ID:          id,
		Title:       "mock@yandex.ru",
		Level:       200,
		Description: "Barabella",
		UpdatedAt:   "Mifcar",
		CreatedAt:   &time,
	}

	resp := &model.ReadOrderResponse{
		ID:          order.ID,
		Title:       order.Title,
		Level:       order.Level,
		Description: order.Description,
		UpdatedAt:   order.UpdatedAt,
		CreatedAt:   order.CreatedAt,
	}

	return resp, nil
}
