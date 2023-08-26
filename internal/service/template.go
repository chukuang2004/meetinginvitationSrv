package service

import (
	"beastpark/meetinginvitationservice/internal/model/DB"
	"beastpark/meetinginvitationservice/internal/model/rpc"
	"beastpark/meetinginvitationservice/pkg/database"
	"beastpark/meetinginvitationservice/pkg/http"
	"beastpark/meetinginvitationservice/pkg/log"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Template struct {
	db  *database.DB
	log *log.Logger
}

func NewTemplate() *Template {

	t := &Template{
		db:  database.GetInstance(),
		log: log.GetInstance(),
	}

	if !t.db.Migrator().HasTable(&DB.Template{}) {
		t.db.Migrator().CreateTable(&DB.Template{})
	}

	return t
}

func (t *Template) Templates(ctx *gin.Context) {

	req := &rpc.TemplatesReq{}

	if err := ctx.BindQuery(req); err != nil {
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	templates := []*DB.Template{}
	if err := t.db.Find(&templates).Error; err != nil {
		t.log.Errorf("read db table@template error, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	resp := &rpc.TemplatesResp{}

	for i := 0; i < len(templates); i++ {
		t := templates[i]
		resp.Templates = append(resp.Templates, &rpc.Template{Name: t.Name, Data: t.Data})
	}

	http.WithResponse(ctx, http.Success, resp)
}

func (t *Template) UploadTemplates(ctx *gin.Context) {
	req := &rpc.UploadTemplatesReq{}

	if err := ctx.BindJSON(req); err != nil {
		t.log.Debugln("upload template error ", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	if req.Name == "" && req.Data == nil {
		t.log.Debugln("upload template params empty ")
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	tmp, err := json.Marshal(req.Data)
	if err != nil {
		t.log.Debugln("upload template params error, ", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	temp := DB.Template{
		Name: req.Name,
	}

	database.GetInstance().Where(DB.Template{Name: req.Name}).First(&temp)
	temp.Data = string(tmp)

	if err := t.db.Save(&temp).Error; err != nil {
		t.log.Errorf("save signIn relation fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	http.WithResponse(ctx, http.Success, nil)
}

func (t *Template) Delete(ctx *gin.Context) {

	req := &rpc.DeleteTemplatesReq{}

	if err := ctx.BindJSON(req); err != nil {
		t.log.Debugln("delete req param fail, ", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	info := DB.Template{
		ID: req.TemplateID,
	}

	if err := t.db.Delete(info).Error; err != nil {
		t.log.Errorf("delete meeting fail, %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	t.log.Debugln("start delete template ", req.TemplateID)

	http.WithResponse(ctx, http.Success, nil)
}
