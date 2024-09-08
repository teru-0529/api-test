package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/sethvargo/go-envconfig"
	"github.com/teru-0529/api-test/fixture"
)

// STRUCT:
type ApiAccesser struct {
	client        *http.Client
	PostgrestHost string `env:"POSTGREST_API_HOST"`
	ReseterHost   string `env:"DB_RESETER_API_HOST"`
}

// FUNCTION: APIアクセス用構造体の生成
func New() *ApiAccesser {
	var accessor ApiAccesser
	envconfig.Process(context.Background(), &accessor)

	accessor.client = &http.Client{}
	return &accessor
}

// FUNCTION: API実行
func (api *ApiAccesser) Execute(exec fixture.ExecuteOperation) ([]byte, int, error) {
	hostName := os.Getenv(exec.HostKey)
	urlPath, _ := url.JoinPath(hostName, exec.Path)

	req, _ := http.NewRequest(exec.Method, urlPath, strings.NewReader(exec.Body))
	req.Header.Set("Content-Type", "application/json")
	for _, head := range exec.Headers {
		req.Header.Set(head.Key, head.Value)
	}

	// API実行
	res, err := api.client.Do(req)
	if err != nil {
		return nil, 500, fmt.Errorf("execute failured: (%w).", err)
	}
	defer res.Body.Close()
	jsonData, err := io.ReadAll(res.Body)
	return jsonData, res.StatusCode, err
}

// FUNCTION: リセット
func (api *ApiAccesser) Reset(resetItems fixture.ResetOperation) error {
	// PROCESS: シーケンス
	for _, resetSeq := range resetItems.Sequences {
		data, _ := json.Marshal(resetSeq.Items)
		if err := api.squenceReset(resetSeq.Schema, data); err != nil {
			return err
		}
	}
	// PROCESS: トランケート
	for _, truncate := range resetItems.Tables {
		data, _ := json.Marshal(truncate.Items)
		if err := api.truncate(truncate.Schema, data); err != nil {
			return err
		}
	}
	// FIXME:
	return nil
}

// FUNCTION: sequence
func (api *ApiAccesser) squenceReset(schema string, sequences []byte) error {
	urlPath, _ := url.JoinPath(api.ReseterHost, "schemas", schema, "action-seq-reset")

	req, _ := http.NewRequest("POST", urlPath, bytes.NewReader(sequences))
	req.Header.Set("Content-Type", "application/json")

	_, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("sequence reset failured: (%w).", err)
	}
	return err
}

// FUNCTION: truncate
func (api *ApiAccesser) truncate(schema string, tables []byte) error {
	urlPath, _ := url.JoinPath(api.ReseterHost, "schemas", schema, "action-truncate")

	req, _ := http.NewRequest("POST", urlPath, bytes.NewReader(tables))
	req.Header.Set("Content-Type", "application/json")

	_, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("truncate failured: (%w).", err)
	}
	return err
}

// FUNCTION: データの登録
// postgRESTによるDBアクセス、JSONで指定したデータを投入(配列で指定可能)。
func (api *ApiAccesser) BulkInsert(schema string, table string, jsonData string) error {
	urlPath, _ := url.JoinPath(api.PostgrestHost, table)

	req, _ := http.NewRequest("POST", urlPath, strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Profile", schema)

	// API実行
	res, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("bulk insert failured: (%w).", err)
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: [%s].", errors.New("Bad request"), jsonData)
	case http.StatusNotFound:
		return fmt.Errorf("%w: [%s.%s].", errors.New("Table not found"), schema, table)
	}
	return nil
}

// FUNCTION: データの読込み
// postgRESTによるDBアクセス、全件検索をしJSONを返す。
func (api *ApiAccesser) GetAll(schema string, table string) ([]byte, error) {
	urlPath, _ := url.JoinPath(api.PostgrestHost, table)

	req, _ := http.NewRequest("GET", urlPath, nil)
	req.Header.Set("Accept-profile", schema)

	// API実行
	res, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("all get failured: (%w).", err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: [%s.%s].", errors.New("Table not found"), schema, table)
	}
	return io.ReadAll(res.Body)
}
