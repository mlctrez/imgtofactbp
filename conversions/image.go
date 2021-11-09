package conversions

import (
	"image"

	"github.com/mlctrez/edgeefy"
	"github.com/nfnt/resize"
)

func ResizeWidth(original image.Image, width int) (img image.Image) {
	return resize.Resize(uint(width), 0, original, resize.Lanczos3)
}

func GrayScale(original image.Image) (img image.Image, err error) {
	pixels, err := edgeefy.GrayPixelsFrommImage(original)
	if err != nil {
		return
	}
	//pixels = edgeefy.CannyEdgeDetect(pixels, false, .5, .1)
	return edgeefy.GrayImageFromGrayPixels(pixels), nil
}
