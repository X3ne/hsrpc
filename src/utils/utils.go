package utils

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

type PlayerInfo struct {
	Player struct {
		UID				string	`json:"uid"`
		Level			int			`json:"level"`
		Nickname	string	`json:"nickname"`
	} `json:"player"`
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

func PanicRecover(r interface{}, guiApp ...fyne.App) {
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
	fmt.Fprintf(file, "Error: %s\n\n=== STACK TRACE ===\n%s\n", r, buf[:stackLen])

	if len(guiApp) > 0 {
		app := guiApp[0]
		win := app.NewWindow("Application Error")

		win.SetIcon(app.Icon())

		errText := widget.NewLabel(fmt.Sprintf("An error occurred: %v\n\nA crash log has been saved to %s", r, logFilePath))
		stackText := widget.NewMultiLineEntry()
		stackText.SetText(string(buf[:stackLen]))

		stackScroll := container.NewScroll(stackText)
		stackScroll.SetMinSize(fyne.NewSize(0, 300))

		win.CenterOnScreen()

		win.SetContent(container.NewVBox(
			errText,
			stackScroll,
		))

		win.SetCloseIntercept(func() {
			os.Exit(1)
		})

		win.Show()
	}
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

func GetPlayerInfos(uuid string) (*PlayerInfo, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.mihomo.me/sr_info_parsed/%s?lang=fr", uuid))
	if err != nil {
		return nil, fmt.Errorf("error getting player infos: %w", err)
	}
	defer resp.Body.Close()

	var playerInfo PlayerInfo
	err = json.NewDecoder(resp.Body).Decode(&playerInfo)
	if err != nil {
		logger.Logger.Error(err)
		return nil, fmt.Errorf("error decoding player infos: %w", err)
	}

	return &playerInfo, nil
}
