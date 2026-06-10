package main

import (
	"fmt"
	"os"
	"path/filepath"

	"sispak-dada/config"
	"sispak-dada/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("File .env tidak ditemukan, menggunakan konfigurasi default sistem")
	}

	config.ConnectDatabase()

	r := gin.Default()

	loadTemplates(r)
	r.Static("/static", "./static")

	routes.SetupRoutes(r)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}

func loadTemplates(r *gin.Engine) {
	var templateFiles []string

	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".html" {
			templateFiles = append(templateFiles, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	r.LoadHTMLFiles(templateFiles...)
}