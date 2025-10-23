package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}
type Alarms struct {
	CheckIntervalSeconds  int        `mapstructure:"check_interval_seconds"`
	RepeatIntervalMinutes int        `mapstructure:"repeat_interval_minutes"`
	Thresholds            Thresholds `mapstructure:"thresholds"`
}

// Thresholds mendefinisikan nilai ambang batas statis untuk alarm.
type Thresholds struct {
	SystemCPUCriticalPercent  float64 `mapstructure:"system_cpu_critical_percent"`
	SystemDiskCriticalPercent float64 `mapstructure:"system_disk_critical_percent"`
	SystemMemCriticalPercent  float64 `mapstructure:"system_mem_critical_percent"`
	WorkerHeartbeatTimeout    int     `mapstructure:"worker_heartbeat_timeout_seconds"`
}

type Settings struct {
	MonitoringIntervalSeconds int    `mapstructure:"monitoring_interval_seconds"`
	Alarms                    Alarms `mapstructure:"alarms"`
}

type Config struct {
	Database   DBConfig  `mapstructure:"database"`
	AppSetting AppConfig `mapstructure:"app"`
	Settings   Settings  `mapstructure:"settings"`
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
