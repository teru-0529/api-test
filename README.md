# api-test

APIのリグレッションテストをゴールデンテストで実施する

## appendix

* [Golangで特定ファイルのテストのみ実行する](https://masarufuruya.hatenadiary.jp/entry/2017/08/26/184550)
* [【Windows】バッチファイルで現在の日付・時刻を取得するdate・timeコマンドの使い方](https://qiita.com/setonao/items/192c5c63074ea72a229b)

### アウトプットのプレーンテキストをMDにコンバートするサンプルコード

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// ファイルを開く
	file, err := os.Open("test_results.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var totalTests, passedTests, failedTests int
	var testDetails []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "=== RUN") {
			totalTests++
		} else if strings.HasPrefix(line, "--- PASS") {
			passedTests++
		} else if strings.HasPrefix(line, "--- FAIL") {
			failedTests++
			testDetails = append(testDetails, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// マークダウン形式で出力
	fmt.Println("# Test Results")
	fmt.Printf("- **Total Tests**: %d\n", totalTests)
	fmt.Printf("- **Passed**: %d\n", passedTests)
	fmt.Printf("- **Failed**: %d\n", failedTests)
	fmt.Println()
	fmt.Println("## Test Details")

	for _, detail := range testDetails {
		fmt.Println("### Failed Test")
		fmt.Printf("- **Error Message**: %s\n", detail)
	}
}
```
