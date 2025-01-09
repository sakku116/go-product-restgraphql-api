package config

import (
	"github.com/spf13/viper"
)

type EnvsSchema struct {
	HOST                   string
	PORT                   int
	LOG_LEVEL              string
	JWT_SECRET_KEY         string
	JWT_EXP_MINS           int
	REFRESH_TOKEN_EXP_MINS int

	INITIAL_USER_USERNAME  string
	INITIAL_USER_PASSWORD  string
	INITIAL_ADMIN_USERNAME string
	INITIAL_ADMIN_PASSWORD string

	MONGO_URL string
	MONGO_DB  string
}

var Envs *EnvsSchema

func envInitiator() {
	Envs = &EnvsSchema{
		HOST:                   viper.GetString("HOST"),
		PORT:                   viper.GetInt("PORT"),
		LOG_LEVEL:              viper.GetString("LOG_LEVEL"),
		JWT_SECRET_KEY:         viper.GetString("JWT_SECRET_KEY"),
		JWT_EXP_MINS:           viper.GetInt("JWT_EXP_MINS"),
		REFRESH_TOKEN_EXP_MINS: viper.GetInt("REFRESH_TOKEN_EXP_MINS"),

		INITIAL_USER_USERNAME:  viper.GetString("INITIAL_USER_USERNAME"),
		INITIAL_USER_PASSWORD:  viper.GetString("INITIAL_USER_PASSWORD"),
		INITIAL_ADMIN_USERNAME: viper.GetString("INITIAL_ADMIN_USERNAME"),
		INITIAL_ADMIN_PASSWORD: viper.GetString("INITIAL_ADMIN_PASSWORD"),

		MONGO_URL: viper.GetString("MONGO_URL"),
		MONGO_DB:  viper.GetString("MONGO_DB"),
	}
}

func InitEnv(filepath string) {
	viper.SetConfigType("env")
	viper.SetConfigFile(filepath)
	if err := viper.ReadInConfig(); err != nil {
		logger.Warningf("error loading environment variables from %s: %w", filepath, err)
	}
	viper.AutomaticEnv()
	envInitiator()
}
