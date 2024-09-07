package main

import (
	"fmt"
	"os"
	"testing"
)

const FIXTURE_DIR = "./testdata/fixture/"
const GOLDEN_DIR = "./testdata/golden/"

// TEST: API実行テスト
func TestApi(t *testing.T) {
	files, _ := os.ReadDir(FIXTURE_DIR)
	for _, file := range files {
		file := file
		// PROCESS: fixtureの生成
		fmt.Println(file.Name())

		t.Run(file.Name(), func(t *testing.T) {
			if false {
				t.Error("error")
			}
		})

	}
}

func TestApiError(t *testing.T) {
	if true {
		t.Error("error2")
	}
}
