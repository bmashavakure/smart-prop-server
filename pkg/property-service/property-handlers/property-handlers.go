package property_handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Brian-Mashavakure/smart-prop-server/pkg/database/connector"
	"github.com/Brian-Mashavakure/smart-prop-server/pkg/database/models"
	"github.com/Brian-Mashavakure/smart-prop-server/pkg/utils"
	"github.com/gin-gonic/gin"
)

type PrefRequest struct {
	EMAIL         string   `json:"email"`
	LOCATIONS     []string `json:"locations"`
	BUDGET        string   `json:"budget"`
	BEDROOMS      uint     `json:"bedrooms"`
	PROPERTY_SIZE uint     `json:"property_size"`
	AMENITIES     []string `json:"amenities"`
}

func GetPreferencesHandler(c *gin.Context) {
	var prefReq PrefRequest
	if err := c.ShouldBindJSON(&prefReq); err != nil {
		fmt.Printf("Error: %v\n", err)
		parseResponse := utils.ReturnJsonResponse("failed", "failed to bind json", nil, map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusBadRequest, parseResponse)
		return
	}

	//get user id
	var user models.User
	result := connector.DB.Where("email = ?", prefReq.EMAIL).First(&user)
	if result.Error != nil {
		log.Printf("Error occurred trying to find user:\n %n", result.Error)
		resultResponse := utils.ReturnJsonResponse("failed", "User not found", nil, map[string]interface{}{"error": "user could not be found in our system"})
		c.JSON(http.StatusBadRequest, resultResponse)
		return
	}

	//locations and amenities json stuff
	locationsJson, err := json.Marshal(prefReq.LOCATIONS)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ReturnJsonResponse(
			"failed",
			"could not encode locations",
			nil,
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	fmt.Printf("Locations: %s\n", string(locationsJson))

	amenitiesJson, err := json.Marshal(prefReq.AMENITIES)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ReturnJsonResponse(
			"failed",
			"could not encode amenities",
			nil,
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	fmt.Printf("Amenities: %s\n", string(amenitiesJson))

	//preference
	preference := models.Preferences{
		UserID:        user.ID,
		LOCATIONS:     locationsJson,
		BEDROOMS:      prefReq.BEDROOMS,
		PROPERTY_SIZE: prefReq.PROPERTY_SIZE,
		AMENITIES:     amenitiesJson,
		BUDGET:        prefReq.BUDGET,
	}

	create := connector.DB.Create(&preference)
	if create.Error != nil {
		c.JSON(http.StatusBadRequest, utils.ReturnJsonResponse(
			"failed",
			"preference could not be created",
			nil,
			map[string]interface{}{"error": err.Error()},
		))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, utils.ReturnJsonResponse(
		"success",
		"User preference saved",
		nil,
		nil,
	))
}
