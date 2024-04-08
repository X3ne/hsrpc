package utils

import (
	"github.com/X3ne/hsrpc/src/logger"
	"github.com/lxn/win"
)

// FindCurrentCharacter
// This function check the brightness of the pixels in the given coordinates to determine the current character
func FindCurrentCharacter(hWnd win.HWND, coords []Rect) int32 {
	whitestPosition := -1
	whitestValue := 0

	for i, coord := range coords {
		pixel, err := GetPixelColor(hWnd, coord)
		if err != nil {
			logger.Logger.Error(err)
			continue
		}

		brightness := int(pixel.R) + int(pixel.G) + int(pixel.B)

		if brightness > whitestValue {
			whitestValue = brightness
			whitestPosition = i
		}
	}

	return int32(whitestPosition)
}
