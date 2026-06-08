package main

import (
	"fmt"
	"os"

	"sispak-dada/config"
	"sispak-dada/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
	if err != nil {
		fmt.Println("File .env tidak ditemukan, menggunakan konfigurasi default sistem")
	}

	config.ConnectDatabase()

	r := gin.Default()

	routes.SetupRoutes(r)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}