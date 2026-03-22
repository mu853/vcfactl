package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// APIのフィルタ表示用ユーティリティ
func filterAPIs(apis []string, args []string) {
	search := ""
	if len(args) > 0 {
		search = strings.ToLower(args[0])
	}
	fmt.Println("Matched APIs:")
	for _, api := range apis {
		if search == "" || strings.Contains(strings.ToLower(api), search) {
			fmt.Printf("  - %s\n", api)
		}
	}
}

var apiProviderCmd = &cobra.Command{
	Use:   "api-provider [search_string]",
	Short: "List and filter Provider APIs",
	Run: func(cmd *cobra.Command, args []string) {
		// プロバイダ管理APIのリスト (VCF Automation 9.x想定)
		apis := []string{
			"GET /iaas/api/orgs",
			"POST /iaas/api/orgs",
			"GET /iaas/api/vdcs",
			"GET /cloudapi/1.0.0/apps",
			"POST /cloudapi/1.0.0/apps",
			"DELETE /cloudapi/1.0.0/apps/{id}",
		}
		filterAPIs(apis, args)
	},
}

var apiOrgCmd = &cobra.Command{
	Use:   "api-org [search_string]",
	Short: "List and filter Organization APIs",
	Run: func(cmd *cobra.Command, args []string) {
		// 組織管理APIのリスト
		apis := []string{
			"GET /iaas/api/projects",
			"POST /iaas/api/projects",
			"GET /iaas/api/machines",
			"GET /plan",
			"POST /plan",
			"DELETE /plan/{id}",
		}
		filterAPIs(apis, args)
	},
}

func init() {
	// ルートコマンドに検索コマンドを登録
	rootCmd.AddCommand(apiProviderCmd)
	rootCmd.AddCommand(apiOrgCmd)
}
