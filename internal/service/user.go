package service

import (
	"beastpark/meetinginvitationservice/internal/model/DB"
	"beastpark/meetinginvitationservice/internal/model/rpc"
	"beastpark/meetinginvitationservice/pkg/config"
	"beastpark/meetinginvitationservice/pkg/database"
	"beastpark/meetinginvitationservice/pkg/http"
	"beastpark/meetinginvitationservice/pkg/log"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	db  *database.DB
	log *log.Logger
}

func NewUser() *User {

	u := &User{
		db:  database.GetInstance(),
		log: log.GetInstance(),
	}

	if !u.db.Migrator().HasTable(&DB.User{}) {
		u.db.Migrator().CreateTable(&DB.User{})
	}

	return u
}

func (u *User) Login(ctx *gin.Context) {

	req := &rpc.LoginReq{}

	if err := ctx.BindJSON(req); err != nil {
		u.log.Errorf("login params error, %s", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	conf, ok := ctx.MustGet("config").(*config.Config)
	if !ok {
		u.log.Errorln("get config fail")
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	openId, unionID := GetWeixin().GetWXOpenID(conf.WX.AppId, conf.WX.Secret, req.Wxcode)
	if openId == "" {
		u.log.Warningln("get wx openid fail")
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	user := DB.User{OpenID: openId, UnionID: unionID}
	if num := database.GetInstance().First(&user).RowsAffected; num == 0 {
		user.EnableMeetingNotify = true
		user.EnableMsgPush = true

		if err := u.db.Save(&user).Error; err != nil {
			u.log.Errorf("update user openid fail, %s", err.Error())
			http.WithResponse(ctx, http.UnknownError, nil)
			return
		}

	}

	resp := &rpc.LoginResp{
		Openid:              openId,
		WXName:              user.WXName,
		WXAvatar:            user.WXAvatar,
		EnableMeetingNotify: user.EnableMeetingNotify,
		EnableMsgPush:       user.EnableMsgPush,
	}

	http.WithResponse(ctx, http.Success, resp)
}

func (u *User) Notify(ctx *gin.Context) {

	req := &rpc.NotifyReq{}

	if err := ctx.BindJSON(req); err != nil {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}
	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	user := DB.User{
		OpenID: req.Openid,
	}
	if err := u.db.Model(&user).Updates(map[string]interface{}{"enable_msg_push": req.EnableMsgPush, "enable_meeting_notify": req.EnableMeetingNotify}).Error; err != nil {
		log.GetInstance().Errorf("db update notify fail %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	http.WithResponse(ctx, http.Success, nil)
}

func (u *User) Push(ctx *gin.Context) {
	req := &rpc.PushReq{}

	if err := ctx.BindJSON(req); err != nil {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	meeting, ok := CheckMeetingIdValid(ctx, req.MeetingID)
	if !ok {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !meeting.EndTime.After(time.Now()) {
		u.log.Infoln("can't push when meeting is finish")
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	users, err := GetUserByMeeting(req.MeetingID)
	if err != nil {
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	GetWeixin().Push2Users(users, meeting)

	http.WithResponse(ctx, http.Success, nil)
}

func (u *User) Update(ctx *gin.Context) {

	req := &rpc.UserUpdateReq{}

	if err := ctx.BindJSON(req); err != nil {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	user := DB.User{OpenID: req.Openid}
	if num := database.GetInstance().First(&user).RowsAffected; num == 1 {

		user.WXName = req.WXName
		user.WXAvatar = req.WXAvatar

		if err := u.db.Save(&user).Error; err != nil {
			u.log.Errorf("update user openid fail, %s", err.Error())
			http.WithResponse(ctx, http.UnknownError, nil)
			return
		}

	}

	http.WithResponse(ctx, http.Success, nil)
}

func CheckOpenIdValid(ctx *gin.Context, openid string) bool {

	user := DB.User{OpenID: openid}
	if num := database.GetInstance().First(&user).RowsAffected; num == 0 {
		log.GetInstance().Errorln("user Openid is not exist")
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return false
	}

	return true
}

func GetUserByMeeting(id string) ([]string, error) {

	results := []DB.Relation{}
	err := database.GetInstance().Where(&DB.Relation{MeetingID: id}).Distinct("open_id").Find(&results).Error
	if err != nil {
		log.GetInstance().Errorf("get meeting id from db fail, %s", err.Error())
		return nil, err
	}

	ids := []string{}
	for _, item := range results {
		ids = append(ids, item.OpenID)
	}

	return ids, nil
}
