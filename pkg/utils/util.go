package utils

import (
	"beastpark/meetinginvitationservice/pkg/log"
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"os"

	"github.com/nfnt/resize"
	uuid "github.com/satori/go.uuid"
	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/image/draw"
)

const (
	DefaultThumbnailWidth  = 400
	DefaultThumbnailHeight = 500
	DefaultQRSize          = 256
)

func GenUUID() string {
	id := uuid.NewV4()
	return id.String()
}

func GenThumbnail(root, src, dst string, width, height uint) error {

	if width == 0 || height == 0 {
		width = DefaultThumbnailWidth
		height = DefaultThumbnailHeight
	}

	file, err := os.Open(root + src)
	if err != nil {
		log.GetInstance().Debugln(err.Error())
		return err
	}

	// decode jpeg into image.Image
	img, err := png.Decode(file)
	if err != nil {
		log.GetInstance().Debugln(err.Error())
		return err
	}
	file.Close()

	outImg := resize.Thumbnail(width, height, img, resize.Lanczos3)

	name := dst

	out, err := os.Create(root + name)
	if err != nil {
		log.GetInstance().Debugln(err.Error())
		return err
	}
	defer out.Close()

	// write new image to file
	err = png.Encode(out, outImg)
	if err != nil {
		log.GetInstance().Debugln(err.Error())
		return err
	}

	return nil
}

// GenQRCodeWithLogo 带logo的二维码图片生成 content-二维码内容   level-容错级别,Low,Medium,High,Highest   size-像素单位
func GenQRCodeWithLogo(content string, size int) (image.Image, error) {

	//sets a fixed image width and height (e.g. 256 yields an 256x256px image).
	// size := 256

	code, err := qrcode.New(content, qrcode.Highest)
	if err != nil {
		return nil, err
	}
	//设置文件大小并创建画板
	qrcodeImg := code.Image(size)
	outImg := image.NewRGBA(qrcodeImg.Bounds())

	buf, _ := base64.RawStdEncoding.DecodeString(logo)
	logoImg, _, err := image.Decode(bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	logoImg = resize.Resize(uint(size/6), uint(size/6), logoImg, resize.Lanczos3)

	//logo和二维码拼接
	draw.Draw(outImg, outImg.Bounds(), qrcodeImg, image.Pt(0, 0), draw.Over)
	offset := image.Pt((outImg.Bounds().Max.X-logoImg.Bounds().Max.X)/2, (outImg.Bounds().Max.Y-logoImg.Bounds().Max.Y)/2)
	draw.Draw(outImg, outImg.Bounds().Add(offset), logoImg, image.Pt(0, 0), draw.Over)

	return outImg, nil
}

func GenQRCodeWithLogo4Base64(content string, size int) (string, error) {

	outImg, err := GenQRCodeWithLogo(content, size)
	if err != nil {
		return "", err
	}

	str := Png2Base64(outImg)

	return str, nil
}

func Png2Base64(img image.Image) string {

	data := bytes.NewBuffer([]byte(""))
	png.Encode(data, img)
	return base64.RawStdEncoding.EncodeToString(data.Bytes())
}

func PicMerge(posterImg image.Image, qrImg image.Image) image.Image {

	outImg := image.NewRGBA(posterImg.Bounds())

	draw.Draw(outImg, outImg.Bounds(), posterImg, image.Pt(0, 0), draw.Over)
	offset := image.Pt((outImg.Bounds().Max.X - qrImg.Bounds().Max.X), (outImg.Bounds().Max.Y - qrImg.Bounds().Max.Y))
	draw.Draw(outImg, outImg.Bounds().Add(offset), qrImg, image.Pt(0, 0), draw.Over)

	return outImg
}

func PosterWithQR(root, src string, content string, size int) error {

	file, err := os.Open(root + src)
	if err != nil {
		log.GetInstance().Debugln(err.Error())
		return err
	}

	// decode jpeg into image.Image
	poster, err := png.Decode(file)
	if err != nil {
		log.GetInstance().Debugln(err.Error())
		return err
	}
	file.Close()

	qr, err := GenQRCodeWithLogo(content, size)
	if err != nil {
		return err
	}

	outImg := PicMerge(poster, qr)
	out, err := os.Create(root + src)
	if err != nil {
		return err
	}
	defer out.Close()

	err = png.Encode(out, outImg)
	if err != nil {
		log.GetInstance().Debugln(err.Error())
		return err
	}

	return nil
}
