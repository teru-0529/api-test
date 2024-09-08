package fixture

import (
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
	HttpStatus string        `yaml:"httpStatus"`
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

// FUNCTION: Fixture構造のパース
func New(path string) (*Fixture, error) {
	// PROCESS: read
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %s", err.Error())
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
func UpdateFiles(path string) ([]string, error) {
	// PROCESS: read
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %s", err.Error())
	}

	// PROCESS: unmarchal
	var files []string
	err = yaml.Unmarshal([]byte(file), &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}
