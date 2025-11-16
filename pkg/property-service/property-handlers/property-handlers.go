package property_handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Brian-Mashavakure/smart-prop-server/pkg/database/connector"
	"github.com/Brian-Mashavakure/smart-prop-server/pkg/database/models"
	genai_service "github.com/Brian-Mashavakure/smart-prop-server/pkg/genai-service"
	"github.com/Brian-Mashavakure/smart-prop-server/pkg/utils"
	"github.com/gin-gonic/gin"
)

type PrefRequest struct {
	USERID        uint     `json:"user_id"`
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
		UserID:        prefReq.USERID,
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

func GetPropertiesHandler(c *gin.Context) {
	userID := c.Request.FormValue("user_id")

	var properties []models.Property
	propertyResult := connector.DB.Find(&properties)
	if propertyResult.Error != nil {
		log.Printf("Error occurred trying to find properties:\n %n", propertyResult.Error)
		c.JSON(http.StatusNotFound, utils.ReturnJsonResponse("failed", "Something went wrong", nil, map[string]interface{}{"error": propertyResult.Error.Error()}))
		return
	}

	var userPref models.Preferences
	prefResult := connector.DB.Where("user_id = ?", userID).First(&userPref)
	if prefResult.Error != nil {
		log.Printf("Error occurred trying to find user:\n %n", propertyResult.Error)
		c.JSON(http.StatusNotFound, utils.ReturnJsonResponse("failed", "Something went wrong finding preference", nil, map[string]interface{}{"error": prefResult.Error.Error()}))
		return
	}

	recommendation, recErr := genai_service.GetPropertyRecommendations(userPref, properties)
	if recErr != nil {
		log.Printf("Error occurred trying to find recomendation:\n %n", propertyResult.Error)
		c.JSON(http.StatusNotFound, utils.ReturnJsonResponse("failed", "Property recommendation failed", nil, map[string]interface{}{"error": recErr.Error()}))
		return
	}
	idsList, idsError := utils.NumbersSeparator(recommendation)
	if idsError != nil {
		log.Printf("Error occurred trying to find recomendation:\n %n", propertyResult.Error)
		c.JSON(http.StatusNotFound, utils.ReturnJsonResponse("failed", "Property recommendation failed", nil, map[string]interface{}{"error": recErr.Error()}))
		return
	}
	finalProperties := utils.FilterProperties(idsList, properties)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, utils.ReturnJsonResponse("success", "Properties found", map[string]interface{}{"properties": finalProperties}, nil))
}
