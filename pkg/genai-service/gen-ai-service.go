package genai_service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Brian-Mashavakure/smart-prop-server/pkg/database/models"
	"google.golang.org/genai"
)

func GenAiClient(prompt string) (string, error) {
	ctx := context.Background()
	aiClient, err := genai.NewClient(ctx, nil)
	if err != nil {
		fmt.Printf("Gen Ai Setup Client Error: %v\n", err.Error)
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
	fmt.Printf("Prompt Response: %s\n", result.Text())

	return result.Text(), nil

}

// property recommendation stuff
func GetPropertyRecommendations(userPref models.Preferences, properties []models.Property) (string, error) {
	// Use WaitGroup to marshal JSON concurrently
	var wg sync.WaitGroup
	var prefsJson, propertiesJson []byte
	var prefsErr, propsErr error

	wg.Add(2)

	// Marshal user preferences in a goroutine
	go func() {
		defer wg.Done()
		prefsJson, prefsErr = json.Marshal(userPref)
	}()

	// Marshal properties in a goroutine
	go func() {
		defer wg.Done()
		propertiesJson, propsErr = json.Marshal(properties)
	}()

	// Wait for both marshaling operations to complete
	wg.Wait()

	// Check for errors
	if prefsErr != nil {
		return "", prefsErr
	}
	if propsErr != nil {
		return "", propsErr
	}

	// Build prompt
	var prompt strings.Builder
	prompt.WriteString("CRITICAL OUTPUT RULE: Your response must contain ONLY an array of property IDs. Do not include any explanations, descriptions, introductions, or concluding remarks. Output format must be exactly: [id1, id2, id3, ...]\n\n")
	prompt.WriteString("Task: Analyze the provided user preferences and available properties to recommend the best 10 properties.\n")
	prompt.WriteString("Criteria:\n")
	prompt.WriteString("- Property size must be equal to or greater than the user's preferred size\n")
	prompt.WriteString("- Location must match the types of areas specified in user preferences\n")
	prompt.WriteString("- Amenities should closely match or be similar to those preferred by the user\n\n")

	prompt.WriteString("---User Preferences---\n")
	prompt.WriteString(string(prefsJson))
	prompt.WriteString("\n\n")

	prompt.WriteString("---Available Properties---\n")
	prompt.WriteString(string(propertiesJson))
	prompt.WriteString("\n\n")

	prompt.WriteString("RESPONSE FORMAT (MANDATORY): Return ONLY the array of property IDs with no other text. Example: [23, 45, 67, 78, 65, 34, 21, 18, 90, 12]")

	//response, responseErr := GenAiClient(prompt.String())
	//if responseErr != nil {
	//	return "", responseErr
	//}

	response, responseErr := MistralHandler(prompt.String())
	if responseErr != nil {
		return "", responseErr
	}

	return response, nil
}
