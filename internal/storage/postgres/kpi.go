package postgres

import "context"

func (s *Storage) RefreshSalesDaily(ctx context.Context, from, to string) (int64, error) {
	const q = `INSERT INTO kpi_sales_daily(day, listing_id, orders, units, revenue, returns)
	           SELECT date(o.created_at), oi.listing_id,
	                  count(DISTINCT o.id), sum(oi.qty), sum(oi.revenue), 0
	           FROM orders o
	           JOIN order_items oi ON oi.order_id=o.id
	           WHERE o.created_at >= $1 AND o.created_at < $2
	           GROUP BY 1,2
	           ON CONFLICT(day, listing_id) DO UPDATE
	             SET orders=EXCLUDED.orders, units=EXCLUDED.units, revenue=EXCLUDED.revenue;`
	ct, err := s.pool.Exec(ctx, q, from, to)
	return ct.RowsAffected(), err
}
