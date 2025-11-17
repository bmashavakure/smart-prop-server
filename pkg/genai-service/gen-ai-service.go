package genai_service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Brian-Mashavakure/smart-prop-server/pkg/database/models"
	"google.golang.org/genai"
)

func GenAiClient(prompt string) (string, error) {
	ctx := context.Background()
	aiClient, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", err
	}

	result, resultErr := aiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil)
	if resultErr != nil {
		return "", resultErr
	}

	return result.Text(), nil

}

// property recommendation stuff
func GetPropertyRecommendations(userPref models.Preferences, properties []models.Property) (string, error) {
	prefsJson, prefsErr := json.Marshal(userPref)
	if prefsErr != nil {
		return "", prefsErr
	}

	propertiesJson, propsErr := json.Marshal(properties)
	if propsErr != nil {
		return "", propsErr
	}

	//prompt build
	var prompt strings.Builder
	prompt.WriteString("You are an expert real estate recommendation agent.\n")
	prompt.WriteString("Your task is to analyze the provided user preferences and available properties and return the best 10 recommendations.\n")
	prompt.WriteString("The are in square in square feet of the property suggested should be equal to or higher than the area in the user preferences and the address of the property suggested should match the types of areas provided in the locations in proximity and type of area from the user preferences.\n")
	prompt.WriteString("In the user preferences pay close attention to the amenities and suggest properties that have the same amenities or amenities that are similar to those preferred by the user.\n")
	prompt.WriteString("Return only the property ids of the best properties in the order of best to least recommendation. Dont add any other text or content to the response just return the ids for the properties\n")

	prompt.WriteString("---User Preferences\n")
	prompt.WriteString(string(prefsJson))
	prompt.WriteString("\n\n")

	prompt.WriteString("---Available Properties---\n")
	prompt.WriteString(string(propertiesJson))

	response, responseErr := GenAiClient(prompt.String())
	if responseErr != nil {
		return "", responseErr
	}

	return response, nil
}
