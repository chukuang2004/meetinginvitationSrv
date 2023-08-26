package app

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"beastpark/meetinginvitationservice/internal/service"
	"beastpark/meetinginvitationservice/pkg/http"
)

type HttpSrv struct {
	*gin.Engine
	port int

	handlerUser     *service.User
	handlerMeeting  *service.Meeting
	handlerOther    *service.Other
	handlerTemplate *service.Template
}

func SetConfig(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", app.conf)
		c.Next()
	}
}

func NewServer(app *App) *HttpSrv {

	srv := &HttpSrv{port: app.conf.Http.Port}

	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(SetConfig(app))
	e.Use(http.LogHandler)

	srv.Engine = e
	srv.handlerUser = service.NewUser()
	srv.handlerMeeting = service.NewMeeting()
	srv.handlerOther = service.NewOther()
	srv.handlerTemplate = service.NewTemplate()

	service.NewWeixin(&app.conf.WX)

	srv.routes()

	return srv
}

func (srv *HttpSrv) Run() error {

	addr := fmt.Sprintf("0.0.0.0:%d", srv.port)

	return srv.Engine.Run(addr)
}

func (srv *HttpSrv) routes() {

	user := srv.Group("user")
	user.POST("login", srv.handlerUser.Login)
	user.POST("notify", srv.handlerUser.Notify)
	user.POST("push", srv.handlerUser.Push)
	user.POST("update", srv.handlerUser.Update)

	template := srv.Group("template")
	template.POST("upload", srv.handlerTemplate.UploadTemplates)
	template.GET("all", srv.handlerTemplate.Templates)

	meeting := srv.Group("meeting")
	meeting.GET("query", srv.handlerMeeting.Query)
	meeting.POST("create", srv.handlerMeeting.Create)
	meeting.POST("create/noposter", srv.handlerMeeting.CreateNoPoster)
	meeting.DELETE("delete", srv.handlerMeeting.Delete)
	meeting.PUT("update", srv.handlerMeeting.Update)
	meeting.POST("view", srv.handlerMeeting.View)
	meeting.GET("following", srv.handlerMeeting.Following)
	meeting.GET("history", srv.handlerMeeting.History)
	meeting.GET("QR4invite", srv.handlerMeeting.QR4invite)
	meeting.POST("signUp", srv.handlerMeeting.SignUp)
	meeting.POST("quit", srv.handlerMeeting.Quit)
	meeting.GET("report", srv.handlerMeeting.Report)
	meeting.GET("QR4signIn", srv.handlerMeeting.QR4signIn)
	meeting.POST("signIn", srv.handlerMeeting.SignIn)

	srv.POST("uploadFile", srv.handlerOther.UploadFile)
}
