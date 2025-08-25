package postgres

import (
	"context"
)

func (s *Storage) UpsertPlatform(ctx context.Context, code, name string) (int, error) {
	const q = `INSERT INTO mp_platform(code,name)
	           VALUES ($1,$2)
	           ON CONFLICT(code) DO UPDATE SET name=EXCLUDED.name
	           RETURNING id;`
	var id int
	return id, s.pool.QueryRow(ctx, q, code, name).Scan(&id)
}

func (s *Storage) GetPlatformID(ctx context.Context, code string) (int, error) {
	const q = `SELECT id FROM mp_platform WHERE code=$1;`
	var id int
	if err := s.pool.QueryRow(ctx, q, code).Scan(&id); err != nil { return 0, err }
	return id, nil
}
