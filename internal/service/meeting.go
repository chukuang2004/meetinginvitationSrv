package service

import (
	"beastpark/meetinginvitationservice/internal/model/DB"
	"beastpark/meetinginvitationservice/internal/model/rpc"
	"beastpark/meetinginvitationservice/pkg/config"
	"beastpark/meetinginvitationservice/pkg/database"
	"beastpark/meetinginvitationservice/pkg/http"
	"beastpark/meetinginvitationservice/pkg/log"
	"beastpark/meetinginvitationservice/pkg/utils"
	"bytes"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Meeting struct {
	db  *database.DB
	log *log.Logger
}

func NewMeeting() *Meeting {

	m := &Meeting{
		db:  database.GetInstance(),
		log: log.GetInstance(),
	}

	if !m.db.Migrator().HasTable(&DB.Meeting{}) {
		m.db.Migrator().CreateTable(&DB.Meeting{})
	}
	if !m.db.Migrator().HasTable(&DB.Relation{}) {
		m.db.Migrator().CreateTable(&DB.Relation{})
	}
	if !m.db.Migrator().HasTable(&DB.Template{}) {
		m.db.Migrator().CreateTable(&DB.Template{})
	}

	return m
}

func (m *Meeting) Query(ctx *gin.Context) {
	req := &rpc.QueryReq{}

	if err := ctx.BindQuery(req); err != nil {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	meeting, ok := CheckMeetingIdValid(ctx, req.MeetingID)
	if !ok {
		return
	}

	info := &rpc.MeetingInfo{
		ID:        meeting.ID,
		Name:      meeting.Name,
		Place:     meeting.Place,
		Poster:    meeting.Poster,
		Thumbnail: meeting.Thumbnail,
		Meeting:   meeting.Data,
		StartTime: meeting.StartTime.String(),
		EndTime:   meeting.EndTime.String(),
	}

	resp := &rpc.QueryResp{
		Meeting: info,
	}
	http.WithResponse(ctx, http.Success, resp)
}

func (m *Meeting) Create(ctx *gin.Context) {

	req := &rpc.CreateReq{}

	if err := ctx.BindJSON(req); err != nil {
		m.log.Debugln("create req error ", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	array := strings.Split(req.Poster, "/")
	size := len(array)
	if size < 3 {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	conf, ok := ctx.MustGet("config").(*config.Config)
	if !ok {
		m.log.Errorf("get config fail")
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}
	file := array[size-1]

	thumbnail := utils.GenUUID() + ".png"
	go func(name string) {
		err := utils.GenThumbnail(conf.Http.SavePath, file, name, uint(0), uint(0))
		if err != nil {
			m.log.Errorf("GenThumbnail fail, %s", err.Error())
			return
		}
		m.log.Infoln("GenThumbnail success, poster ", req.Poster)
	}(thumbnail)

	thumbnail = "https://" + conf.Http.DNS + "/download/" + thumbnail

	id := utils.GenUUID()
	meeting := &DB.Meeting{
		ID:        id,
		CreaterID: req.Openid,
		Name:      req.Name,
		Place:     req.Place,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Data:      req.Meeting,
		Poster:    req.Poster,
		Thumbnail: thumbnail,
	}

	if err := m.db.Save(meeting).Error; err != nil {
		m.log.Errorf("save meeting fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	r := &DB.Relation{
		OpenID:    req.Openid,
		MeetingID: id,
		Action:    "create",
		EndTime:   meeting.EndTime,
	}

	if err := m.db.Save(r).Error; err != nil {
		m.log.Errorf("save create relation fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	go func() {
		if err := utils.PosterWithQR(conf.Http.SavePath, file, "http://"+conf.Http.DNS+"/mis/signUp?meetingID="+id, utils.DefaultQRSize); err != nil {
			m.log.Errorf("add qr into poster fail, %s", err.Error())
			return
		}
		m.log.Infoln("PosterWithQR success, poster ", req.Poster)
	}()

	resp := &rpc.CreateResp{}
	resp.ID = id
	http.WithResponse(ctx, http.Success, resp)
}

func (m *Meeting) CreateNoPoster(ctx *gin.Context) {

	req := &rpc.CreateReq{}

	if err := ctx.BindJSON(req); err != nil {
		m.log.Debugln("create req error ", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	var temp DB.InvitationTemplate
	err := temp.Decode(bytes.NewBufferString(req.Meeting).Bytes())
	if err != nil {
		m.log.Debugln("decode InvitationTemplate fail, ", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	for _, page := range temp {
		if page.Index == 1 {
			poster := NewPoster("gfdhgjkkjhlj", page)

			req.Poster = poster.Get()
			m.log.Debugln("poster ", req.Poster)

			break
		}
	}

	array := strings.Split(req.Poster, "/")
	size := len(array)
	if size < 3 {
		m.log.Debugln("poster ", req.Poster)
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	conf, ok := ctx.MustGet("config").(*config.Config)
	if !ok {
		m.log.Errorf("get config fail")
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}
	file := array[size-1]

	thumbnail := utils.GenUUID() + ".png"
	go func(name string) {
		err := utils.GenThumbnail(conf.Http.SavePath, file, name, uint(0), uint(0))
		if err != nil {
			m.log.Errorf("GenThumbnail fail, %s", err.Error())
			return
		}
		m.log.Infoln("GenThumbnail success, poster ", req.Poster)
	}(thumbnail)

	thumbnail = "https://" + conf.Http.DNS + "/download/" + thumbnail

	id := utils.GenUUID()
	meeting := &DB.Meeting{
		ID:        id,
		CreaterID: req.Openid,
		Name:      req.Name,
		Place:     req.Place,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Data:      req.Meeting,
		Poster:    req.Poster,
		Thumbnail: thumbnail,
	}

	if err := m.db.Save(meeting).Error; err != nil {
		m.log.Errorf("save meeting fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	r := &DB.Relation{
		OpenID:    req.Openid,
		MeetingID: id,
		Action:    "create",
		EndTime:   meeting.EndTime,
	}

	if err := m.db.Save(r).Error; err != nil {
		m.log.Errorf("save create relation fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	go func() {
		if err := utils.PosterWithQR(conf.Http.SavePath, file, "http://"+conf.Http.DNS+"/mis/signUp?meetingID="+id, utils.DefaultQRSize/2); err != nil {
			m.log.Errorf("add qr into poster fail, %s", err.Error())
			return
		}
		m.log.Infoln("PosterWithQR success, poster ", req.Poster)
	}()

	type CreateResp struct {
		ID  string `json:"ID"`
		Url string `json:"poster"`
	}

	resp := &CreateResp{
		ID:  id,
		Url: req.Poster,
	}

	http.WithResponse(ctx, http.Success, resp)
}

func (m *Meeting) Delete(ctx *gin.Context) {

	req := &rpc.DeleteReq{}

	if err := ctx.BindJSON(req); err != nil {
		m.log.Debugln("delete req param fail, ", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	info, ok := CheckMeetingIdValid(ctx, req.MeetingID)
	if !ok {
		return
	}

	if info.CreaterID != req.Openid {
		m.log.Debugln("don't have update right")
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if err := m.db.Delete(info).Error; err != nil {
		m.log.Errorf("delete meeting fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	m.log.Debugln("start delete meeting ", req.MeetingID)

	http.WithResponse(ctx, http.Success, nil)
}

func (m *Meeting) Update(ctx *gin.Context) {

	req := &rpc.UpdateReq{}

	if err := ctx.BindJSON(req); err != nil {
		m.log.Debugln("update req param fail, ", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	info, ok := CheckMeetingIdValid(ctx, req.MeetingID)
	if !ok {
		return
	}

	if info.CreaterID != req.Openid {
		m.log.Debugln("don't have update right")
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	array := strings.Split(req.Poster, "/")
	size := len(array)
	if size < 3 {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	conf, ok := ctx.MustGet("config").(*config.Config)
	if !ok {
		m.log.Errorf("get config fail")
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}
	file := array[size-1]

	thumbnail := utils.GenUUID() + ".png"
	go func(name string) {
		err := utils.GenThumbnail(conf.Http.SavePath, file, name, uint(0), uint(0))
		if err != nil {
			m.log.Errorf("GenThumbnail fail, %s", err.Error())
			return
		}
		m.log.Infoln("GenThumbnail success, poster ", req.Poster)
	}(thumbnail)

	thumbnail = "https://" + conf.Http.DNS + "/download/" + thumbnail

	meeting := &DB.Meeting{
		ID:        req.MeetingID,
		CreaterID: req.Openid,
		Name:      req.Name,
		Place:     req.Place,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Data:      req.Meeting,
		Poster:    req.Poster,
		Thumbnail: thumbnail,
	}

	if err := m.db.Omit("created_at").Save(meeting).Error; err != nil {
		m.log.Errorf("update meeting fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	go func() {
		if err := utils.PosterWithQR(conf.Http.SavePath, file, "http://"+conf.Http.DNS+"/mis/signUp?meetingID="+req.MeetingID, utils.DefaultQRSize); err != nil {
			m.log.Errorf("add qr into poster fail, %s", err.Error())
			return
		}
		m.log.Infoln("PosterWithQR success, poster ", req.Poster)
	}()

	http.WithResponse(ctx, http.Success, nil)
}

func (m *Meeting) Following(ctx *gin.Context) {

	req := &rpc.FollowingReq{}

	if err := ctx.BindQuery(req); err != nil {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	// // db.Distinct("name", "age").Order("name, age desc").Find(&results)
	// // db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)
	results := []DB.Relation{}
	err := m.db.Where(&DB.Relation{OpenID: req.Openid}).Where("end_time > ?", time.Now()).Distinct("meeting_id").Find(&results).Error
	if err != nil {
		m.log.Errorf("get meeting id from db fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	ids := []string{}
	for _, item := range results {
		ids = append(ids, item.MeetingID)
	}

	resp := &rpc.FollowingResp{}
	if len(ids) == 0 {
		http.WithResponse(ctx, http.Success, resp)
		return
	}

	meetings := []DB.Meeting{}
	if err = m.db.Find(&meetings, ids).Error; err != nil {
		m.log.Errorf("get meeting id from db fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	for _, meeting := range meetings {
		info := &rpc.MeetingInfo{
			ID:        meeting.ID,
			Name:      meeting.Name,
			Place:     meeting.Place,
			Poster:    meeting.Poster,
			Thumbnail: meeting.Thumbnail,
			Meeting:   meeting.Data,
			StartTime: meeting.StartTime.String(),
			EndTime:   meeting.EndTime.String(),
		}
		if meeting.CreaterID == req.Openid {
			info.IsOwner = true
		}
		resp.Meeting = append(resp.Meeting, info)
	}

	http.WithResponse(ctx, http.Success, resp)
}

func (m *Meeting) History(ctx *gin.Context) {

	req := &rpc.HistoryReq{}

	if err := ctx.BindQuery(req); err != nil {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	// db.Distinct("name", "age").Order("name, age desc").Find(&results)
	// db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)
	results := []DB.Relation{}
	err := m.db.Where(&DB.Relation{OpenID: req.Openid}).Distinct("meeting_id").Find(&results).Error
	if err != nil {
		m.log.Errorf("get meeting id from db fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	ids := []string{}
	for _, item := range results {
		ids = append(ids, item.MeetingID)
	}

	resp := &rpc.FollowingResp{}
	if len(ids) == 0 {
		http.WithResponse(ctx, http.Success, resp)
		return
	}

	meetings := []DB.Meeting{}
	if err = m.db.Find(&meetings, ids).Error; err != nil {
		m.log.Errorf("get meeting id from db fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	for _, meeting := range meetings {
		info := &rpc.MeetingInfo{
			ID:        meeting.ID,
			Name:      meeting.Name,
			Place:     meeting.Place,
			Poster:    meeting.Poster,
			Thumbnail: meeting.Thumbnail,
			Meeting:   meeting.Data,
			StartTime: meeting.StartTime.String(),
			EndTime:   meeting.EndTime.String(),
		}
		if meeting.CreaterID == req.Openid {
			info.IsOwner = true
		}
		resp.Meeting = append(resp.Meeting, info)
	}

	http.WithResponse(ctx, http.Success, resp)
}

// add into poster
func (m *Meeting) QR4invite(ctx *gin.Context) {

	req := &rpc.QR4inviteReq{}

	if err := ctx.BindQuery(req); err != nil {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	conf, ok := ctx.MustGet("config").(*config.Config)
	if !ok {
		m.log.Errorf("get config fail")
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	qr, err := utils.GenQRCodeWithLogo4Base64("http://"+conf.Http.DNS+"/mis/signUp?meetingID="+req.MeetingID, utils.DefaultQRSize)
	if err != nil {
		m.log.Debugln(err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	resp := &rpc.QR4SignInResp{
		QRCode: qr,
	}

	http.WithResponse(ctx, http.Success, resp)
}

func (m *Meeting) View(ctx *gin.Context) {
	req := &rpc.ViewReq{}

	if err := ctx.BindJSON(req); err != nil {
		m.log.Debugln("SignUp req param fail, ", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	m.log.Infoln(req)
	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	meetInfo, ok := CheckMeetingIdValid(ctx, req.MeetingID)
	if !ok {
		return
	}

	if !meetInfo.EndTime.After(time.Now()) {
		m.log.Debugln("SignUp meeting is expired")
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	r := &DB.Relation{
		OpenID:    req.Openid,
		MeetingID: req.MeetingID,
		Action:    "view",
		EndTime:   meetInfo.EndTime,
	}
	if err := m.db.Save(r).Error; err != nil {
		m.log.Errorf("save view relation fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	http.WithResponse(ctx, http.Success, nil)
}

func (m *Meeting) SignUp(ctx *gin.Context) {

	req := &rpc.SignUpReq{}

	if err := ctx.BindJSON(req); err != nil {
		m.log.Debugln("SignUp req param fail, ", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	m.log.Infoln(req)
	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	meetInfo, ok := CheckMeetingIdValid(ctx, req.MeetingID)
	if !ok {
		return
	}

	if !meetInfo.EndTime.After(time.Now()) {
		m.log.Debugln("SignUp meeting is expired")
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	r := &DB.Relation{
		OpenID:     req.Openid,
		MeetingID:  req.MeetingID,
		Action:     "signUp",
		Name:       req.Name,
		Phone:      req.PhoneNum,
		PartnerNum: req.PartnerNum,
		EndTime:    meetInfo.EndTime,
	}
	if err := m.db.Save(r).Error; err != nil {
		m.log.Errorf("save signup relation fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	http.WithResponse(ctx, http.Success, nil)
}

func (m *Meeting) Quit(ctx *gin.Context) {

	req := &rpc.QuitReq{}

	if err := ctx.BindJSON(req); err != nil {
		m.log.Debugln("SignUp req param fail, ", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	info, ok := CheckMeetingIdValid(ctx, req.MeetingID)
	if !ok {
		return
	}

	if info.CreaterID == req.Openid {
		m.log.Debugln("don't have quit right")
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if err := m.db.Where(DB.Relation{
		OpenID:    req.Openid,
		MeetingID: req.MeetingID,
	}).Delete(&DB.Relation{}).Error; err != nil {
		m.log.Errorf("delete signup relation fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	http.WithResponse(ctx, http.Success, nil)
}

func (m *Meeting) Report(ctx *gin.Context) {

	req := &rpc.ReportReq{}

	if err := ctx.BindQuery(req); err != nil {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	_, ok := CheckMeetingIdValid(ctx, req.MeetingID)
	if !ok {
		return
	}

	count := func(action string) int {

		var count int64
		err := m.db.Model(&DB.Relation{}).Where(&DB.Relation{MeetingID: req.MeetingID, Action: action}).Distinct("open_id").Count(&count).Error
		if err != nil {
			m.log.Errorf("get meeting id from db fail, %s", err.Error())
			return 0
		}

		return int(count)
	}

	resp := &rpc.ReportResp{
		View:   count("view"),
		SignUp: count("signUp"),
		SignIn: count("signIn"),
	}

	http.WithResponse(ctx, http.Success, resp)
}

func (m *Meeting) QR4signIn(ctx *gin.Context) {

	req := &rpc.QR4SignInReq{}

	if err := ctx.BindQuery(req); err != nil {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	conf, ok := ctx.MustGet("config").(*config.Config)
	if !ok {
		m.log.Errorf("get config fail")
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	qr, err := utils.GenQRCodeWithLogo4Base64("http://"+conf.Http.DNS+"/mis/signIn?meetingID="+req.MeetingID, utils.DefaultQRSize)
	if err != nil {
		m.log.Debugln(err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	resp := &rpc.QR4SignInResp{
		QRCode: qr,
	}
	http.WithResponse(ctx, http.Success, resp)
}

func (m *Meeting) SignIn(ctx *gin.Context) {

	req := &rpc.SignInReq{}

	if err := ctx.BindJSON(req); err != nil {
		m.log.Debugln("SignIn req param fail, ", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	meetInfo, ok := CheckMeetingIdValid(ctx, req.MeetingID)
	if !ok {
		return
	}

	r := &DB.Relation{
		OpenID:    req.Openid,
		MeetingID: req.MeetingID,
		Action:    "signIn",
		EndTime:   meetInfo.EndTime,
	}

	if err := m.db.Save(r).Error; err != nil {
		m.log.Errorf("save signIn relation fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	http.WithResponse(ctx, http.Success, nil)
}

func CheckMeetingIdValid(ctx *gin.Context, meetingID string) (*DB.Meeting, bool) {

	meeting := DB.Meeting{ID: meetingID}
	if num := database.GetInstance().First(&meeting).RowsAffected; num == 0 {
		log.GetInstance().Errorln("meeting id is not exist")
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return nil, false
	}

	return &meeting, true
}
