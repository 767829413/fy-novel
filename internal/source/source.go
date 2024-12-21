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

	// 构建文件路径
	filePath := fmt.Sprintf("rule/rule%d.json", sourceId)
	// 读取文件内容
	ruleData, err := ruleFS.ReadFile(filePath)
	if err != nil {
		return nil
	}

	// 解析 JSON 到 Rule 结构体
	var rule model.Rule
	err = json.Unmarshal(ruleData, &rule)
	if err != nil {
		return nil
	}

	// 使用 LoadOrStore 方法来确保线程安全的写入
	actual, _ := ruleCache.LoadOrStore(sourceId, &rule)
	return actual.(*model.Rule)
}
