package main

import (
	"github.com/gin-gonic/gin"
	"go-asset/service"
)

func main() {
	r := gin.Default()

	r.POST("/seats", service.PostSeatsCsv)
	r.GET("/seats/:persent", service.GetSeatsByPersent)
	r.GET("/image/:key", service.GetBestPathImage)
	//r.Run("0.0.0.0:8080")
	r.RunTLS("0.0.0.0:8080", "cmd/ssl.crt", "cmd/ssl.key")
}