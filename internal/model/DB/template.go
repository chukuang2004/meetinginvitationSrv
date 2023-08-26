package DB

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Template struct {
	gorm.Model
	ID   int `gorm:"primary_key"`
	Name string
	Data string `gorm:"size:1024*1024"`
}

func (*Template) TableName() string {
	return "template"
}

type Text struct {
	Tag        string `json:"tag"` //name place time(start ~ end)
	X          int    `json:"x"`
	Y          int    `json:"y"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	BaseLine   string `json:"baseLine"`
	Center     bool   `json:"center"`
	Text       string `json:"text"`
	FontSize   int    `json:"fontSize"`
	LineHeight int    `json:"line-height"`
	TextIndent int    `json:"textIndent"`
	Color      string `json:"color"`
	TextAlign  string `json:"textAlign"`
}

type InvitationPage struct {
	Index      int     `json:"index"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	BGUrl      string  `json:"backgroundUrl"`
	PixelRatio float32 `json:"pixelRatio"`
	Img        []*struct {
		Tag          string `json:"tag"` //for poster
		X            int    `json:"x"`
		Y            int    `json:"y"`
		Width        int    `json:"width"`
		Height       int    `json:"height"`
		URL          string `json:"url"`
		BorderRadius string `json:"borderRadius"` //圆角
		ZIndex       int    `json:"zIndex"`       //上下堆叠顺序  大的在上面
	} `json:"images"`
	Video []*struct {
		X        int    `json:"x"`
		Y        int    `json:"y"`
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		URL      string `json:"url"`
		Autoplay bool   `json:"autoplay"` //
		ZIndex   int    `json:"zIndex"`   //上下堆叠顺序
	} `json:"videos"`
	Text []*struct {
		X            int    `json:"x"`
		Y            int    `json:"y"`
		Width        int    `json:"width"`
		Height       int    `json:"height"`
		BGColor      string `json:"backgroundColor"`
		BorderRadius string `json:"borderRadius"` //圆角
		Tag          string `json:"tag"`          //name place time(start ~ end)
		BaseLine     string `json:"baseLine"`
		TextAlign    string `json:"textAlign"`
		Text         string `json:"text"`
		FontSize     int    `json:"fontSize"`
		LineHeight   int    `json:"line-height"`
		TextIndent   int    `json:"textIndent"`
		Color        string `json:"color"`
	} `json:"texts"`
}

type InvitationTemplate []*InvitationPage

func (i *InvitationTemplate) Encode() ([]byte, error) {

	data, err := json.Marshal(i)
	return data, err
}

func (i *InvitationTemplate) Decode(data []byte) error {

	err := json.Unmarshal(data, i)

	return err
}
