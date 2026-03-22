package cmd

import (
	"fmt"
	"os"

	"vcfactl/api"

	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post [path]",
	Short: "Execute POST API with a JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		file, _ := cmd.Flags().GetString("file")

		var payload []byte
		var err error
		if file != "" {
			payload, err = os.ReadFile(file)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				return
			}
		}

		// apiパッケージのExecuteAPIを呼び出す
		api.ExecuteAPI("POST", path, payload)
	},
}

func init() {
	// ルートコマンドにpostコマンドを登録
	rootCmd.AddCommand(postCmd)
	// -f, --file フラグの定義
	postCmd.Flags().StringP("file", "f", "", "JSON payload file path")
}
