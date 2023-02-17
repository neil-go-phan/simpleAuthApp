package main

import (
	"example/web-service-gin/controllers"
	// "example/web-service-gin/middlewares"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
	_ = r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))
	// r.Use(middlewares.JSONAppErrorReporter())
	r.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	userRepo := controllers.New()

	authRoutes := r.Group("auth")
	{
		authRoutes.POST("sign-up", userRepo.SignUp)
		authRoutes.POST("sign-in", userRepo.SignIn)
		authRoutes.GET("users", userRepo.GetUsers)
	}
	
	// r.PUT("/users/:id", userRepo.UpdateUser)
	// r.DELETE("/users/:id", userRepo.DeleteUser)

	return r
}



