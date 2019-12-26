package main

import (
	"flag"
	"fmt"
	"strings"
)

type CommandParser struct {
	CurrentCommand    string
	SiteTitle         string
	SiteFolderPath    string
	MonitorFolderPath string
	IndexPageSize     string
	//MonitorInterval   int64
}

func (cpp *CommandParser) ParseCommand() bool {
	//Set All Arguments
	flag.StringVar(&cpp.CurrentCommand, "Command", "", "Command you want to run")
	flag.StringVar(&cpp.SiteTitle, "SiteTitle", "", "Title of the site to update")
	flag.StringVar(&cpp.SiteFolderPath, "SiteFolder", "", "Path of Site Project folder to update")
	flag.StringVar(&cpp.IndexPageSize, "IndexPageSize", "Normal", "The size of index and more page, normal page contains 20 items,small page contains 10 items, very small page contains 5 items, big page contains 30 items (Default Normal)")
	//flag.Int64Var(&cpp.MonitorInterval, "MonitorInterval", 600, "Monitor Internal, second as unit, defualt 600 seconds")
	flag.StringVar(&cpp.MonitorFolderPath, "MonitorFolder", "", "Folder to monitor, changes in this folder will be monitored and updated to target site project")

	//Parse
	flag.Parse()

	//Trim all String properties
	cpp.CurrentCommand = strings.TrimSpace(cpp.CurrentCommand)
	cpp.SiteFolderPath = strings.TrimSpace(cpp.SiteFolderPath)
	cpp.SiteTitle = strings.TrimSpace(cpp.SiteTitle)
	cpp.IndexPageSize = strings.TrimSpace(cpp.IndexPageSize)
	cpp.MonitorFolderPath = strings.TrimSpace(cpp.MonitorFolderPath)

	//To Upper
	cpp.CurrentCommand = strings.ToUpper(cpp.CurrentCommand)

	var ret bool
	ret = true
	cpp.CurrentCommand = strings.ToUpper(cpp.CurrentCommand)
	//Check Properties of New Site Project
	switch cpp.CurrentCommand {
	case COMMAND_NEWMONITOR:
		if cpp.SiteFolderPath == "" {
			fmt.Println("Site Folder is empty")
			return false
		}
		if cpp.SiteTitle == "" {
			fmt.Println("SiteTitle is empty, cannot create site ")
			ret = false
		}

		if cpp.MonitorFolderPath == "" {
			fmt.Println("Monitor Folder is empty, will not create site")
			ret = false
		}

	case COMMAND_RUNMONITOR:
		if cpp.MonitorFolderPath == "" {
			fmt.Println("Monitor Folder is empty, will not create site")
			ret = false
		}
		/*
			if cpp.MonitorInterval <= 0 {
				fmt.Println("Monitor Internal should bigger than 0")
				ret = false
			}
		*/
	case COMMAND_LISTNORMALFILE:
		if cpp.MonitorFolderPath == "" {
			fmt.Println("Monitor Folder is empty, will not create site")
			ret = false
		}
	default:
	}
	return ret
}
