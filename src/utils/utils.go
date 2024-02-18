package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/X3ne/hsrpc/src/consts"
	"github.com/X3ne/hsrpc/src/logger"
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
	img := robotgo.CaptureImg(rect.X, rect.Y, 1, 1)
	if img == nil {
		return color.RGBA{}, nil
	}

	pixelColor := img.At(0, 0).(color.RGBA)
	return pixelColor, nil
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

func PanicRecover(r interface{}) {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config dir:", err)
		return
	}
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("crash_%s.log", timestamp)
	logFilePath := filepath.Join(appDataDir, consts.AppDataDir, "crash", filename)

	dir := filepath.Dir(logFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Logger.Error("Error creating crash log directory:", err)
			return
		}
	}

	file, err := os.Create(logFilePath)
	if err != nil {
		logger.Logger.Error("Error creating crash log file:", err)
		return
	}
	defer file.Close()

	buf := make([]byte, 1<<16)
	stackLen := runtime.Stack(buf, true)
	fmt.Fprintf(file, "=== STACK TRACE ===\n%s\n", buf[:stackLen])
}
