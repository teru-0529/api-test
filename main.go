package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

const RESULT_DIR = "./testResult/"

const PASS = "🟢 PASS"
const FAIL = "🔴 FAIL"
const SKIP = "🟡 SKIP"

// FUNCTION: テスト結果サマリ出力
func main() {
	now := jstNow()

	// PROCESS: 読込みファイル
	inFile, err := os.Open("plain.out")
	if err != nil {
		log.Printf("cannot open file: %v", err)
		return
	}
	defer inFile.Close()

	r := Result{}
	f := FailTrace{hasFail: false}
	tracePhase := true

	// PROCESS: ファイル分析
	s := bufio.NewScanner(inFile)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "=== RUN ") { // INFO: RUN
			f.createTrace(line)
		} else if strings.HasPrefix(line, "=== NAME") { // INFO: NAME
			f.createTrace(line)
		} else if strings.HasPrefix(line, "--- ") { // INFO: 総合結果(status/elapse)
			r.newResult(line)
			tracePhase = false
		} else if strings.HasPrefix(line, "    --- ") { // INFO: テスト結果(status/name/elapse)
			r.addDetail(line)
		} else if strings.HasSuffix(line, "skipped the test.") { // INFO: スキップはトレースしない
		} else if tracePhase { // INFO: トレースの時
			f.addTrace(line)
		}
	}
	if s.Err() != nil {
		// non-EOF error.
		log.Fatal(s.Err())
	}

	var suffix string
	if f.hasFail {
		suffix = "(FAILURE)"
	}
	outPath := path.Join(RESULT_DIR, fmt.Sprintf("Result-%s%s.md", now.Format("060102-150405"), suffix))

	// PROCESS: 書込みファイル
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Printf("cannot create file: %v", err)
		return
	}
	defer outFile.Close()

	outFile.WriteString(fmt.Sprintf("# %sAPI TestingReport\n\n", r.Icon))

	outFile.WriteString(fmt.Sprintf("- **📆operation datetime**: %s\n", now.Format("2006/01/02 15:04:05")))
	outFile.WriteString("- **📄summary**:\n\n")

	outFile.WriteString("  | STATUS | ELAPSED | PASS | FAIL | SKIP | TOTAL |\n")
	outFile.WriteString("  |---|--:|--:|--:|--:|--:|\n")
	outFile.WriteString(fmt.Sprintf("  | %s | %s | %d | %d | %d | %d |\n",
		r.Status, r.Elapse, r.PassedCount, r.FailedCount, r.SkippedCount, len(r.Details)))

	if f.hasFail {
		outFile.WriteString("\n- **📑failed tests**:")
		for _, detail := range f.Details {
			if !detail.hasFail {
				continue
			}
			outFile.WriteString(fmt.Sprintf("\n  - **%s**\n\n", detail.Name))
			outFile.WriteString("    ```dat\n")
			for _, rec := range detail.Values {
				outFile.WriteString(fmt.Sprintf("    %s\n", rec))
			}
			outFile.WriteString("    ```\n")
		}
	}

	outFile.WriteString("\n- **📁all results**:\n\n")
	outFile.WriteString("  <details>\n\n")
	outFile.WriteString("  <summary>Click for test details</summary>\n\n")
	outFile.WriteString("  | STATUS | ELAPSED | TEST NAME |\n")
	outFile.WriteString("  |--|--:|--|\n")
	for _, detail := range r.Details {
		outFile.WriteString(fmt.Sprintf("  | %s | %s | %s |\n", detail.Status, detail.Elapse, detail.Name))
	}
	outFile.WriteString("\n  </details>\n")

}

// FUNCTION: 日本現在時間取得
func jstNow() time.Time {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	return time.Now().In(jst)
}

// FUNCTION: ステータス装飾
func decorateStatus(status string) string {
	switch status {
	case "PASS":
		return PASS
	case "FAIL":
		return FAIL
	case "SKIP":
		return SKIP
	}
	return "N/A"
}

// STRUCT: テスト結果
type Result struct {
	Status       string
	Elapse       string
	Icon         string
	PassedCount  int
	FailedCount  int
	SkippedCount int
	Details      []ResultDetail
}

type ResultDetail struct {
	Status string
	Name   string
	Elapse string
}

// FUNCTION: 結果サマリ
func (r *Result) newResult(record string) {
	re := regexp.MustCompile(`^--- (.+?): (.+?) \((\d+\.\d+s)\)$`)
	matches := re.FindStringSubmatch(record)
	r.Status = decorateStatus(matches[1])
	r.Elapse = matches[3]
	if r.Status == FAIL {
		r.Icon = "📕"
	} else {
		r.Icon = "📘"
	}
}

// FUNCTION: 結果詳細
func (r *Result) addDetail(record string) {
	re := regexp.MustCompile(`^    --- (.+?): TestApi/(.+?) \((\d+\.\d+s)\)$`)
	matches := re.FindStringSubmatch(record)
	switch matches[1] {
	case "PASS":
		r.PassedCount++
	case "FAIL":
		r.FailedCount++
	case "SKIP":
		r.SkippedCount++
	}
	r.Details = append(r.Details, ResultDetail{
		Status: decorateStatus(matches[1]),
		Name:   matches[2],
		Elapse: matches[3],
	})
}

// STRUCT: エラートレース
type FailTrace struct {
	hasFail bool
	Details []FailDetail
}

type FailDetail struct {
	hasFail bool
	Name    string
	Values  []string
}

// FUNCTION: トレース作成
func (f *FailTrace) createTrace(record string) {
	re := regexp.MustCompile(`^=== .{4}  (.+?)$`)
	matches := re.FindStringSubmatch(record)
	f.Details = append(f.Details, FailDetail{Name: matches[1], Values: []string{}})
}

// FUNCTION: エラートレース追加
func (f *FailTrace) addTrace(record string) {
	f.hasFail = true
	detail := &f.Details[len(f.Details)-1]
	detail.hasFail = true
	detail.Values = append(detail.Values, record)
}
