package qrgen

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"log"
	"os"
)

func QrFromUrl(url string, hash string) {
	qrcode, err := qr.Encode(url, qr.M, qr.Auto)
	if err != nil {
		log.Fatalf(err.Error())
	}
	qrcode, err = barcode.Scale(qrcode, 300, 300)
	file, err := os.Create("./storage/qr/" + hash + ".png")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	err = png.Encode(file, qrcode)
	if err != nil {
		log.Fatal(err)
		return
	}
}
