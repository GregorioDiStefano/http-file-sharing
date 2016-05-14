package main

import (
	"github.com/GregorioDiStefano/go-file-storage/controllers"
	"github.com/GregorioDiStefano/go-file-storage/helpers"
	"github.com/GregorioDiStefano/go-file-storage/models"
	"github.com/gin-gonic/gin"
	"os"
)

func init() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		panic("No CONFIG_FILE set")
	}
	helpers.ParseConfig(configFile)
	models.DB.OpenDatabaseFile()
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	go deleteUnusedFile()

	router.GET("/", controller.IndexPage)
	router.PUT("/:filename", controller.SimpleUpload)
	router.GET("/:key/:filename", controller.DownloadFile)
	router.DELETE("/:key/:delete_key/:filename", controller.DeleteFile)
	router.Run(helpers.Config.ServerPort)
}
