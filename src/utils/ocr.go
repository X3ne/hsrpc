package utils

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"github.com/X3ne/hsrpc/src/logger"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
)

type OCRConfig struct {
	ExecutablePath	*string
}

type OCRManager struct {
	config	OCRConfig
}

var OcrManager *OCRManager

func InitOcr(cfg OCRConfig, hWnd win.HWND) {
	OcrManager = &OCRManager{
		config: cfg,
	}
}

func preprocessImage(img image.Image) image.Image {
	grayImg := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			grayImg.Set(x, y, grayColor)
		}
	}

	threshold := 180
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

func (m *OCRManager) StartOcr(imageBytes []byte) (string, error) {
	executablePath := *m.config.ExecutablePath
	if executablePath == "" {
		return "", errors.New("tesseract executable path is not set")
	}

	whitelistChars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz, "

	cmdArgs := []string{
		"-l", "eng",
		"-c", "tessedit_char_whitelist=" + whitelistChars,
		"stdin", "stdout",
	}

	cmd := exec.Command(executablePath, cmdArgs...)

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}
	defer stdinPipe.Close()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(stdinPipe, bytes.NewReader(imageBytes))
	if err != nil {
		return "", err
	}

	err = stdinPipe.Close()
	if err != nil {
		return "", err
	}

	var outputBuf bytes.Buffer
	_, err = io.Copy(&outputBuf, stdoutPipe)
	if err != nil {
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	outputText := strings.ReplaceAll(outputBuf.String(), "\n", " ")
	outputText = strings.ReplaceAll(outputText, "\r", " ")
	outputText = regexp.MustCompile(`\s+`).ReplaceAllString(outputText, " ")

	return outputText, nil
}

func (m *OCRManager) WindowOcr(rect Rect, job string, preprocess bool) (string, image.Image) {
	image := robotgo.CaptureImg(rect.X, rect.Y, rect.Width, rect.Height)
	if image == nil {
		return "", nil
	}

	if preprocess {
		image = preprocessImage(image)
	}

	// grayImg := ConvertToGrayscale(image)

	SaveImg(image, "tmp/" + job + ".png") // Save the image for debugging purposes

	buf := new(bytes.Buffer)
	err := png.Encode(buf, image)
	if err != nil {
		logger.Logger.Errorf("["+job+"] "+"Error: %s", err)
		return "", nil
	}
	text, err := m.StartOcr(buf.Bytes())
	if err != nil {
		logger.Logger.Errorf("["+job+"] "+"Error: %s", err)
		return "", nil
	}

	text = strings.TrimSpace(text)

	return text, image
}
