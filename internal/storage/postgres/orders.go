package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) UpsertOrder(ctx context.Context, platformID int, mpOrderID, status, region string, createdAt string) (int64, error) {
	const q = `INSERT INTO orders(platform_id, mp_order_id, created_at, status, buyer_region)
	           VALUES ($1,$2,$3,$4,$5)
	           ON CONFLICT(platform_id, mp_order_id) DO UPDATE
	             SET created_at=EXCLUDED.created_at,
	                 status=EXCLUDED.status,
	                 buyer_region=EXCLUDED.buyer_region
	           RETURNING id;`
	var id int64
	return id, s.pool.QueryRow(ctx, q, platformID, mpOrderID, createdAt, status, region).Scan(&id)
}

type OrderItemUpsert struct {
	ListingID int
	Qty       int
	Price     float64
	Discount  float64
}

func (s *Storage) UpsertOrderItems(ctx context.Context, orderID int64, items []OrderItemUpsert) error {
	tx, err := s.pool.Begin(ctx); if err != nil { return err }
	defer func() { _ = tx.Rollback(ctx) }()

	b := &pgx.Batch{}
	for _, it := range items {
		const q = `INSERT INTO order_items(order_id, listing_id, qty, price, discount)
		           VALUES ($1,$2,$3,$4,$5)
		           ON CONFLICT(order_id, listing_id) DO UPDATE
		             SET qty=EXCLUDED.qty, price=EXCLUDED.price, discount=EXCLUDED.discount;`
		b.Queue(q, orderID, it.ListingID, it.Qty, it.Price, it.Discount)
	}
	br := tx.SendBatch(ctx, b)
	if err := br.Close(); err != nil { return err }
	return tx.Commit(ctx)
}
