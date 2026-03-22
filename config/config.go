package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error finding home directory:", err)
		os.Exit(1)
	}

	// ~/.config ディレクトリの作成
	configDir := filepath.Join(home, ".config")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
	}

	viper.AddConfigPath(configDir)
	viper.SetConfigName("vcfactl")
	viper.SetConfigType("json") // JSON形式に指定

	// ファイルが存在しない場合は空のJSONを作成
	configPath := filepath.Join(configDir, "vcfactl.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.WriteFile(configPath, []byte("{}"), 0644)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
	}
}
