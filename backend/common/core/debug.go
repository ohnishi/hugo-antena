package core

import (
	"encoding/json"
	"fmt"
)

// PrintJSON デバッグ用にJSONにして表示する
func PrintJSON(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(bytes))
}
