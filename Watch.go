package main

import (
	"errors"
	"fmt"
	"ipsd_vsc/Utils"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watch struct {
	watch *fsnotify.Watcher
}

func FindInList(str string, strList []string) bool {
	if str == "" {
		return false
	}

	if strList == nil {
		return false
	}

	if len(strList) == 0 {
		return false
	}

	for _, s := range strList {
		if s == str {
			return true
		}
	}

	return false
}

//监控目录
func (w *Watch) watchDir(dir string) {
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = w.watch.Add(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	go func() {
		var updateFileList []string
		for {
			select {
			case ev := <-w.watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						fmt.Println("Copy/Generate ", ev.Name)
						//这里获取新创建文件的信息，如果是目录，则加入监控中
						if FindInList(ev.Name, updateFileList) == false {
							updateFileList = append(updateFileList, ev.Name)
						}
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							w.watch.Add(ev.Name)
						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						if FindInList(ev.Name, updateFileList) == false {
							fmt.Println("Updating ", ev.Name)
							updateFileList = append(updateFileList, ev.Name)
						}
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						fmt.Println("Deleting ", ev.Name)
					}
				}
			case err := <-w.watch.Errors:
				{
					fmt.Println("error : ", err)
					return
				}
			}
		}
	}()
}
func (w *Watch) removeWatch(dir string) {
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = w.watch.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func GetOutputFolderPath(siteFolderPath string) (string, error) {
	outputFolderPath := filepath.Join(siteFolderPath, "Output")

	if Utils.PathIsExist(outputFolderPath) == false {
		return "", errors.New("RunMonitor.MonitorCompile.GetOutputFolderPath: output folder path not exist " + outputFolderPath)
	}

	return filepath.EvalSymlinks(outputFolderPath)
}

func MonitorCompile(siteFolderPath string) *Watch {

	outputDir, errGetOutputFolderPath := GetOutputFolderPath(siteFolderPath)

	if errGetOutputFolderPath != nil {
		var errMsg = "RunMonitor: Cannot find output folder " + outputDir
		Utils.Logger.Println(errMsg)
		return nil
	}

	if Utils.PathIsExist(outputDir) == false {
		var errMsg = "RunMonitor: Output Folder not exist " + outputDir
		Utils.Logger.Println(errMsg)
		return nil
	}

	watch, _ := fsnotify.NewWatcher()
	w := Watch{
		watch: watch,
	}
	w.watchDir(outputDir)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	return &w
}

func StopMonitor(w *Watch, siteFolderPath string) {
	outputDir, errGetOutputFolderPath := GetOutputFolderPath(siteFolderPath)

	if errGetOutputFolderPath != nil {
		var errMsg = "RunMonitor: Cannot find output folder " + outputDir
		Utils.Logger.Println(errMsg)
		return
	}

	if Utils.PathIsExist(outputDir) == false {
		var errMsg = "RunMonitor: Output Folder not exist " + outputDir
		Utils.Logger.Println(errMsg)
		return
	}
	w.removeWatch(outputDir)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}
