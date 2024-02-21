package gui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/X3ne/hsrpc/src/config"
	"github.com/X3ne/hsrpc/src/consts"
	"github.com/X3ne/hsrpc/src/internal/bundle"
	"github.com/X3ne/hsrpc/src/logger"
	"github.com/X3ne/hsrpc/src/utils"
	"github.com/creativeprojects/go-selfupdate"
)

func debounce(callback func(string), duration time.Duration) func(string) {
	var timer *time.Timer
	return func(s string) {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(duration, func() {
			callback(s)
		})
	}
}

func createRectForm(rect *utils.Rect, onChanged func()) *widget.Form {
	xEntry := widget.NewEntry()
	yEntry := widget.NewEntry()
	widthEntry := widget.NewEntry()
	heightEntry := widget.NewEntry()

	xEntry.SetText(strconv.Itoa(int(rect.X)))
	yEntry.SetText(strconv.Itoa(int(rect.Y)))
	widthEntry.SetText(strconv.Itoa(int(rect.Width)))
	heightEntry.SetText(strconv.Itoa(int(rect.Height)))

	xEntry.OnChanged = debounce(func(s string) {
		rect.X = stringToInt(s)
		onChanged()
	}, 500*time.Millisecond)


	yEntry.OnChanged = debounce(func(s string) {
		rect.Y = stringToInt(s)
		onChanged()
	}, 500*time.Millisecond)

	widthEntry.OnChanged = debounce(func(s string) {
		rect.Width = stringToInt(s)
		onChanged()
	}, 500*time.Millisecond)

	heightEntry.OnChanged = debounce(func(s string) {
		rect.Height = stringToInt(s)
		onChanged()
	}, 500*time.Millisecond)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "X", Widget: xEntry},
			{Text: "Y", Widget: yEntry},
			{Text: "Width", Widget: widthEntry},
			{Text: "Height", Widget: heightEntry},
		},
	}

	return form
}

func (g *GUI) saveConfig() {
	if err := config.SaveConfig(g.RPCApp.Config); err != nil {
		logger.Logger.Error(err)
	}
}

func getPresenceString(isGatewayConnected bool) string {
	if isGatewayConnected {
		return "Connected to Discord Gateway"
	}
	return "Not connected to Discord Gateway"
}

func stringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Logger.Error(err)
	}
	return i
}

func ImportIcon() []byte {
	return bundle.Get("icon.png")
}

func Restart() {
	executable, err := os.Executable()
	if err != nil {
		logger.Logger.Fatalf("failed to get executable path: %v", err)
	}

	cmd := exec.Command(executable)

	cmd.Args = append(cmd.Args, os.Args[1:]...)

	cmd.Env = os.Environ()

	cmd.Dir, err = os.Getwd()
	if err != nil {
		logger.Logger.Fatalf("failed to get working directory: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		logger.Logger.Fatalf("failed to restart: %v", err)
	}

	os.Exit(0)
}

func HandleUpdateError(message string, err error, label *widget.Label, updateCompleted chan bool, window fyne.Window) {
	if err != nil {
		logger.Logger.Errorf(message, err)
	} else {
		logger.Logger.Errorf(message)
	}
	label.SetText(message)
	window.SetFixedSize(false)
	updateCompleted <- false
}

func GetLatestUpdate() (*selfupdate.Release, error) {
	if consts.Version == "" {
		return nil, fmt.Errorf("version is not set")
	}

	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("X3ne/hsrpc"))
	if err != nil {
		return nil, fmt.Errorf("error occurred while checking for updates: %s", err)
	}
	if !found {
		return nil, fmt.Errorf("no update found")
	}

	if latest.LessOrEqual(consts.Version) {
		logger.Logger.Infof("current version (%s) is the latest", latest)
		return nil, nil
	}

	logger.Logger.Infof("update available: %s", latest.Version)
	return latest, nil
}
