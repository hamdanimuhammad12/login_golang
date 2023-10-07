package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Configurations exported
type Configurations struct {
	Server          ServerConfigurations
	Database        DatabaseConfigurations
	EXAMPLE_PATH    string
	EXAMPLE_VAR     string
	SSL_PRIVATE_KEY string
	SSL_PUBLIC_KEY  string

	ENV_TYPE           string
	DO_SPACES_KEY      string
	DO_SPACES_SECRET   string
	DO_SPACES_ENDPOINT string
	DO_SPACES_REGION   string
	DO_SPACES_BUCKET   string
	DO_SPACES_SSL      string

	LINK_ONESINYAL   string
	APP_ID_ONESINYAL string
	KEY_ONESINYAL    string

	SMTP_HOST   string
	SMTP_PORT   int
	EMAIL       string
	SENDER_NAME string
	PASSWORD    string

	LINKVERIFY    string
	URL_DASHBOARD string
	hostname      string

	LINK_EXTERNAL string
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Hostname string
	Port     int
	Ssl_Port int
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

func ConnectDB(config Configurations) *gorm.DB {

	URL := "host=" + config.Database.DBHost + " user=" + config.Database.DBUser + " password=" + config.Database.DBPassword + " dbname=" + config.Database.DBName + " port=" + config.Database.DBPort + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(URL), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	return db
}
