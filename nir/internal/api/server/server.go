package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/newRational/soc/infrastructure/logger"
	"github.com/newRational/soc/infrastructure/redis"
	"github.com/newRational/soc/internal/model"
)

const QueryParamId = "id"

type service interface {
	Add(context.Context, *model.AddOrderRequest) error
}

type Server struct {
	service service
	cache   redis.Cache
}

func NewServer(service service, cache redis.Cache) Server {
	return Server{
		service: service,
		cache:   cache,
	}
}

var respChans = make(map[uint64]chan model.Message)
var mu sync.RWMutex

func (s *Server) Add(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	if string(body) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var order model.AddOrderRequest
	if err = json.Unmarshal(body, &order); err != nil {
		//fmt.Println(err)
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mu.Lock()
	respChans[order.ID] = make(chan model.Message)
	mu.Unlock()

	err = s.service.Add(req.Context(), &order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	var rm model.Message
	mu.RLock()
	select {
	case rm = <-respChans[order.ID]:
		break
	case <-time.After(5 * time.Second):
		mu.RUnlock()
		if rm.Title == "" {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Oops, something went wrong"))
			logger.Warn("The timeout has triggered")
			return
		}
	}
	mu.RUnlock()

	if rm.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This order has already been bought"))
		return
	}

	resp := &model.AddOrderResponse{
		ID:          rm.ID,
		Title:       rm.Title,
		Level:       rm.Level,
		Description: rm.Description,
		UpdatedAt:   rm.UpdatedAt,
		CreatedAt:   rm.CreatedAt,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (s *Server) Get(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)[QueryParamId]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := s.cache.Get(req.Context(), req.URL.Path)
	if err == nil {
		w.WriteHeader(resp.Status)
		w.Write(resp.Data)
		return
	}

	data, status := s.orderGet(req.Context(), idInt)

	err = s.cache.Set(req.Context(), req.URL.Path,
		&model.Response{
			Status: status,
			Data:   data,
		},
	)
	if err != nil {
		//log.Println(err)
		logger.Error(err)
	}

	w.WriteHeader(status)
	w.Write(data)
}

func (s *Server) orderGet(ctx context.Context, id uint64) ([]byte, int) {
	point, err := s.service.Get(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrObjectNotFound) {
			return []byte(err.Error()), http.StatusNotFound
		}
		return nil, http.StatusInternalServerError
	}

	pointJson, err := json.Marshal(point)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return pointJson, http.StatusOK
}
