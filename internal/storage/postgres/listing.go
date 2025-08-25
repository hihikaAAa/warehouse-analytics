package postgres

import "context"

func (s *Storage) UpsertListing(ctx context.Context, productID, platformID int, mpSKU string) (int, error) {
	const q = `INSERT INTO listing(product_id,platform_id,mp_sku)
	           VALUES ($1,$2,$3)
	           ON CONFLICT(product_id,platform_id) DO UPDATE SET mp_sku=EXCLUDED.mp_sku
	           RETURNING id;`
	var id int
	return id, s.pool.QueryRow(ctx, q, productID, platformID, mpSKU).Scan(&id)
}

func (s *Storage) GetListingIDByPlatformSKU(ctx context.Context, platformID int, mpSKU string) (int, error) {
	const q = `SELECT id FROM listing WHERE platform_id=$1 AND mp_sku=$2;`
	var id int
	return id, s.pool.QueryRow(ctx, q, platformID, mpSKU).Scan(&id)
}
