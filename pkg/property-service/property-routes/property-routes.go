package property_routes

import (
	"github.com/Brian-Mashavakure/smart-prop-server/pkg/auth-service/middleware"
	property_handlers "github.com/Brian-Mashavakure/smart-prop-server/pkg/property-service/property-handlers"
	"github.com/gin-gonic/gin"
)

func PropertyRoutes(router *gin.Engine) {
	api := router.Group("/smart-prop-api/prop/")

	api.POST("user-preferences", property_handlers.GetPreferencesHandler, middleware.JWTMiddleware())
	api.POST("get-properties", property_handlers.GetPropertiesHandler, middleware.JWTMiddleware())
	api.POST("create-booking", property_handlers.BookingHandler, middleware.JWTMiddleware())
	api.POST("get-bookings", property_handlers.GetBookingsHandler, middleware.JWTMiddleware())

}
