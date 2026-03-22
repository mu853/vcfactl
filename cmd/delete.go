package cmd

import (
	"vcfactl/api"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [path]",
	Short: "Execute DELETE API",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		// apiパッケージのExecuteAPIを呼び出す
		api.ExecuteAPI("DELETE", path, nil)
	},
}

func init() {
	// ルートコマンドにdeleteコマンドを登録
	rootCmd.AddCommand(deleteCmd)
}
