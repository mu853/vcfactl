package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func isOrgAPI(path string) bool {
	if strings.HasPrefix(path, "/cloudapi/") {
		return false
	}
	providerPaths := []string{
		"/iaas/api/orgs",
		"/iaas/api/vdcs",
	}
	for _, p := range providerPaths {
		if strings.HasPrefix(path, p) {
			return false
		}
	}
	orgPaths := []string{
		"/plan",
		"/relocation",
		"/iaas/api/projects",
		"/iaas/api/machines",
	}
	for _, p := range orgPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// ▼ 追加: HTTPリクエストを行う共通関数（リファクタリング）
func doHTTPRequest(method string, path string, payload []byte) ([]byte, error) {
	endpoint := viper.GetString("provider.endpoint")
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint not set. Run 'vcfactl config set-provider' first")
	}

	url := endpoint + path
	isOrg := isOrgAPI(path)

	token, err := GetBearerToken(isOrg)
	if err != nil {
		return nil, fmt.Errorf("authentication error: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("request creation error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	if strings.HasPrefix(path, "/iaas/api/") || strings.HasPrefix(path, "/relocation/") {
		req.Header.Set("Accept", "application/json")
	} else {
		req.Header.Set("Accept", "application/json;version=40.0")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API execution error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %s\nError Details: %s", resp.Status, string(body))
	}

	return body, nil
}

func printDebug(format string, a ...interface{}) {
	if viper.GetBool("debug") {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}

// 汎用API実行 (生のJSONを出力)
func ExecuteAPI(method string, path string, payload []byte) {
	printDebug("Executing %s %s...\n", method, viper.GetString("provider.endpoint")+path)

	body, err := doHTTPRequest(method, path, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	PrintPrettyJSON(body)
}

// ▼ 修正: リソース指定のGET実行 (本物のAPIから取得してテーブル表示)
func ExecuteResourceGet(resource string) {
	var path string

	// リソース名と実際のAPIパスのマッピング
	switch resource {
	case "org":
		path = "/iaas/api/orgs"
	case "vdc":
		path = "/iaas/api/vdcs"
	case "project":
		path = "/iaas/api/projects"
	case "vm":
		path = "/iaas/api/machines"
	case "app":
		path = "/cloudapi/1.0.0/apps"
	default:
		fmt.Printf("Unknown resource type: %s\n", resource)
		return
	}

	printDebug("Fetching resource: %s (API: GET %s)\n", resource, path)
	body, err := doHTTPRequest("GET", path, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// JSONのパース
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Failed to parse response JSON")
		return
	}

	// IaaS系APIは "content"、VCD系APIは "values" という配列にデータが入っている
	var items []interface{}
	if val, ok := result["content"]; ok {
		items = val.([]interface{})
	} else if val, ok := result["values"]; ok {
		items = val.([]interface{})
	} else {
		fmt.Println("No items found in response.")
		return
	}

	// テーブルのヘッダ (IDにはURNが入ることもあるため幅を広めに確保)
	fmt.Printf("%-40s | %-30s\n", "ID", "Name")
	fmt.Println(strings.Repeat("-", 75))

	// 各要素から id と name を抽出して表示
	for _, item := range items {
		if m, ok := item.(map[string]interface{}); ok {
			id := ""
			name := ""
			if m["id"] != nil {
				id = fmt.Sprintf("%v", m["id"])
			}
			if m["name"] != nil {
				name = fmt.Sprintf("%v", m["name"])
			}
			fmt.Printf("%-40s | %-30s\n", id, name)
		}
	}
}

// JSONの整形出力
func PrintPrettyJSON(data []byte) {
	if len(data) == 0 {
		printDebug("(No Content)\n")
		return
	}
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		fmt.Println(string(data))
		return
	}
	fmt.Println(prettyJSON.String())
}
