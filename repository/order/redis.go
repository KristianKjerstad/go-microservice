package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"project/model"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	//add to redis database by converting struct to a json string
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Failed to encode order %w", err)
	}
	key := orderIDKey(order.OrderID)

	insertTransaction := r.Client.TxPipeline()

	res := insertTransaction.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		insertTransaction.Discard()
		return fmt.Errorf("Failed to encode set %w", err)
	}

	if err := insertTransaction.SAdd(ctx, "orders", key).Err(); err != nil {
		insertTransaction.Discard()
		return fmt.Errorf("Failed to add orders set: %w", err)
	}

	if _, err := insertTransaction.Exec(ctx); err != nil {
		return fmt.Errorf("Failed to exec: %w", err)
	}
	return nil
}

var ErrNotExist = errors.New("Order does not exist")

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIDKey(id)
	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("Get order: %w", err)
	}

	var order model.Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("Failed to decode order json: %w", err)
	}

	return order, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := orderIDKey(id)

	deleteTransaction := r.Client.TxPipeline()

	err := deleteTransaction.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		deleteTransaction.Discard()
		return ErrNotExist
	} else if err != nil {
		deleteTransaction.Discard()
		return fmt.Errorf("Get order: %w", err)
	}

	if err := deleteTransaction.SRem(ctx, "orders", key).Err(); err != nil {
		deleteTransaction.Discard()
		return fmt.Errorf("Failed to remove from orders set: %w", err)
	}

	if _, err := deleteTransaction.Exec(ctx); err != nil {
		return fmt.Errorf("Failed to exec: %w", err)
	}
	return nil

}

func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("Failed to encode order: %w", err)
	}

	key := orderIDKey(order.OrderID)

	err = r.Client.SetXX(ctx, key, srting(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("Set order: %w", err)
	}
	return nil

}

type FindAllPage struct {
	Size   uint
	Offset uint
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, paeg FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "orders", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("Failed to get order ids: %w", err)
	}

	results, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("Failed to get orders: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{Orders: []model.Order{}}, nil
	}

	orders := make([]model.Order, len(results))

	for i, result := range results {
		result := result.(string)
		var order model.Order

		err := json.Unmarshal([]byte(result), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("Failed to decode order json: %w", err)
		}
		orders[i] = order
	}

	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil

}
