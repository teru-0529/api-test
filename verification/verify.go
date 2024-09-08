package verification

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Cside/jsondiff"
	"github.com/mattn/go-jsonpointer"
)

// FUNCTION: Json検証
func JsonVerify(t *testing.T, jsonData []byte, goldenPath string, update bool, excludes []string) {
	t.Helper()

	// PROCESS: 検証対象外項目の除外
	var obj interface{}
	json.Unmarshal(jsonData, &obj)
	for _, exclude := range excludes {
		obj = removeExcludeJson(t, obj, exclude)
	}
	jsonData, _ = json.MarshalIndent(&obj, "", "  ")

	// PROCESS: goldenファイルの更新
	if update {
		os.WriteFile(goldenPath, jsonData, 0644)
	}

	// PROCESS: goldenファイルの参照
	expected, _ := os.ReadFile(goldenPath)

	// PROCESS: 検証
	if diff := jsondiff.Diff(jsonData, expected); diff != "" {
		t.Errorf("response jsons are not correct. diff:\n%s", diff)
	}
	// log.Println(excludes)
}

// // FUNCTION: JSON除外項目の対応(除外項目配列)
// func removeExcludeJson(t *testing.T, org string, excludes []string) string {
// 	t.Helper()

// 	var obj interface{}
// 	json.Unmarshal([]byte(org), &obj)

// 	// PROCESS: 項目ごとの対応
// 	for _, exclude := range excludes {
// 		obj = removeExclude(t, obj, exclude)
// 	}

// 	result, _ := json.MarshalIndent(&obj, "", "  ")
// 	return string(result)
// }

// FUNCTION: JSON除外項目の対応
func removeExcludeJson(t *testing.T, obj interface{}, exclude string) interface{} {
	t.Helper()

	// PROCESS: ワイルドカード文字列($)の検索
	if strings.Contains(exclude, "$") {
		// PROCESS: 再帰処理による除外

		wcNum := strings.Index(exclude, "$")
		result, err := jsonpointer.Get(obj, exclude[:wcNum])
		if err != nil { // PROCESS: 親構造が存在しなければ対応せずそのまま戻す
			fmt.Println(err)

		} else {
			switch result := result.(type) {
			case []interface{}: // PROCESS: ワイルドカード構造が配列の場合：各配列要素に対してプレフィックスの除外処理を再帰的に実行し結果で入れ替える
				for i := range len(result) {
					sub, _ := jsonpointer.Get(obj, fmt.Sprintf("%s%d", exclude[:wcNum], i))
					sub = removeExcludeJson(t, sub, exclude[wcNum+1:])

					jsonpointer.Set(obj, fmt.Sprintf("%s%d", exclude[:wcNum], i), sub)
				}

			default: // PROCESS: ワイルドカード構造が配列以外の場合、対応せずそのまま返す
				fmt.Println("not useable wildcard excluded array type.")
			}

		}
		return obj

	} else {
		// PROCESS: 通常の除外処理
		result, err := jsonpointer.Remove(obj, exclude)
		if err != nil {
			return obj
		} else {
			return result
		}
	}
}
