package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"resuming/api"
	"resuming/database"
	"resuming/scheduler"
	systemconfig "resuming/system-config"
	"resuming/tool"
)

func main() {
	// set this during production
	// gin.SetMode(gin.ReleaseMode)

	if err := tool.SetupValkey(); err != nil {
		log.Fatalln("Failed to setup Valkey:", err)
	}
	if err := database.DatabaseConnect(); err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	if err := database.TableConnect(); err != nil {
		log.Fatalln("Failed to migrate database tables:", err)
	}

	scheduler.FirstSync()
	go scheduler.FullSync()

	// let the server settle before first sync
	time.Sleep(5 * time.Second)

	router := api.APIConnect()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{systemconfig.FrontendUri},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	log.Println("Backend User Running at PORT 5321")
	if err := router.Run(":" + systemconfig.BackendPort); err != nil {
		log.Println("Failed to start server:", err)
	}
}
