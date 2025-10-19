package config

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type IAppConfig interface {
	GetCloudEnv() string
	IsDebugMode() bool
	GetBaseURL() string
	GetDBConnectionString() string
	GetAuthHashPepper() string
	GetJWTSecretKey() string
	GetMJAPIKeyPublic() string
	GetMJAPIKeyPrivate() string
	GetCorsAllowedOrigin() string
	GetCookieDomain() string
	GetSystemAPIKey() string
}

type AppConfig struct {
	cloudEnv           string
	debugMode          bool
	baseURL            string
	dBConnectionString string
	authHashPepper     string
	jwtSecretKey       string
	mjAPIKeyPublic     string
	mjAPIKeyPrivate    string
	corsAllowedOrigin  string
	cookieDomain       string
	systemAPIKey       string
}

func NewAppConfig() IAppConfig {

	appConfig := &AppConfig{}
	// Required environment variables
	appConfig.cloudEnv = os.Getenv("CLOUD_ENV")

	// Default values
	appConfig.debugMode = os.Getenv("DEBUG_MODE") == "true"
	appConfig.baseURL = "http://localhost:3000"
	appConfig.dBConnectionString = ""
	appConfig.authHashPepper = ""
	appConfig.jwtSecretKey = ""
	appConfig.mjAPIKeyPublic = ""
	appConfig.mjAPIKeyPrivate = ""
	appConfig.corsAllowedOrigin = "http://localhost:3000"
	appConfig.cookieDomain = "localhost"
	appConfig.systemAPIKey = ""

	if appConfig.cloudEnv == "" {
		log.Fatal("[CLOUD_ENV] is required")
	}

	// Load any additional variables from the environment and override the secret manager values
	appConfig.baseURL = os.Getenv("BASE_URL")
	appConfig.dBConnectionString = os.Getenv("DB_CONNECTION_STRING")
	appConfig.authHashPepper = os.Getenv("AUTH_HASH_PEPPER")
	appConfig.jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	appConfig.mjAPIKeyPublic = os.Getenv("MJ_APIKEY_PUBLIC")
	appConfig.mjAPIKeyPrivate = os.Getenv("MJ_APIKEY_PRIVATE")
	appConfig.systemAPIKey = os.Getenv("SYSTEM_API_KEY")

	// Load optional configuration with defaults
	if corsOrigin := os.Getenv("CORS_ALLOWED_ORIGIN"); corsOrigin != "" {
		appConfig.corsAllowedOrigin = corsOrigin
	}
	if cookieDomain := os.Getenv("COOKIE_DOMAIN"); cookieDomain != "" {
		appConfig.cookieDomain = cookieDomain
	}

	errorList := ""

	if appConfig.dBConnectionString == "" {
		errorList += "[DB_CONNECTION_STRING]\n"
	}

	if appConfig.authHashPepper == "" {
		errorList += "[AUTH_HASH_PEPPER]\n"
	}

	if appConfig.jwtSecretKey == "" {
		errorList += "[JWT_SECRET_KEY]\n"
	}

	if appConfig.mjAPIKeyPublic == "" {
		errorList += "[MJ_APIKEY_PUBLIC]\n"
	}

	if appConfig.mjAPIKeyPrivate == "" {
		errorList += "[MJ_APIKEY_PRIVATE]\n"
	}

	if appConfig.systemAPIKey == "" {
		errorList += "[SYSTEM_API_KEY]\n"
	}

	if errorList != "" {
		errorList = "Missing environment variables:\n" + errorList
		panic(errorList)
	}

	return appConfig
}

func (a *AppConfig) GetCloudEnv() string {
	return a.cloudEnv
}

func (a *AppConfig) IsDebugMode() bool {
	return a.debugMode
}

func (a *AppConfig) GetBaseURL() string {
	return a.baseURL
}

func (a *AppConfig) GetDBConnectionString() string {
	return a.dBConnectionString
}

func (a *AppConfig) GetAuthHashPepper() string {
	return a.authHashPepper
}

func (a *AppConfig) GetJWTSecretKey() string {
	return a.jwtSecretKey
}

func (a *AppConfig) GetMJAPIKeyPublic() string {
	return a.mjAPIKeyPublic
}

func (a *AppConfig) GetMJAPIKeyPrivate() string {
	return a.mjAPIKeyPrivate
}

func (a *AppConfig) GetCorsAllowedOrigin() string {
	return a.corsAllowedOrigin
}

func (a *AppConfig) GetCookieDomain() string {
	return a.cookieDomain
}

func (a *AppConfig) GetSystemAPIKey() string {
	return a.systemAPIKey
}
