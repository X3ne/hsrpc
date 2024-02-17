package utils

import (
	"github.com/X3ne/hsrpc/src/logger"
)

// This function check the brightness of the pixels in the given coordinates to determine the current character
func FindCurrentCharacter(coords []Rect) int32 {
	whitestPosition := -1
	whitestValue := 0
	minBrightness := 450 // TODO: determine this value to prevent false results

	for i, coord := range coords {
		pixel, err := GetPixelColor(&coord) // TODO: replace this function by this one https://pkg.go.dev/github.com/go-vgo/robotgo#GetPixelColor
		if err != nil {
			logger.Logger.Error(err)
			continue
		}

		brightness := int(pixel.R) + int(pixel.G) + int(pixel.B)

		if brightness >= minBrightness && brightness > whitestValue {
			whitestValue = brightness
			whitestPosition = i
		}
	}

	return int32(whitestPosition)
}
