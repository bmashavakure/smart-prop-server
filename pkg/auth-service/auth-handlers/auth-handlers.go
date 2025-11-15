package auth_handlers

import (
	"log"
	"net/http"
	"time"

	auth_utils "github.com/Brian-Mashavakure/smart-prop-server/pkg/auth-service/auth-utils"
	"github.com/Brian-Mashavakure/smart-prop-server/pkg/database/connector"
	"github.com/Brian-Mashavakure/smart-prop-server/pkg/database/models"
	"github.com/Brian-Mashavakure/smart-prop-server/pkg/utils"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	name := c.Request.FormValue("name")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")

	hashedPassword := auth_utils.HashPassword(string(password))

	errorMessage := "Failed To Register User"

	var existingUser models.User

	oldEmail := connector.DB.Where("email = ?", email).First(&existingUser)
	if oldEmail.Error == nil {
		log.Printf("User Already Exists:\n %n", oldEmail.Error)
		oldEmailResponse := utils.ReturnJsonResponse("failed", errorMessage, nil, map[string]interface{}{"error": "user already exists"})
		c.JSON(http.StatusOK, oldEmailResponse)
		return
	}

	user := models.User{
		NAME:     name,
		EMAIL:    email,
		PASSWORD: hashedPassword,
	}

	result := connector.DB.Create(&user)
	if result.Error != nil {
		log.Printf("Error occurred trying to create user:\n %n", result.Error)
		createError := utils.ReturnJsonResponse("failed", errorMessage, nil, map[string]interface{}{"error": "failed to create user in db"})
		c.JSON(http.StatusInternalServerError, createError)
		return
	}

	currentDateTime := time.Now().Format("20060102150405")
	tokenString, tokenErr := auth_utils.GenerateJWTToken(currentDateTime, email)
	if tokenErr != nil {
		log.Printf("Error occurred trying to create user:\n %n", tokenErr)
		tokenError := utils.ReturnJsonResponse("failed", errorMessage, nil, map[string]interface{}{"error": "failed to generate token"})
		c.JSON(http.StatusInternalServerError, tokenError)
		return
	}

	finalResponse := utils.ReturnJsonResponse("success", "User created successfully", map[string]interface{}{"id": user.ID, "token": tokenString}, nil)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, finalResponse)
}

func LoginHandler(c *gin.Context) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")

	var user models.User
	result := connector.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		log.Printf("Error occurred trying to find user:\n %n", result.Error)
		resultResponse := utils.ReturnJsonResponse("failed", "User not found", nil, map[string]interface{}{"error": "user could not be found in our system"})
		c.JSON(http.StatusBadRequest, resultResponse)
		return
	}

	//compare passwords
	comparison := auth_utils.ComparePasswordAndHash(user.PASSWORD, password)
	if comparison == false {
		log.Println("Passwords do not match")
		passwordResponse := utils.ReturnJsonResponse("failed", "password mismatch", nil, map[string]interface{}{"error": "incorrect password provided"})
		c.JSON(http.StatusBadRequest, passwordResponse)
		return
	}

	currentDateTime := time.Now().Format("20060102150405")
	tokenString, err := auth_utils.GenerateJWTToken(currentDateTime, user.EMAIL)
	if err != nil {
		log.Println("Failed to generate token")
		tokenResponse := utils.ReturnJsonResponse("failed", "bad request", nil, map[string]interface{}{"error": "something went wrong"})
		c.JSON(http.StatusBadRequest, tokenResponse)
		return
	}

	finalResponse := utils.ReturnJsonResponse("success", "login successful", map[string]interface{}{"message": "login successful", "token": tokenString}, nil)
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, finalResponse)

}
