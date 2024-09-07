package main

import (
	"log"
	"os"
	"path"
	"slices"
	"testing"

	"github.com/joho/godotenv"
	"github.com/teru-0529/api-test/dbaccess"
	"github.com/teru-0529/api-test/fixture"
)

const FIXTURE_DIR = "./testdata/fixture/"
const GOLDEN_DIR = "./testdata/golden/"
const UPDATE_FILES_LIST_PATH = "./testdata/updateFiles.yaml"

// TEST: API実行テスト
func TestApi(t *testing.T) {
	// PROCESS: configの呼び出し
	leadEnv()
	dbAccesser := dbaccess.New()
	updateFiles, err := fixture.UpdateFiles(UPDATE_FILES_LIST_PATH)
	if err != nil {
		t.Fatal(err)
	}

	// FIXME:
	log.Println(dbAccesser.PostgrestHost)
	log.Println(dbAccesser.ReseterHost)
	// FIXME:

	files, _ := os.ReadDir(FIXTURE_DIR)
	for _, file := range files {
		file := file
		update := slices.Contains(updateFiles, file.Name())
		// PROCESS: fixtureの生成
		fix, err := fixture.New(path.Join(FIXTURE_DIR, file.Name()))
		if err != nil {
			t.Errorf("fixture parse error:[%s], %v", file.Name(), err)
		}

		// FIXME:
		log.Println(fix.Name)
		log.Println(os.Getenv(fix.Execute.HostKey))
		log.Println(update)
		// FIXME:

		t.Run(fix.Name, func(t *testing.T) {
			// FIXME:
			if false {
				t.Error("error")
			}
			// FIXME:
		})

	}
}

// FUNCTION: UpdateFileリストの取得
func leadEnv() {
	// envファイルのロード
	_, err := os.Stat(".env")
	if !os.IsNotExist(err) {
		godotenv.Load()
		log.Print("loaded environment variables from .env file.")
	}
}

// FIXME:
func TestApiError(t *testing.T) {
	if true {
		t.Error("error2")
	}
}
