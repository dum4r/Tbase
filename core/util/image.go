package util

import (
	"bytes"
	"embed"
	"image"
	"image/png"
)

var assets *embed.FS

func SetAssets(a *embed.FS) { assets = a }

func GetImage(filePath string) image.Image {
	file, err := assets.Open(filePath)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return img
}

func IcontoByte(img image.Image) (data []byte, width int, height int) {
	bounds := img.Bounds()
	width = bounds.Max.X
	height = bounds.Max.Y
	data = make([]byte, width*height*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			index := y*width*4 + x*4
			data[index] = byte(b >> 8)
			data[index+1] = byte(g >> 8)
			data[index+2] = byte(r >> 8)
			data[index+3] = byte(a >> 8)
		}
	}
	return
}

func ImgtoByte(img image.Image) (data []byte, width int, height int) {
	bounds := img.Bounds()
	width = bounds.Max.X
	height = bounds.Max.Y
	data = make([]byte, width*height*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			data[y*width*4+x*4+3] = byte(a >> 8)
			data[y*width*4+x*4+0] = byte(b >> 8)
			data[y*width*4+x*4+1] = byte(g >> 8)
			data[y*width*4+x*4+2] = byte(r >> 8)
		}
	}
	return
}

func ImgtoByte2(img image.Image) (data []byte, width int, height int) {
	bounds := img.Bounds()
	width = bounds.Max.X
	height = bounds.Max.Y
	data = ConvertImageToJPEG(img)
	return
}

func ConvertImageToJPEG(img image.Image) []byte {
	// Create a new buffer.
	buf := new(bytes.Buffer)

	// Encode the image to the buffer.
	err := png.Encode(buf, img)
	if err != nil {
		panic(err)
	}

	// Return the buffer as a byte slice.
	return buf.Bytes()
}
