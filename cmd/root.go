package cmd

import (
	"fmt"
	"os"

	"vcfactl/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "vcfactl",
	Short: "VCF Automation 9.x CLI Tool",
}

func Execute() {
	cobra.OnInitialize(config.InitConfig)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug output")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}
