package Configuration

import (
	"errors"
	"ipsd/Utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

func GetQuickHelpPath() (string, error) {
	executionPath, errPath := GetCurrentPath()

	if errPath != nil {
		return "", errors.New("Cannot get path of current executable")
	}

	var helpFilePath = filepath.Join(executionPath, "ipsdQuickHelp.txt")
	if Utils.PathIsExist(helpFilePath) == false {
		return "", errors.New("Cannot find QuickHelp.txt at path " + helpFilePath)
	}
	return helpFilePath, nil
}
