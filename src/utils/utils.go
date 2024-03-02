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
	"github.com/lxn/win"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"golang.org/x/sys/windows/registry"
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

func GetPixelColor(hWnd win.HWND, rect Rect) (color.RGBA, error) {
	var winRect win.RECT
	win.GetWindowRect(hWnd, &winRect)

	img := robotgo.CaptureImg(rect.X + int(winRect.Left), rect.Y + int(winRect.Top), 100, 1)
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
	appDataDir, err := GetAppPath()
	if err != nil {
		fmt.Println("Error getting user config dir:", err)
		return
	}
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("crash_%s.log", timestamp)
	logFilePath := filepath.Join(appDataDir, consts.CrashDir, filename)

	err = CreateDir(logFilePath)
	if err != nil {
		logger.Logger.Error(err)
		return
	}

	file, err := os.Create(logFilePath)
	if err != nil {
		logger.Logger.Error("error creating crash log file:", err)
		return
	}
	defer file.Close()

	buf := make([]byte, 1<<16)
	stackLen := runtime.Stack(buf, true)
	fmt.Fprintf(file, "=== STACK TRACE ===\n%s\n", buf[:stackLen])
}

func GetAppPath() (string, error) {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("error getting user config dir: %w", err)
	}

	return filepath.Join(appDataDir, consts.AppDataDir), nil
}

func CreateDir(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0750)
		if err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
	}

	return nil
}

func GetWindowsAccentColor() (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Accent`, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("error opening registry key: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetIntegerValue("StartColorMenu")
	if err != nil {
		return "", fmt.Errorf("error querying registry value: %w", err)
	}

	colorHex := fmt.Sprintf("#%06X", value&0xFFFFFF)

	return colorHex, nil
}
