package service

import (
	"beastpark/meetinginvitationservice/internal/model/DB"
	"beastpark/meetinginvitationservice/pkg/http"
	"beastpark/meetinginvitationservice/pkg/log"
	"encoding/json"
)

type Payload struct {
	Name  string `json:"name"`
	Place string `json:"place"`
	Pic   string `json:"pic"`
	Time  string `json:"time"`
}

type Poster struct {
	Payload string `json:"payload"`
	UUID    string `json:"uuid"`
}

func NewPoster(id string, i *DB.InvitationPage) *Poster {

	if i.Index != 1 {
		log.GetInstance().Debugln("not first page")
		return nil
	}

	data, err := newPayload(i).Encode()
	if err != nil {
		log.GetInstance().Debugln("get poster payload fail,", err.Error())
		return nil
	}

	p := &Poster{
		Payload: string(data),
		UUID:    id,
	}

	return p
}

func (p *Poster) Get() string {

	url := "http://127.0.0.1:9900/poster"

	reqHeader := make(map[string]string)
	reqHeader["Content-Type"] = "application/json"
	reqHeader["token"] = "wertyuicvbnm,.fghjklghj"

	type Resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Url string `json:"url"`
		} `json:"data"`
	}
	resp := &Resp{}
	err := http.CustomPost(url, reqHeader, p, resp)
	if err != nil {
		log.GetInstance().Errorln("get poster fail, ", err.Error())
	}

	return resp.Data.Url
}

func newPayload(i *DB.InvitationPage) *Payload {

	p := &Payload{}
	for _, t := range i.Text {
		switch t.Tag {
		case "name":
			p.Name = t.Text
		case "place":
			p.Place = t.Text
		case "time":
			p.Time = t.Text
		}
	}

	for _, img := range i.Img {
		if img.Tag == "pic" {
			p.Pic = img.URL
		}
	}

	return p
}

func (p *Payload) Encode() ([]byte, error) {

	data, err := json.Marshal(p)

	return data, err
}
