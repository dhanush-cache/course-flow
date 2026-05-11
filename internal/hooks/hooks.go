package hooks

import (
	"fmt"
)

type HookFunc func(folderPath string) error

var Registry = map[string]HookFunc{}

func init() {
	Registry["sql"] = SQL
	Registry["node"] = Node
}

func Run(courseID string, folderPath string) error {
	hook, exists := Registry[courseID]
	if !exists {
		return nil
	}

	fmt.Printf("Running hook for course ID: %s\n", courseID)
	return hook(folderPath)
}
