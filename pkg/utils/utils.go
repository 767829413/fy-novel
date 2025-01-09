package utils

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"regexp"
	"time"
	"unicode"
)

// PrintAsJSON 接受任意类型的输入，将其转换为JSON字符串并打印
func PrintAsJSON(v interface{}) {
	var output string
	switch v := v.(type) {
	case string:
		output = v
	case int, int8, int16, int32, int64, float32, float64, bool:
		output = fmt.Sprintf("%v", v)
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			fmt.Println("Error converting to JSON:", err)
			return
		}
		output = string(jsonBytes)
	}
	fmt.Println(output)
}

// WriteAsJSON 接受任意类型的输入，将其转换为JSON字符串并写入指定文件
func WriteAsJSON(v interface{}, filename *string) error {
	var output string
	switch v := v.(type) {
	case string:
		output = v
	case int, int8, int16, int32, int64, float32, float64, bool:
		output = fmt.Sprintf("%v", v)
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("Error converting to JSON: %v", err)
		}
		output = string(jsonBytes)
	}

	// 写入文件

	realFileName := fmt.Sprintf("./tmp/%d.json", time.Now().UnixNano())
	if filename != nil {
		realFileName = *filename
	}
	err := os.WriteFile(realFileName, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v", err)
	}

	return nil
}

// CleanBlank removes all whitespace characters from the input string.
func CleanBlank(str string) string {
	result := []rune{}
	for _, c := range str {
		if !unicode.IsSpace(c) {
			result = append(result, c)
		}
	}
	return string(result)
}

// SpinWaitWithExponentialBackoff 是一个通用的自旋等待方法，带有指数退避机制。
// 它会在指定的最大等待时间内，反复检查给定的条件函数是否为真。
// 如果条件函数在超时之前返回真，则方法返回 true；如果超时，则返回 false。
// 每次条件检查失败后，等待的间隔时间会按指定的倍数增加。
//
// 参数:
// - condition: 一个函数类型的参数，用于检查某个条件是否满足。返回 true 表示条件满足。
// - interval: 初始的等待间隔时间。
// - maxWaitTime: 最大的等待时间。如果超过这个时间，方法将返回 false。
// - backoffMultiplier: 每次条件检查失败后，等待间隔时间的倍数。用于控制指数退避的增长速度。
func SpinWaitWithExponentialBackoff(
	condition func() bool,
	interval, maxWaitTime time.Duration,
	backoffMultiplier float64,
) bool {
	timeout := time.After(maxWaitTime)

	for {
		select {
		case <-timeout:
			return false
		default:
			if condition() {
				return true
			}
			time.Sleep(interval)
			interval = time.Duration(float64(interval) * backoffMultiplier) // 使用调用方提供的倍数
		}
	}
}

// 使用最大重试次数
// 重试间隔 1 2 4 8 。。。。
func SpinWaitMaxRetryAttempts(condition func() bool, maxRetryAttempts int) bool {
	interval, maxWaitTime, backoffMultiplier := calculateBackoffParameters(maxRetryAttempts)
	return SpinWaitWithExponentialBackoff(
		condition,
		interval,
		maxWaitTime,
		backoffMultiplier,
	)
}

func calculateBackoffParameters(maxRetries int) (time.Duration, time.Duration, float64) {
	initialInterval := 1 * time.Second
	backoffMultiplier := 2.0

	maxWaitTime := 1 * time.Second
	// 设置一个上限，例如1小时，以防止等待时间过长
	maxAllowedWaitTime := 1 * time.Hour
	for i := 0; i < maxRetries; i++ {
		maxWaitTime *= 2
	}
	if maxWaitTime > maxAllowedWaitTime {
		maxWaitTime = maxAllowedWaitTime
	}
	return initialInterval, maxWaitTime, backoffMultiplier
}

func StringToUniqueHash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

// GetGroup0 获得匹配的字符串，获得正则中分组0的内容
//
// regex: 匹配的正则
// content: 被匹配的内容
// 返回值: 匹配后得到的字符串，未匹配返回空字符串
func GetGroup0(regex string, content string) string {
	return Get(regex, content, 0)
}

// GetGroup1 获得匹配的字符串，获得正则中分组1的内容
//
// regex: 匹配的正则
// content: 被匹配的内容
// 返回值: 匹配后得到的字符串，未匹配返回空字符串
func GetGroup1(regex string, content string) string {
	return Get(regex, content, 1)
}

// Get 获得匹配的字符串
//
// regex: 匹配的正则
// content: 被匹配的内容
// groupIndex: 匹配正则的分组序号
// 返回值: 匹配后得到的字符串，未匹配返回空字符串
func Get(regex string, content string, groupIndex int) string {
	re, err := regexp.Compile(regex)
	if err != nil {
		return ""
	}

	matches := re.FindStringSubmatch(content)
	if len(matches) > groupIndex {
		return matches[groupIndex]
	}

	return ""
}
