package conversions

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"strings"
)

// Base64ToImage base64 encoded image data to image.Image.
func Base64ToImage(input string) (img image.Image, typ string, err error) {
	if !strings.HasPrefix(input, "data:image/") {
		err = errors.New("input must start with data:image/")
		return
	}
	var imageData []byte
	if imageData, err = Base64ToByte(input); err != nil {
		return
	} else {
		return image.Decode(bytes.NewBuffer(imageData))
	}
}

// Base64ToByte decodes the base64 data to a byte array
func Base64ToByte(input string) (result []byte, err error) {
	if !strings.HasPrefix(input, "data:") {
		err = errors.New("input must start with data: ")
		return
	}
	parts := strings.SplitN(input, ",", 2)
	if !strings.HasSuffix(parts[0], ";base64") {
		err = errors.New("data encoding must be base64")
		return
	}
	return base64.StdEncoding.DecodeString(parts[1])
}

func ImageToBase64(input image.Image) string {
	b := &bytes.Buffer{}
	b.WriteString("data:image/png;base64,")
	_ = png.Encode(base64.NewEncoder(base64.StdEncoding, b), input)
	return b.String()
}
