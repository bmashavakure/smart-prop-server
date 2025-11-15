package utils

import (
	"strconv"
	"strings"

	"github.com/Brian-Mashavakure/smart-prop-server/pkg/database/models"
)

// return json response for api responses
func ReturnJsonResponse(status string, message string, data map[string]interface{}, errors map[string]interface{}) ApiResponse {
	response := ApiResponse{
		Status:  status,
		Message: message,
		Data:    data,
		Errors:  errors,
	}
	return response
}

// property ids seperator for ai response
func NumbersSeparator(input string) ([]uint, error) {
	lines := strings.Split(input, "\n")
	numbers := make([]uint, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		num, err := strconv.ParseUint(line, 10, 0)
		if err != nil {
			return nil, err
		}

		numbers = append(numbers, uint(num))
	}

	return numbers, nil
}

// properties filter
func FilterProperties(ids []uint, props []models.Property) []models.Property {
	lookup := make(map[uint]models.Property, len(props))
	for _, p := range props {
		lookup[p.ID] = p
	}

	result := make([]models.Property, 0, len(ids))
	for _, id := range ids {
		if p, ok := lookup[id]; ok {
			result = append(result, p)
		}
	}
	return result
}
