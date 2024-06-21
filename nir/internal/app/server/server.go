package server

import (
	"context"

	"github.com/newRational/soc/infrastructure/logger"
	"github.com/newRational/soc/internal/model"
)

type sender interface {
	SendMessage(message model.Message) error
}

type service interface {
	Add(ctx context.Context, order *model.AddOrderRequest) (*model.AddOrderResponse, error)
}

type Server struct {
	service service
	sender  sender
	msgChan <-chan model.Message
}

func NewServer(service service, sender sender, msgChan chan model.Message) Server {
	return Server{
		service: service,
		msgChan: msgChan,
		sender:  sender,
	}
}

func (s *Server) Listen(ctx context.Context) {
	for msg := range s.msgChan {
		order := &model.AddOrderRequest{
			ID:          msg.ID,
			Title:       msg.Title,
			Level:       msg.Level,
			Description: msg.Description,
			UpdatedAt:   msg.UpdatedAt,
			CreatedAt:   msg.CreatedAt,
		}

		var rm model.Message
		resp, err := s.service.Add(ctx, order)
		if err != nil {
			rm = model.Message{
				ID: order.ID,
			}
		} else {
			rm = model.Message{
				ID:          resp.ID,
				Title:       resp.Title,
				Level:       resp.Level,
				Description: resp.Description,
				UpdatedAt:   resp.UpdatedAt,
				CreatedAt:   resp.CreatedAt,
			}
		}

		err = s.sender.SendMessage(rm)
		if err != nil {
			//fmt.Println("App send sync message error: ", err)
			logger.Error("App send sync message error: ", err)
		}
	}
	//log.Println("App channel closed, exiting listen goroutine")
	logger.Info("App channel closed, exiting listen goroutine")
}
