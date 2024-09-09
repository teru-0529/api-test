package fixture

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// TITLE: Fixture構造体
type Fixture struct {
	DataType    string `yaml:"dataType"`
	Version     string `yaml:"version"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	Reset        ResetOperation        `yaml:"reset"`
	Setup        []SetupTable          `yaml:"setupTable"`
	Execute      ExecuteOperation      `yaml:"execute"`
	Verification VerificationOperation `yaml:"verification"`
}

type ResetOperation struct {
	Sequences []ResetItem `yaml:"sequences"`
	Tables    []ResetItem `yaml:"tables"`
}

type SetupTable struct {
	Schema string `yaml:"schema"`
	Table  string `yaml:"table"`
	Body   string `yaml:"body"`
}

type ExecuteOperation struct {
	HostKey string     `yaml:"hostKey"`
	Method  string     `yaml:"method"`
	Path    string     `yaml:"path"`
	Headers []MapValue `yaml:"headers"`
	Body    string     `yaml:"body"`
}

type VerificationOperation struct {
	HttpStatus int           `yaml:"httpStatus"`
	Result     ResultVerify  `yaml:"execResult"`
	Tables     []TableVerify `yaml:"tables"`
}

type ResetItem struct {
	Schema string   `yaml:"schema"`
	Items  []string `yaml:"items"`
}

type MapValue struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type ResultVerify struct {
	IsCheck  bool     `yaml:"check"`
	Excludes []string `yaml:"exclude"`
}

type TableVerify struct {
	Schema   string   `yaml:"schema"`
	Table    string   `yaml:"table"`
	Excludes []string `yaml:"exclude"`
}

// TITLE: Settings構造体
type Settings struct {
	WhiteList    []string `yaml:"whiteList"`
	UpdateGorden []string `yaml:"updateGorden"`
	PartialTest  bool
}

// FUNCTION: Fixture構造のパース
func New(path string) (*Fixture, error) {
	// PROCESS: read
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	// PROCESS: unmarchal
	var fix Fixture
	err = yaml.Unmarshal([]byte(file), &fix)
	if err != nil {
		return nil, err
	}
	return &fix, nil
}

// FUNCTION: UpdateFileリストの取得
func NewSettings(path string) (*Settings, error) {
	// PROCESS: read
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	// PROCESS: unmarchal
	var files Settings
	err = yaml.Unmarshal([]byte(file), &files)
	if err != nil {
		return nil, err
	}
	files.PartialTest = len(files.WhiteList) > 0
	return &files, nil
}

func (f *Fixture) WriteSpecification(path string) error {
	// PROCESS: ファイル準備
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer file.Close()

	// PROCESS: 基本情報
	file.WriteString(fmt.Sprintf("# %s\n", f.Name))
	file.WriteString(fmt.Sprintf("\n%s\n", f.Description))

	// PROCESS: reset
	file.WriteString("\n## テスト対象DBの初期化\n")
	for _, seq := range f.Reset.Sequences {
		file.WriteString("\n```http\n# @name sdequenceリセット\n")
		file.WriteString(fmt.Sprintf("POST {{$dotenv DB_RESETER_API_HOST}}/schemas/%s/action-seq-reset HTTP/1.1\n", seq.Schema))
		file.WriteString("Content-Type: application/json\n\n")
		file.WriteString(fmt.Sprint(objToJson(seq.Items)))
		file.WriteString("\n```\n")
	}

	for _, table := range f.Reset.Tables {
		file.WriteString("\n```http\n# @name table初期化\n")
		file.WriteString(fmt.Sprintf("POST {{$dotenv DB_RESETER_API_HOST}}/schemas/%s/action-truncate HTTP/1.1\n", table.Schema))
		file.WriteString("Content-Type: application/json\n\n")
		file.WriteString(fmt.Sprint(objToJson(table.Items)))
		file.WriteString("\n```\n")
	}

	// PROCESS: setTable
	file.WriteString("\n## テストデータ登録\n")
	for _, table := range f.Setup {
		file.WriteString(fmt.Sprintf("\n```http\n# @name bulk insert(%s.%s)\n", table.Schema, table.Table))
		file.WriteString(fmt.Sprintf("POST {{$dotenv POSTGREST_API_HOST}}/%s HTTP/1.1\n", table.Table))
		file.WriteString("Content-Type: application/json\n")
		file.WriteString(fmt.Sprintf("Content-Profile: %s\n\n", table.Schema))
		file.WriteString(fmt.Sprint(jsonFormat(table.Body)))
		file.WriteString("\n```\n")
	}

	// PROCESS: execute
	file.WriteString("\n## テスト実行\n")
	file.WriteString("\n```http\n# @name execute API\n")
	file.WriteString(fmt.Sprintf("%s {{$dotenv %s}}%s HTTP/1.1\n", f.Execute.Method, f.Execute.HostKey, f.Execute.Path))
	file.WriteString("Content-Type: application/json\n")
	for _, head := range f.Execute.Headers {
		file.WriteString(fmt.Sprintf("%s: %s\n", head.Key, head.Value))
	}
	if f.Execute.Body != "" {
		file.WriteString("\n")
		file.WriteString(fmt.Sprint(jsonFormat(f.Execute.Body)))
		file.WriteString("\n")
	}
	file.WriteString("```\n")

	// PROCESS: verification
	file.WriteString("\n## 検証対象\n")

	file.WriteString(fmt.Sprintf("\n- API実行結果の http status が(%d)であること\n", f.Verification.HttpStatus))
	if f.Verification.Result.IsCheck {
		file.WriteString("- API実行結果が正しいこと\n\n")
	}
	for _, table := range f.Verification.Tables {
		file.WriteString(fmt.Sprintf("- テーブル (%s.%s) のデータが正しいこと\n", table.Schema, table.Table))

		file.WriteString(fmt.Sprintf("\n```http\n# @name get all data(%s.%s)\n", table.Schema, table.Table))
		file.WriteString(fmt.Sprintf("GET {{$dotenv POSTGREST_API_HOST}}/%s HTTP/1.1\n", table.Table))
		file.WriteString(fmt.Sprintf("Accept-Profile: %s\n", table.Schema))
		file.WriteString("```\n\n")

	}

	// PROCESS: フッター
	file.WriteString("-----\n")
	file.WriteString(fmt.Sprintf("api-test-specification.(v%s)\n", f.Version))
	return nil
}

func jsonFormat(val string) string {
	var obj interface{}
	json.Unmarshal([]byte(val), &obj)
	return objToJson(obj)
}

func objToJson(obj interface{}) string {
	result, _ := json.MarshalIndent(&obj, "", "  ")
	return string(result)
}
