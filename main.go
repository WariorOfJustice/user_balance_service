package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"user_balance_service/app/core"

	"context"
	"os"

	"user_balance_service/app/handlers"
)

func init() {
	if err := godotenv.Load(); err != nil {
		println("Отсутствует .env файл")
		os.Exit(1)
	}
    println("Переменные окружения загружены")
}

func main() {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("USER_NAME"),
		os.Getenv("USER_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
        println(err.Error())
		println("Не могу подключиться к базе, help")
		os.Exit(1)
	}
    println("Подключились к базе")
	defer conn.Close(context.Background())
	core.SetDb(conn)

	engine := gin.Default()
	addRoutes(engine)
	engine.Run()
}

func addRoutes(engine *gin.Engine) {
	engine.GET("get_user_balance", handlers.GetUserBalance)
	engine.POST("add_user_balance", handlers.AddUserBalance)
	engine.POST("reserve_user_balance", handlers.ReserveUserBalance)
	engine.POST("add_revenue_to_company", handlers.AddRevenueToCompany)
}
