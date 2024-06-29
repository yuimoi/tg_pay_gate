package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"tg_pay_gate/internal/models"
	"tg_pay_gate/internal/router"
	"tg_pay_gate/internal/utils/config"
	"tg_pay_gate/internal/utils/db"
	"tg_pay_gate/internal/utils/tg_bot"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	port := flag.String("port", "8086", "http运行端口")
	flag.Parse()

	config.LoadAllConfig()
	db.InitDB()

	tg_bot.InitTGBot()
	go router.RunTgBot()

	// 数据库迁移，不能放在db.go因为models有一个applyFilter需要用到db.DB，考虑分理出models
	if err := db.DB.AutoMigrate(
		&models.Order{},
	); err != nil {
		panic(err)
	}
	r := router.SetupGinRoutes()

	runningPath := fmt.Sprintf("0.0.0.0:%s", *port)
	fmt.Printf("http将运行在端口: %s\n", runningPath)

	r.Run(runningPath)
}
