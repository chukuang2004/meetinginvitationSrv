package service

import (
	"beastpark/meetinginvitationservice/internal/model/DB"
	"beastpark/meetinginvitationservice/pkg/config"
	"beastpark/meetinginvitationservice/pkg/http"
	"beastpark/meetinginvitationservice/pkg/log"
	"bytes"
	"fmt"
	"time"
)

type Weixin struct {
	log         *log.Logger
	conf        *config.WXConfig
	accessToken string
	done        chan bool
}

var w *Weixin = nil

func NewWeixin(c *config.WXConfig) *Weixin {

	if w == nil {
		w = &Weixin{
			log:  log.GetInstance(),
			conf: c,
			done: make(chan bool),
		}
	}

	go func() {
		w.updateWXAccessToken()
	}()

	return w
}

func GetWeixin() *Weixin {

	if w == nil {
		log.GetInstance().Debugln("weixin instance is not exist")
	}

	return w
}

func (w *Weixin) Push2Users(users []string, meeting *DB.Meeting) {

	for _, id := range users {
		w.sendWXNotify(w.accessToken, id, meeting)
	}
}

func (w *Weixin) sendWXNotify(token, toOpenID string, meeting *DB.Meeting) {

	type Value struct {
		Val string `json:"value"`
	}

	type Notify struct {
		Touser      string `json:"touser"`
		Template_id string `json:"template_id"`
		Page        string `json:"page"`
		State       string `json:"miniprogram_state"`
		Lang        string `json:"lang"`
		Date        struct {
			Thing4    Value `json:"thing4"`
			Thing11   Value `json:"thing11"`
			Thing6    Value `json:"thing6"`
			Date3     Value `json:"date3"`
			Character Value `json:"character_string17"`
		} `json:"data"`
	}

	invitation := &DB.InvitationTemplate{}
	err := invitation.Decode(bytes.NewBufferString(meeting.Data).Bytes())
	if err != nil {
		w.log.Debugln("meeting invitation decode fail, ", err.Error())
		return
	}
	req := &Notify{
		Touser:      toOpenID,
		Template_id: "tbzvNdqleMmJbD2T5CM9QMgkwVXWMzh77H56WqCerfQ",
		Page:        "pages/preview/preview?meetingID=" + meeting.ID,
		State:       w.conf.State,
		Lang:        "zh_CN",
	}
	req.Date.Thing4.Val = meeting.Name
	req.Date.Thing11.Val = "xxxxx"
	req.Date.Thing6.Val = meeting.Place
	req.Date.Date3.Val = meeting.StartTime.Format("2006-10-12 12:23:12")

	type Resp struct {
		MsgID   int    `json:"msgid"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	resp := &Resp{}
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s", token)

	for i := 0; i < 3; i++ {
		err := http.Post(url, req, resp)
		if err != nil {
			w.log.Warning("send wx msg fail, ", err.Error())
		} else if resp.ErrCode != 0 {
			w.log.Warning("send wx msg fail, ", resp.ErrMsg)
		} else {
			break
		}
	}
}

func (w *Weixin) updateWXAccessToken() {

	token, expired := w.getWXAccessToken(w.conf.AppId, w.conf.Secret)

	w.accessToken = token

	ticker := time.NewTicker(time.Duration(expired) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.done:
			w.log.Infoln("quit updateWXAccessToken")
			return
		case <-ticker.C:
			token, expired = w.getWXAccessToken(w.conf.AppId, w.conf.Secret)
			w.accessToken = token
			ticker.Reset(time.Duration(expired) * time.Second)
		}
	}

}

func (w *Weixin) getWXAccessToken(appid, secret string) (string, int) {

	// GET https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		appid, secret)
	type Resp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	var resp Resp
	for i := 0; i < 3; i++ {
		err := http.Get(url, &resp)
		if err != nil {
			w.log.Warning("get wx access token fail, ", err.Error())
		} else if resp.ErrCode != 0 {
			w.log.Warning("get wx access token fail, ", resp.ErrMsg)
		} else {
			break
		}
	}

	return resp.AccessToken, resp.ExpiresIn
}

func (w *Weixin) GetWXOpenID(appid, secret, code string) (string, string) {

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appid, secret, code)
	type Resp struct {
		Openid     string `json:"openid"`
		SessionKey string `json:"session_key"`
		Unionid    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	var resp Resp
	for i := 0; i < 3; i++ {
		err := http.Get(url, &resp)
		if err != nil {
			w.log.Warning("get wx openid fail, ", err.Error())
		} else if resp.ErrCode != 0 {
			w.log.Warning("get wx openid fail, ", resp.ErrMsg)
		} else {
			break
		}
	}

	w.log.Debugf("get wx openid %s,unionid %s", resp.Openid, resp.Unionid)

	return resp.Openid, resp.Unionid
}
