package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"beastpark/meetinginvitationservice/internal/app"
	"beastpark/meetinginvitationservice/pkg/config"
)

func main() {

	conf := config.NewConfig()

	app := &cli.App{
		Name: "MeetingInvitationService",
		Flags: []cli.Flag{
			&cli.IntFlag{Name: "http.port", Value: 8088, Usage: "listen port", Destination: &conf.Http.Port},
			&cli.StringFlag{Name: "http.dns", Value: "api.web.cn", Usage: "http dns", Destination: &conf.Http.DNS},
			&cli.StringFlag{Name: "http.savepath", Value: "/mnt/bucket/", Usage: "http file save path", Destination: &conf.Http.SavePath},

			&cli.StringFlag{Name: "db.user", Value: "root", Usage: "db user", Destination: &conf.DB.User},
			&cli.StringFlag{Name: "db.pwd", Value: "123456qaz", Usage: "db pwd", Destination: &conf.DB.PWD},
			&cli.StringFlag{Name: "db.addr", Value: "127.0.0.1:3306", Usage: "db addr", Destination: &conf.DB.Addr},
			&cli.StringFlag{Name: "db.name", Value: "mis", Usage: "database name", Destination: &conf.DB.Name},

			&cli.StringFlag{Name: "log.level", Value: "debug", Usage: "log level", Destination: &conf.Log.Level},

			&cli.StringFlag{Name: "wx.appid", Value: "wx456789", Usage: "wx appid", Destination: &conf.WX.AppId},
			&cli.StringFlag{Name: "wx.secret", Value: "sdagfjhkjkfsgfgjhjk.fgh", Usage: "wx secret", Destination: &conf.WX.Secret},
			&cli.StringFlag{Name: "wx.state", Value: "formal", Usage: "wx state", Destination: &conf.WX.State},
		},
		Action: func(c *cli.Context) error {

			app := app.NewApp(conf)
			if app == nil {
				return nil
			}
			err := app.Run()

			return err
		},

		Authors: []*cli.Author{{Name: "David Lee", Email: "chukuang2004@163.com"}},
	}

	app.Run(os.Args)

}
