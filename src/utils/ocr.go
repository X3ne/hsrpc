package utils

import (
	"bytes"
	"image"
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
	ExecutablePath	string
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

func (m *OCRManager) StartOcr(imageBytes []byte) (string, error) {
	cmd := exec.Command(m.config.ExecutablePath, "-c", "tessedit_char_whitelist=0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz, ", "stdin", "stdout")
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

	stdinPipe.Close()

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

func (m *OCRManager) WindowOcr(rect Rect, job string) (string, image.Image) {
	image := robotgo.CaptureImg(rect.X, rect.Y, rect.Width, rect.Height)
	if image == nil {
		return "", nil
	}

	grayImg := ConvertToGrayscale(image)

	SaveImg(grayImg, "tmp/" + job + ".png") // Save the image for debugging purposes

	buf := new(bytes.Buffer)
	err := png.Encode(buf, grayImg)
	if err != nil {
		logger.Logger.Errorf("["+job+"] "+"Error: %s", err)
	}
	text, err := m.StartOcr(buf.Bytes())
	if err != nil {
		logger.Logger.Errorf("["+job+"] "+"Error: %s", err)
	}

	return text, grayImg
}
