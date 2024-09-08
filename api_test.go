package main

import (
	"log"
	"os"
	"path"
	"slices"
	"testing"

	"github.com/joho/godotenv"
	"github.com/teru-0529/api-test/api"
	"github.com/teru-0529/api-test/fixture"
)

const FIXTURE_DIR = "./testdata/fixture/"
const GOLDEN_DIR = "./testdata/golden/"
const UPDATE_FILES_LIST_PATH = "./testdata/updateFiles.yaml"

// TEST: API実行テスト
func TestApi(t *testing.T) {
	// PROCESS: configの呼び出し
	leadEnv()
	apiAccesser := api.New()
	updateFiles, err := fixture.UpdateFiles(UPDATE_FILES_LIST_PATH)
	if err != nil {
		t.Fatal(err)
	}

	files, _ := os.ReadDir(FIXTURE_DIR)
	for _, file := range files {
		file := file
		update := slices.Contains(updateFiles, file.Name())

		// PROCESS: fixtureの生成
		fix, err := fixture.New(path.Join(FIXTURE_DIR, file.Name()))
		if err != nil {
			t.Errorf("fixture parse error:[%s], (%v).", file.Name(), err)
			log.Printf("fixture parse error:[%s], (%v).", file.Name(), err)
			log.Println("skip test")
			continue
		}

		// FIXME:
		log.Println(os.Getenv(fix.Execute.HostKey))
		// FIXME:

		t.Run(fix.Name, func(t *testing.T) {
			log.Println(fix.Name)
			if update {
				log.Println(" - (*) golden file update.")
			}

			// PROCESS: Dbのリセット(対象テーブルのtruncate/sequenceの初期化)
			if err = apiAccesser.Reset(fix.Reset); err != nil {
				log.Println(" - reset NG")
				t.Fatalf("reset failured: (%v).", err)
			}
			log.Println(" - reset OK")

			// PROCESS: テストデータのInsert
			for _, item := range fix.Setup {
				if err = apiAccesser.BulkInsert(item.Schema, item.Table, item.Body); err != nil {
					log.Println(" - setupTable NG")
					log.Printf("   - %v", err)
					t.Fatalf("setup failured: (%v).", err)
				}
			}
			log.Println(" - setupTable OK")

			// PROCESS: API実行

			// FIXME:
			if false {
				t.Error("error")
			}
			// FIXME:
		})

	}
}

// FUNCTION: 環境変数への設定(.envファイルがある場合のみ)
func leadEnv() {
	// envファイルのロード
	_, err := os.Stat(".env")
	if !os.IsNotExist(err) {
		godotenv.Load()
		log.Print("loaded environment variables from .env file.")
	}
}
