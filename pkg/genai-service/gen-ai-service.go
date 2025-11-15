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
