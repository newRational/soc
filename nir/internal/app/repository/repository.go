package repository

import (
	"context"

	"github.com/newRational/soc/internal/app/db"
	"github.com/newRational/soc/internal/model"
)

type OrdersRepo struct {
	db db.DBops
}

func NewOrders(database db.DBops) *OrdersRepo {
	return &OrdersRepo{db: database}
}

func (r *OrdersRepo) Add(ctx context.Context, order *model.Order) error {
	var id uint64
	err := r.db.ExecQueryRow(ctx, "INSERT INTO orders(id, title, level, description, updated_at, created_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;",
		order.ID, order.Title, order.Level, order.Description, order.UpdatedAt, order.CreatedAt).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
