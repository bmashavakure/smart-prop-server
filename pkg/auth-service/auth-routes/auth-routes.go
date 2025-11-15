package auth_routes

import (
	auth_handlers "github.com/Brian-Mashavakure/smart-prop-server/pkg/auth-service/auth-handlers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	api := router.Group("/smart-prop-api/auth/")

	api.POST("register-user", auth_handlers.RegisterHandler)
	api.POST("login-user", auth_handlers.LoginHandler)
}
