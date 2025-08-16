package main 

import(
	"github.com/joho/godotenv"

	"github.com/hihikaAAa/warehouse-analytics/internal/config"
)

const(
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)
func main(){
	_ = godotenv.Load("local.env")
	cfg := config.MustLoad()
	
	//TODO: logger
	//TODO: init storage
	//TODO: init router
	//TODO: run server
}