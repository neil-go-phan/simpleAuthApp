package main

import (
	"example/web-service-gin/controllers"
	"example/web-service-gin/database"
	"example/web-service-gin/models"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	r := SetupRouter()
	_ = r.Run(":8080")
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))
	db := database.GetDB()
	repository := models.CreateRepository(db)
	server := controllers.NewServer(repository)
	// r.Use(middlewares.JSONAppErrorReporter())
	r.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})	

	authRoutes := r.Group("auth")
	{
		authRoutes.POST("sign-up", server.SignUp)
		authRoutes.POST("sign-in", server.SignIn)
		authRoutes.GET("users", server.GetUsers)
	}


	return r
}



