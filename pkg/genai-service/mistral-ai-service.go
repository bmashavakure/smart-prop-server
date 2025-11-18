package genai_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type MistralChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MistralChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func MistralHandler(promptString string) (string, error) {
	//err := godotenv.Load()
	//if err != nil {
	//	fmt.Println("Error loading .env file")
	//}

	apiKey := os.Getenv("MISTRAL_API_KEY")
	url := "https://api.mistral.ai/v1/chat/completions"

	// mistral request
	reqData := MistralChatRequest{
		Model:       "mistral-small",
		Temperature: 0.2,
		Messages: []Message{
			{Role: "user", Content: promptString},
		},
	}

	body, marshalErr := json.Marshal(reqData)
	if marshalErr != nil {
		fmt.Printf("Error marshalling request body: %s\n", marshalErr.Error())
		return "", fmt.Errorf("failed to marshal request: %w", marshalErr)
	}

	// request
	client := &http.Client{Timeout: 12 * time.Second}
	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if reqErr != nil {
		fmt.Printf("Error creating request: %s\n", reqErr.Error())
		return "", fmt.Errorf("failed to create request: %w", reqErr)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, respErr := client.Do(req)
	if respErr != nil {
		fmt.Printf("Error executing request: %s\n", respErr.Error())
		return "", fmt.Errorf("failed to execute request: %w", respErr)
	}
	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Printf("Error reading response: %s\n", readErr.Error())
		return "", fmt.Errorf("failed to read response body: %w", readErr)
	}

	var mistralResp MistralChatResponse
	unmarshalErr := json.Unmarshal(respBody, &mistralResp)
	if unmarshalErr != nil {
		fmt.Printf("Error unmarshalling response: %s\n", unmarshalErr.Error())
		return "", fmt.Errorf("failed to unmarshal response: %w", unmarshalErr)
	}

	if len(mistralResp.Choices) == 0 {
		fmt.Printf("No choices found in response: %s\n", string(respBody))
		return "", fmt.Errorf("no choices returned from Mistral API")
	}

	content := mistralResp.Choices[0].Message.Content
	fmt.Printf("Mistral AI Response: %s\n", content)

	return content, nil
}
