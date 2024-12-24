package source

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync"

	"fy-novel/internal/model"
)

//go:embed rule/*.json
var ruleFS embed.FS

var ruleCache sync.Map

func GetRuleBySourceID(sourceId int) *model.Rule {
	if v, ok := ruleCache.Load(sourceId); ok {
		return v.(*model.Rule)
	}

	// Build the file path
	filePath := fmt.Sprintf("rule/rule%d.json", sourceId)
	// Read the file content
	ruleData, err := ruleFS.ReadFile(filePath)
	if err != nil {
		return nil
	}

	// Parse JSON into Rule struct
	var rule model.Rule
	err = json.Unmarshal(ruleData, &rule)
	if err != nil {
		return nil
	}

	// Use LoadOrStore method to ensure thread-safe writing
	actual, _ := ruleCache.LoadOrStore(sourceId, &rule)
	return actual.(*model.Rule)
}
