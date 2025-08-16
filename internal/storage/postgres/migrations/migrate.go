package migrate

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	db, err := goose.OpenDBWithDriver("postgres", pool.Config().ConnString())
	if err != nil {
		return err
	}
	defer db.Close()
	goose.SetBaseFS(nil)
	return goose.Up(db, dir)
}
