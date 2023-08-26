package utils

import (
	"encoding/base64"
	"image/png"
	"os"
	"testing"
)

// go test -v -run Test_GenQRCodeWithLogo util_test.go
func Test_GenQRCodeWithLogo(t *testing.T) {
	logo, err := GenQRCodeWithLogo4Base64("http://web.cn/mis/signUp?meetingID=4de79469-5ca5-4cd5-b51d-88ec887ab28a", 256*3)
	if err != nil {
		t.Errorf("%s\n", err.Error())
	}

	buf, _ := base64.RawStdEncoding.DecodeString(logo)
	os.WriteFile("signUp.png", buf, 0644)
}

func Test_PicMerge(t *testing.T) {

	file1, err := os.Open("poster.png")
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}

	// decode jpeg into image.Image
	poster, err := png.Decode(file1)
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	file1.Close()

	file2, err := os.Open("signUp.png")
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}

	// decode jpeg into image.Image
	qr, err := png.Decode(file2)
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	file2.Close()

	outImg := PicMerge(poster, qr)

	out, err := os.Create("merge.png")
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
	defer out.Close()

	err = png.Encode(out, outImg)
	if err != nil {
		t.Errorf("%s\n", err.Error())
		return
	}
}
