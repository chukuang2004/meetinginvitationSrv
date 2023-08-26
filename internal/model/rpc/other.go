package rpc

type UploadFileReq struct {
	Openid string `json:"openid"`
}

type UploadFileResp struct {
	Url string `json:"url"`
}
