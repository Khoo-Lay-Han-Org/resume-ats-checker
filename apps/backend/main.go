package main

import (
	"log"
	"time"

	echomw "github.com/labstack/echo/v4/middleware"
	"resuming/api"
	"resuming/database"
	"resuming/scheduler"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func main() {
	if err := tool.SetupValkey(); err != nil {
		log.Fatalln("Failed to setup Valkey:", err)
	}
	if err := database.DatabaseConnect(); err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	scheduler.FirstSync()
	go scheduler.FullSync()

	time.Sleep(5 * time.Second)

	router := api.APIConnect()
	router.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins:     []string{systemconfig.FrontendUri},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour / time.Second),
	}))

	log.Println("Backend User Running at PORT 5321")
	if err := router.Start(":" + systemconfig.BackendPort); err != nil {
		log.Fatalln("Failed to start server:", err)
	}
}
