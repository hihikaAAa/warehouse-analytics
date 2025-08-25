package postgres

import "context"

func (s *Storage) GetCursor(ctx context.Context, platformID int, resource string) (string, bool, error) {
	const q = `SELECT cursor_value FROM etl_cursors WHERE platform_id=$1 AND resource=$2;`
	var val string
	if err := s.pool.QueryRow(ctx, q, platformID, resource).Scan(&val); err != nil {
		if err.Error() == "no rows in result set" { return "", false, nil }
		return "", false, err
	}
	return val, true, nil
}

func (s *Storage) SetCursor(ctx context.Context, platformID int, resource, value string) error {
	const q = `INSERT INTO etl_cursors(platform_id,resource,cursor_value)
	           VALUES ($1,$2,$3)
	           ON CONFLICT(platform_id,resource) DO UPDATE
	             SET cursor_value=EXCLUDED.cursor_value, updated_at=now();`
	_, err := s.pool.Exec(ctx, q, platformID, resource, value)
	return err
}
