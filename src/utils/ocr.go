package utils

import (
	"errors"
	"image"
	"image/color"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/X3ne/hsrpc/src/consts"
	"github.com/X3ne/hsrpc/src/logger"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
)

type OCRConfig struct {
	ExecutablePath      *string
	PreprocessThreshold *int
}

type OCRManager struct {
	config OCRConfig
	HWND   win.HWND
}

var OcrManager *OCRManager

func InitOcr(cfg OCRConfig, hWnd win.HWND) {
	OcrManager = &OCRManager{
		config: cfg,
		HWND:   hWnd,
	}
}

func preprocessImage(img image.Image, threshold int) image.Image {
	grayImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			grayImg.Set(x, y, grayColor)
		}
	}

	binarizedImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			if grayImg.GrayAt(x, y).Y > uint8(threshold) {
				binarizedImg.SetGray(x, y, color.Gray{Y: 255})
			} else {
				binarizedImg.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	return binarizedImg
}

func (m *OCRManager) StartOcr(path string) (string, error) {
	executablePath := *m.config.ExecutablePath
	if executablePath == "" {
		return "", errors.New("tesseract executable path is not set")
	}

	whitelistChars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz, ():"

	cmd := exec.Command(executablePath, path, "stdout", "-l", "eng", "-c", "tessedit_char_whitelist="+whitelistChars)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	body, err := cmd.Output()
	if err != nil {
		return "", err
	}

	outputText := strings.ReplaceAll(string(body), "\n", " ")
	outputText = strings.ReplaceAll(outputText, "\r", " ")
	outputText = strings.ReplaceAll(outputText, "\t", " ")

	return outputText, nil
}

func (m *OCRManager) WindowOcr(rect Rect, job string, preprocess bool) (string, image.Image) {
	var winRect win.RECT
	win.GetWindowRect(m.HWND, &winRect)
	image := robotgo.CaptureImg(int(winRect.Left)+rect.X, int(winRect.Top)+rect.Y, rect.Width, rect.Height)
	if image == nil {
		return "", nil
	}

	if preprocess {
		image = preprocessImage(image, *m.config.PreprocessThreshold)
	}

	appPath, err := GetAppPath()
	if err != nil {
		logger.Logger.Errorf("["+job+"] "+"error: %s", err)
		return "", nil
	}

	imagePath := filepath.Join(appPath, consts.TmpDir, job+".png")
	err = CreateDir(imagePath)
	if err != nil {
		logger.Logger.Errorf("["+job+"] "+"error: %s", err)
		return "", nil
	}

	SaveImg(image, imagePath)

	text, err := m.StartOcr(imagePath)
	if err != nil {
		logger.Logger.Errorf("["+job+"] "+"error: %s", err)
		return "", nil
	}

	text = strings.TrimSpace(text)

	return text, image
}
