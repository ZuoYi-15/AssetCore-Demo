package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig
	MySQL MySQLConfig
	Redis RedisConfig
	Kafka KafkaConfig
	JWT   JWTConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port int
}

type MySQLConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers []string
	Enabled bool
}

type JWTConfig struct {
	Secret string
	Issuer string
}

func Load() *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	v.SetEnvPrefix("ASSET_CORE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)
	_ = v.ReadInConfig()

	return &Config{
		App: AppConfig{
			Name: v.GetString("app.name"),
			Env:  v.GetString("app.env"),
			Port: v.GetInt("app.port"),
		},
		MySQL: MySQLConfig{
			DSN:          v.GetString("mysql.dsn"),
			MaxOpenConns: v.GetInt("mysql.max_open_conns"),
			MaxIdleConns: v.GetInt("mysql.max_idle_conns"),
		},
		Redis: RedisConfig{
			Addr:     v.GetString("redis.addr"),
			Password: v.GetString("redis.password"),
			DB:       v.GetInt("redis.db"),
		},
		Kafka: KafkaConfig{
			Brokers: v.GetStringSlice("kafka.brokers"),
			Enabled: v.GetBool("kafka.enabled"),
		},
		JWT: JWTConfig{
			Secret: v.GetString("jwt.secret"),
			Issuer: v.GetString("jwt.issuer"),
		},
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "asset-core")
	v.SetDefault("app.env", "dev")
	v.SetDefault("app.port", 8080)
	v.SetDefault("mysql.dsn", "asset_core:asset_core_password@tcp(127.0.0.1:3306)/asset_core?charset=utf8mb4&parseTime=True&loc=Local")
	v.SetDefault("mysql.max_open_conns", 50)
	v.SetDefault("mysql.max_idle_conns", 10)
	v.SetDefault("redis.addr", "127.0.0.1:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("kafka.brokers", []string{"127.0.0.1:9092"})
	v.SetDefault("kafka.enabled", false)
	v.SetDefault("jwt.secret", "change-me")
	v.SetDefault("jwt.issuer", "asset-core")
}
