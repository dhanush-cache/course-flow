package dreamsofcode

import (
	"encoding/json"
	"fmt"
	"os"
)

type courseWrapper struct {
	Course Course `json:"course"`
}

// LoadCourse loads the JSON file and returns a pointer to Course.
func LoadCourse(path string) (*Course, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open course file: %w", err)
	}
	defer file.Close()

	var wrapper courseWrapper
	if err := json.NewDecoder(file).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode course json: %w", err)
	}

	return &wrapper.Course, nil
}
