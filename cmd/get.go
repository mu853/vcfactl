package cmd

import (
	"strings"

	"vcfactl/api"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [path_or_resource]",
	Short: "Execute GET API or fetch specific resource",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]
		if strings.HasPrefix(arg, "/") {
			api.ExecuteAPI("GET", arg, nil)
		} else {
			api.ExecuteResourceGet(arg)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
