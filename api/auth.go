package api

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// 引数 isOrg に応じて、プロバイダ用 / 組織用 の認証を自動で切り替える
func GetBearerToken(isOrg bool) (string, error) {
	endpoint := viper.GetString("provider.endpoint")
	var user, password, orgName, loginURL string

	if isOrg {
		// 組織APIの場合は、組織ユーザーとしてログイン
		user = viper.GetString("org.user")
		password = viper.GetString("org.password")
		orgName = viper.GetString("org.name") // 例: "org-01"
		loginURL = fmt.Sprintf("%s/cloudapi/1.0.0/sessions", endpoint)

		if user == "" || password == "" || orgName == "" {
			return "", fmt.Errorf("organization credentials are not set completely")
		}
	} else {
		// プロバイダAPIの場合は、プロバイダ管理者としてログイン
		user = viper.GetString("provider.user")
		password = viper.GetString("provider.password")
		orgName = "system"
		loginURL = fmt.Sprintf("%s/cloudapi/1.0.0/sessions/provider", endpoint)

		if user == "" || password == "" {
			return "", fmt.Errorf("provider credentials are not set completely")
		}
	}

	authString := fmt.Sprintf("%s@%s:%s", user, orgName, password)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))

	req, err := http.NewRequest("POST", loginURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Basic "+encodedAuth)
	req.Header.Set("Accept", "application/json;version=40.0")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to login API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("login failed with status: %d (user: %s@%s)", resp.StatusCode, user, orgName)
	}

	token := resp.Header.Get("x-vmware-vcloud-access-token")
	if token == "" {
		token = resp.Header.Get("X-Vmware-Vcloud-Access-Token")
	}

	if token == "" {
		return "", fmt.Errorf("bearer token not found in response")
	}

	return token, nil
}
