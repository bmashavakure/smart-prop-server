package main

import (
	auth_routes "github.com/Brian-Mashavakure/smart-prop-server/pkg/auth-service/auth-routes"
	"github.com/Brian-Mashavakure/smart-prop-server/pkg/database/connector"
	"github.com/Brian-Mashavakure/smart-prop-server/pkg/database/models"
	property_routes "github.com/Brian-Mashavakure/smart-prop-server/pkg/property-service/property-routes"
	"github.com/gin-gonic/gin"
)

func main() {
	connector.Connector()

	connector.DB.AutoMigrate(models.User{}, models.Preferences{}, models.Property{}, models.Booking{})

	router := gin.Default()
	auth_routes.AuthRoutes(router)
	property_routes.PropertyRoutes(router)

	router.Run("localhost:8080")
}
