package main 

import(
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/hihikaAAa/warehouse-analytics/internal/config"
	"github.com/hihikaAAa/warehouse-analytics/internal/lib/logger/handlers/slogpretty"
	"github.com/hihikaAAa/warehouse-analytics/internal/storage/postgres"
	"github.com/hihikaAAa/warehouse-analytics/internal/lib/logger/sl"
)

const(
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main(){
	_ = godotenv.Load("local.env")
	cfg := config.MustLoad()
	
	log := setupLogger(cfg.Env)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	storage,err := postgres.New(ctx,cfg.DB.DSN)
	if err != nil{
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}else{
		log.Info("Storage initialized")
	}
	defer storage.Close()
	_ = storage
	
	// TODO: init router
	// TODO: run server

	<-ctx.Done()
	log.Info("shutting down...")
	time.Sleep(100 * time.Millisecond)
}

func setupLogger(env string) *slog.Logger{
	var log *slog.Logger

	switch env{
		case envLocal:
			log = setupPrettySlog()
		case envDev:
			log = slog.New(slog.NewJSONHandler(os.Stdout,&slog.HandlerOptions{Level: slog.LevelDebug}),)
		case envProd:
			log = slog.New(slog.NewJSONHandler(os.Stdout,&slog.HandlerOptions{Level: slog.LevelInfo}),)
		}
		return log
}

func setupPrettySlog() *slog.Logger{
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}