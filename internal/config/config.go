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
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

// LoadConfig membaca konfigurasi dari file atau environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./configs") // Path ke file config
	viper.SetConfigName("config")    // Nama file (tanpa ekstensi)
	viper.SetConfigType("yaml")      // Tipe file

	viper.AutomaticEnv() // Baca juga dari environment variable jika ada

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
