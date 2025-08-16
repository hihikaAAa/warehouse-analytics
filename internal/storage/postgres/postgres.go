package postgres

import(
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/hihikaAAa/warehouse-analytics/internal/storage/postgres/migrations"
)

type Storage struct{
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string)(*Storage, error){
	const op = "storage.postgres.New"
	
	cfg,err := pgxpool.ParseConfig(dsn)
	if err !=nil{
		return nil, fmt.Errorf("%s: %w",op, err)
	}
	cfg.MaxConns = 12
	cfg.HealthCheckPeriod = 30 * time.Second
	cfg.MaxConnLifetime = 2 * time.Hour

	pool, err := pgxpool.NewWithConfig(ctx,cfg)
	if err != nil{
		return nil, fmt.Errorf("%s : %w",op,err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w",op, err)
	}

	if err := migrate.RunMigrations(ctx, pool, "migrations"); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s, %w", op,err)
	}
	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() { 
	s.pool.Close() 
}

func (s *Storage) Pool() *pgxpool.Pool { 
	return s.pool 
}