package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var setProviderCmd = &cobra.Command{
	Use:   "set-provider [name]",
	Short: "Set provider context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		endpoint, _ := cmd.Flags().GetString("endpoint")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")

		viper.Set("provider.name", name)
		viper.Set("provider.endpoint", endpoint)
		viper.Set("provider.user", user)
		viper.Set("provider.password", password)
		viper.WriteConfig()
		fmt.Printf("Provider context '%s' has been saved to ~/.config/vcfactl.json\n", name)
	},
}

var setOrgCmd = &cobra.Command{
	Use:   "set-org [name]",
	Short: "Set organization context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")

		viper.Set("org.name", name)
		viper.Set("org.user", user)
		viper.Set("org.password", password)
		viper.WriteConfig()
		fmt.Printf("Organization context '%s' has been saved to ~/.config/vcfactl.json\n", name)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setProviderCmd)
	configCmd.AddCommand(setOrgCmd)

	setProviderCmd.Flags().StringP("endpoint", "e", "", "Provider API Endpoint")
	setProviderCmd.Flags().StringP("user", "u", "", "Provider Username")
	setProviderCmd.Flags().StringP("password", "p", "", "Provider Password")

	setOrgCmd.Flags().StringP("user", "u", "", "Organization Username")
	setOrgCmd.Flags().StringP("password", "p", "", "Organization Password")
}
