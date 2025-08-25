package postgres

import "context"

func (s *Storage) UpsertProduct(ctx context.Context, sku, title, barcode, category string) (int, error) {
	const q = `INSERT INTO product(sku,title,barcode,category)
	           VALUES ($1,$2,$3,$4)
	           ON CONFLICT(sku) DO UPDATE
	             SET title=EXCLUDED.title, barcode=EXCLUDED.barcode, category=EXCLUDED.category
	           RETURNING id;`
	var id int
	return id, s.pool.QueryRow(ctx, q, sku, title, barcode, category).Scan(&id)
}
