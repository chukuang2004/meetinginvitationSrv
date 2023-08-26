package config

import (
	"beastpark/meetinginvitationservice/pkg/database"
	"beastpark/meetinginvitationservice/pkg/http"
	"beastpark/meetinginvitationservice/pkg/log"
)

type WXConfig struct {
	AppId  string
	Secret string
	State  string
}

type Config struct {
	WX   WXConfig
	DB   database.Config
	Log  log.Config
	Http http.Config
}

func NewConfig() *Config {
	return &Config{}
}
