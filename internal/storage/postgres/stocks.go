package postgres

import "context"

func (s *Storage) UpsertStock(ctx context.Context, listingID int, warehouse string, qty int, updatedAt string) error {
	const q = `INSERT INTO stock(listing_id, warehouse, qty, updated_at)
	           VALUES ($1,$2,$3,$4)
	           ON CONFLICT(listing_id, warehouse) DO UPDATE
	             SET qty=EXCLUDED.qty,
	                 updated_at = GREATEST(stock.updated_at, EXCLUDED.updated_at);`
	_, err := s.pool.Exec(ctx, q, listingID, warehouse, qty, updatedAt)
	return err
}
