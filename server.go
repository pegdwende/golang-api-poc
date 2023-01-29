package main

import (
	"fmt"
	"gitlab/pragmaticreviews/golang-gin-poc/controller"
	"gitlab/pragmaticreviews/golang-gin-poc/middlewares"
	"gitlab/pragmaticreviews/golang-gin-poc/service"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	//gindump "github.com/tpkeeper/gin-dump"
)

var (
	videoService    service.VideoService       = service.New()
	videoController controller.VideoController = controller.New(videoService)
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

}

func main() {

	setupLogOutput()
	server := gin.New()

	server.Static("/css", "./templates/css")
	server.LoadHTMLGlob("templates/*.html")

	server.Use(gin.Recovery(), middlewares.Looger()) //, gindump.Dump()

	apiRoutes := server.Group("/api", middlewares.BasicAuth())
	{
		apiRoutes.GET("/videos", func(ctx *gin.Context) {
			fmt.Println("printing videos")
			fmt.Println(videoController.FindAll())
			ctx.JSON(200, videoController.FindAll())
		})

		apiRoutes.POST("/video", func(ctx *gin.Context) {

			fmt.Println("saving videos")

			err := videoController.Save(ctx)

			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Video Input is Valid!!!"})
			}
		})
	}

	viewRoutes := server.Group("/view")
	{
		viewRoutes.GET("/videos", videoController.ShowAll)
	}

	server.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "ok!!",
		})
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	server.Run(":" + port)

}
