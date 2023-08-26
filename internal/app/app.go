package app

import (
	"beastpark/meetinginvitationservice/pkg/config"
	"beastpark/meetinginvitationservice/pkg/database"
	"beastpark/meetinginvitationservice/pkg/log"
)

type App struct {
	conf *config.Config
	log  *log.Logger
}

func NewApp(conf *config.Config) *App {

	logInst := log.GetInstance()
	logInst.Init(&conf.Log)

	logInst.Infof("config: %+v", *conf)

	err := database.GetInstance().Init(&conf.DB)
	if err != nil {
		logInst.Errorf("db init fail, %s", err.Error())
		return nil
	}

	application := &App{
		conf: conf,
		log:  logInst,
	}

	return application
}

func (a *App) Run() error {

	httpSrv := NewServer(a)

	err := httpSrv.Run()

	return err
}
