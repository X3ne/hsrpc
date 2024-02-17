package utils

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/go-vgo/robotgo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type Rect struct {
	X, Y, Width, Height	int
}

// This function is useful to mitigate OCR errors by finding the closest correspondence to the given text
func FindClosestCorrespondence(text string, candidates []Data) Data {
	const threshold = 5

	minDistance := len(text)
	var closest Data

	for _, candidate := range candidates {
		distance := levenshtein.DistanceForStrings([]rune(text), []rune(candidate.Value), levenshtein.DefaultOptions)
		if distance < minDistance {
			minDistance = distance
			closest = candidate
		}
	}

	if minDistance > threshold {
		return Data{}
	}

	return closest
}

func GetPixelColor(rect *Rect) (color.RGBA, error) {
	img := robotgo.CaptureImg(rect.X, rect.Y, rect.Width, rect.Height)
	if img == nil {
		return color.RGBA{}, nil
	}

	pixelColor := img.At(0, 0).(color.RGBA)
	return pixelColor, nil
}

func ConvertToGrayscale(src image.Image) image.Image {
	bounds := src.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := src.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor).(color.Gray)
			grayImg.Set(x, y, grayColor)
		}
	}

	return grayImg
}

func SaveImg(img image.Image, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
			return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
			return err
	}

	return nil
}

