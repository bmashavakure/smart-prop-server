package property_handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
	PROPERTY_SIZE float64  `json:"property_size"`
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

	// Use WaitGroup to fetch properties and preferences concurrently
	var wg sync.WaitGroup
	var properties []models.Property
	var userPref models.Preferences
	var propertyErr, prefErr error

	wg.Add(2)

	// Fetch properties in a goroutine
	go func() {
		defer wg.Done()
		result := connector.DB.Find(&properties)
		if result.Error != nil {
			propertyErr = result.Error
		}
	}()

	// Fetch user preferences in a goroutine
	go func() {
		defer wg.Done()
		result := connector.DB.Where("user_id = ?", userID).First(&userPref)
		if result.Error != nil {
			prefErr = result.Error
		}
	}()

	// Wait for both queries to complete
	wg.Wait()

	// Check for errors after both queries complete
	if propertyErr != nil {
		log.Printf("Error occurred trying to find properties:\n %v", propertyErr)
		c.JSON(http.StatusNotFound, utils.ReturnJsonResponse("failed", "Something went wrong", nil, map[string]interface{}{"error": propertyErr.Error()}))
		return
	}

	if prefErr != nil {
		log.Printf("Error occurred trying to find user preferences:\n %v", prefErr)
		c.JSON(http.StatusNotFound, utils.ReturnJsonResponse("failed", "Something went wrong finding preference", nil, map[string]interface{}{"error": prefErr.Error()}))
		return
	}

	// Get AI recommendations
	recommendation, recErr := genai_service.GetPropertyRecommendations(userPref, properties)
	if recErr != nil {
		log.Printf("Error occurred trying to find recommendation:\n %v", recErr)
		c.JSON(http.StatusNotFound, utils.ReturnJsonResponse("failed", "Property recommendation failed", nil, map[string]interface{}{"error": recErr.Error()}))
		return
	}

	fmt.Printf("Recommendations: %s\n", recommendation)

	idsList, idsError := utils.NumbersSeparator(recommendation)
	if idsError != nil {
		log.Printf("Error occurred trying to parse recommendation ids:\n %v", idsError)
		c.JSON(http.StatusNotFound, utils.ReturnJsonResponse("failed", "Property recommendation failed", nil, map[string]interface{}{"error": idsError.Error()}))
		return
	}

	finalProperties := utils.FilterProperties(idsList, properties)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, utils.ReturnJsonResponse("success", "Properties found", map[string]interface{}{"properties": finalProperties}, nil))
}

type BookingReq struct {
	PropertyID   uint   `json:"property_id"`
	BookingDate  string `json:"booking_date"`
	BookingTime  string `json:"booking_time"`
	CheckoutDate string `json:"checkout_date"`
	CheckoutTime string `json:"checkout_time"`
	UserID       uint   `json:"user_id"`
}

func BookingHandler(c *gin.Context) {
	var req BookingReq

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Error: %v\n", err)
		parseResponse := utils.ReturnJsonResponse("failed", "failed to bind json", nil, map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusBadRequest, parseResponse)
		return
	}

	// Parse and validate dates
	bookingDate, bookingDateErr := time.Parse("2006-01-02", req.BookingDate)
	if bookingDateErr != nil {
		log.Printf("Invalid booking date format: %v\n", bookingDateErr)
		dateFormatResponse := utils.ReturnJsonResponse("failed", "invalid booking date format", nil, map[string]interface{}{"error": "booking date must be in YYYY-MM-DD format"})
		c.JSON(http.StatusBadRequest, dateFormatResponse)
		return
	}

	checkoutDate, checkoutDateErr := time.Parse("2006-01-02", req.CheckoutDate)
	if checkoutDateErr != nil {
		log.Printf("Invalid checkout date format: %v\n", checkoutDateErr)
		dateFormatResponse := utils.ReturnJsonResponse("failed", "invalid checkout date format", nil, map[string]interface{}{"error": "checkout date must be in YYYY-MM-DD format"})
		c.JSON(http.StatusBadRequest, dateFormatResponse)
		return
	}

	// Check if checkout date is before booking date
	if checkoutDate.Before(bookingDate) {
		log.Printf("Checkout date (%s) is before booking date (%s)\n", req.CheckoutDate, req.BookingDate)
		dateOrderResponse := utils.ReturnJsonResponse("failed", "invalid date range", nil, map[string]interface{}{"error": "checkout date cannot be before booking date"})
		c.JSON(http.StatusBadRequest, dateOrderResponse)
		return
	}

	// Check if there's already a booking for this property on the same date
	var existingBooking models.Booking
	bookingCheck := connector.DB.Where("property_id = ? AND booking_date = ?", req.PropertyID, req.BookingDate).First(&existingBooking)
	if bookingCheck.Error == nil {
		log.Printf("Property already booked on this date: Property ID %d, Date %s\n", req.PropertyID, req.BookingDate)
		conflictResponse := utils.ReturnJsonResponse("failed", "property already booked on that day", nil, map[string]interface{}{"error": "this property is already booked for the selected date"})
		c.JSON(http.StatusConflict, conflictResponse)
		return
	}

	// Create the booking
	booking := models.Booking{
		PropertyID:   req.PropertyID,
		BookingDate:  req.BookingDate,
		BookingTime:  req.BookingTime,
		CheckoutDate: req.CheckoutDate,
		CheckoutTime: req.CheckoutTime,
		UserID:       req.UserID,
	}

	result := connector.DB.Create(&booking)
	if result.Error != nil {
		log.Printf("Error occurred trying to create booking:\n %v", result.Error)
		createError := utils.ReturnJsonResponse("failed", "failed to create booking", nil, map[string]interface{}{"error": result.Error.Error()})
		c.JSON(http.StatusInternalServerError, createError)
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, utils.ReturnJsonResponse("success", "Booking created successfully", map[string]interface{}{"booking_id": booking.ID}, nil))
}

func GetBookingsHandler(c *gin.Context) {
	userID := c.Request.FormValue("user_id")

	if userID == "" {
		log.Println("user_id parameter is missing")
		missingParamResponse := utils.ReturnJsonResponse("failed", "user_id is required", nil, map[string]interface{}{"error": "user_id parameter is missing"})
		c.JSON(http.StatusBadRequest, missingParamResponse)
		return
	}

	var bookings []models.Booking
	bookingsResult := connector.DB.Where("user_id = ?", userID).Find(&bookings)
	if bookingsResult.Error != nil {
		log.Printf("Error occurred trying to find bookings:\n %v", bookingsResult.Error)
		c.JSON(http.StatusInternalServerError, utils.ReturnJsonResponse("failed", "failed to retrieve bookings", nil, map[string]interface{}{"error": bookingsResult.Error.Error()}))
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, utils.ReturnJsonResponse("success", "Bookings retrieved successfully", map[string]interface{}{"bookings": bookings, "count": len(bookings)}, nil))
}
