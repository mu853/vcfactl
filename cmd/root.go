package cmd

import (
	"fmt"
	"os"

	"vcfactl/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vcfa",
	Short: "VCF Automation 9.x CLI Tool",
}

func Execute() {
	cobra.OnInitialize(config.InitConfig)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
