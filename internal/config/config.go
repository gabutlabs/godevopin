package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}
type Config struct {
	Database   DBConfig  `mapstructure:"database"`
	AppSetting AppConfig `mapstructure:"app"`
}

type DBConfig struct {
	Directory   string `mapstructure:"directory"`
	MainPath    string `mapstructure:"main_path"`
	MetricsPath string `mapstructure:"metrics_path"`
	LogsPath    string `mapstructure:"logs_path"`
}

// LoadConfig membaca konfigurasi dari file atau environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./configs")             // Path lokal untuk development
	viper.AddConfigPath(".")                     // Direktori saat ini
	viper.AddConfigPath("/opt/devopin")          // Path instalasi global (Linux/macOS)
	viper.AddConfigPath("/etc/devopin")          // Path konfigurasi global
	viper.AddConfigPath("$HOME/.config/devopin") // Path konfigurasi user

	viper.SetConfigName("config") // Nama file (tanpa ekstensi)
	viper.SetConfigType("yaml")   // Tipe file

	viper.AutomaticEnv() // Baca juga dari environment variable jika ada

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
