// IPSCM project main.go
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"ipsd/Configuration"
	"ipsd/Monitor"
	"ipsd/Utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Try2FindSpFile(siteFolderPath string) (string, error) {
	if Utils.PathIsExist(siteFolderPath) == false {
		var errMsg = "Try2FindSpFile: Site Folder not exist"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	var spCount int
	spCount = 0
	var spFileName string
	spFileName = ""

	files, _ := ioutil.ReadDir(siteFolderPath)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sp") {
			spFileName = f.Name()
			spCount++
			if spCount > 1 {
				var errMsg = "Try2FindSpFile: More than 1 .sp file"
				Utils.Logger.Println(errMsg)
				return "", errors.New(errMsg)
			}
		}
	}
	return spFileName, nil
}

func Try2FindSmFile(monitorFolderPath string) (string, error) {
	if Utils.PathIsExist(monitorFolderPath) == false {
		return "", errors.New("Try2FindSmFile: Monitor Folder not exist")
	}

	var smCount int
	smCount = 0
	var smFileName string
	smFileName = ""

	files, _ := ioutil.ReadDir(monitorFolderPath)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sm") {
			smFileName = f.Name()
			smCount++
			if smCount > 1 {
				return "", errors.New("Try2FindSmFile: More than 1 .sm file")
			}
		}
	}

	if smFileName != "" {

		return filepath.Join(monitorFolderPath, smFileName), nil
	}
	return "", errors.New("There is no sm file in folder " + monitorFolderPath)
}

func Dispatch(cp CommandParser) (bool, error) {

	switch cp.CurrentCommand {
	case COMMAND_NEWMONITOR:
		return NewMonitor(cp.MonitorFolderPath, cp.SiteFolderPath, cp.SiteTitle)

	case COMMAND_RUNMONITOR:
		//MonitorSite(cp.MonitorFolderPath, cp.IndexPageSize, cp.MonitorInterval)
		RunMonitor(cp.MonitorFolderPath, cp.IndexPageSize)
		return true, nil
	case COMMAND_LISTNORMALFILE:
		ListNormalFiles(cp.MonitorFolderPath)
	default:
		DisplayHelp()
		Utils.Logger.Println("Command not found " + cp.CurrentCommand)
		return false, errors.New("Main.Command not found " + cp.CurrentCommand)
	}
	return true, nil
}

func DisplayHelp() {
	helpContent, errHelp := GetQuickHelpInformation()
	if errHelp != nil {
		Utils.Logger.Println("Main.DisplayHelp: Cannot get quick help information")
	} else {
		fmt.Println("")
		fmt.Println("QuickHelp")
		fmt.Println("==========")
		fmt.Println(helpContent)
	}
}

func GetQuickHelpInformation() (string, error) {
	quickHelpFilePath, errPath := Configuration.GetQuickHelpPath()
	if errPath != nil {
		Utils.Logger.Println("GetQuickHelp: " + errPath.Error())
		return "", errPath
	}
	bHelpContent, errRead := ioutil.ReadFile(quickHelpFilePath)

	if errRead != nil {
		Utils.Logger.Println("GetQuickHelp: " + errRead.Error())
		return "", errRead
	}
	var sHelpContent = string(bHelpContent)
	return sHelpContent, nil
}

func ListNormalFiles(monitorFolderPath string) {
	if monitorFolderPath == "" {
		Utils.Logger.Println("RunMonitor: Monitor Folder Path is empty")
		return
	}

	if Utils.PathIsExist(monitorFolderPath) == false {
		Utils.Logger.Println("RunMonitor: Monitor Folder doesn't exist")
		return
	}

	monitorDefinitionFilePath, _ := Try2FindSmFile(monitorFolderPath)
	if monitorDefinitionFilePath == "" {
		Utils.Logger.Println("RunMonitor: No monitor definition file found in folder  " + monitorFolderPath)
		return
	}

	//Load Monitor
	var smp = Monitor.NewSiteMonitor()

	_, errLoad := smp.LoadFromFile(monitorDefinitionFilePath)

	if errLoad != nil {
		Utils.Logger.Println("RunMonitor: Cannot Load Monitor from file " + monitorDefinitionFilePath)
		return
	}

	Monitor.IPSC_ListFile(smp.SiteFolderPath, smp.SiteTitle)

}

func NewMonitor(monitorFolderPath, siteFolderPath, siteTitle string) (bool, error) {
	//Check variables
	if monitorFolderPath == "" {
		return false, errors.New("New Monitor: Monitor Folder Path is empty")
	}

	if siteFolderPath == "" {
		return false, errors.New("New Monitor: Target Site Folder Path is empty")
	}

	if siteTitle == "" {
		return false, errors.New("New Monitor: Target Site Title is empty")
	}

	if Utils.PathIsExist(monitorFolderPath) == false {
		return false, errors.New("New Monitor: Monitor Folder doesn't exist")
	}

	if Utils.PathIsExist(siteFolderPath) == false {
		return false, errors.New("New Monitor: Target Site Folder doesn't exist")
	}

	//Initialize Site Monitor definition file
	monitorDefinitionFilePath, _ := Try2FindSmFile(monitorFolderPath)
	if monitorDefinitionFilePath != "" {
		return false, errors.New("New Monitor: There is already a monitor in folder " + monitorFolderPath)
	}

	var smp = Monitor.NewSiteMonitor_WithArgs(monitorFolderPath, siteFolderPath, siteTitle, Utils.CurrentTime())

	//Export files to monitor folder
	var monitorSiteFilePath = filepath.Join(siteFolderPath, siteTitle+".sp")

	if Utils.PathIsExist(monitorSiteFilePath) == false {
		monitorSiteFilePath, _ = Try2FindSpFile(siteFolderPath)
	}

	markdownPages, htmlPages, linkPages, bExportSite, errExportSite := Monitor.IPSC_ExportSite(siteFolderPath, siteTitle, monitorFolderPath)
	var markdownMonitorFolder = filepath.Join(monitorFolderPath, "Markdown")
	var htmlMonitorFolder = filepath.Join(monitorFolderPath, "Html")
	var linkMonitorFolder = filepath.Join(monitorFolderPath, "Link")
	var fileMonitorFolder = filepath.Join(monitorFolderPath, "Files")

	if bExportSite && errExportSite == nil {
		for _, markdownPage := range markdownPages {
			var mdPage Monitor.MarkdownPage
			mdPage.ID = markdownPage.ID
			mdPage.LastModified = markdownPage.LastModified
			var fName = filepath.Base(markdownPage.FilePath)
			mdPage.FilePath = filepath.Join(markdownMonitorFolder, fName)
			smp.AddMarkdown(mdPage)
		}

		for _, htmlPage := range htmlPages {
			var htmPage Monitor.HtmlPage
			htmPage.ID = htmlPage.ID
			htmPage.LastModified = htmlPage.LastModified
			var fName = filepath.Base(htmlPage.FilePath)
			htmPage.FilePath = filepath.Join(htmlMonitorFolder, fName)
			smp.AddHtml(htmPage)
		}

		for _, linkPage := range linkPages {
			smp.AddLink(linkPage)
		}

		if Utils.PathIsExist(fileMonitorFolder) {
			smp.NormalFiles = smp.GetNormalFileList()
		}

	}

	if Utils.PathIsExist(markdownMonitorFolder) == false {
		Utils.MakeFolder(markdownMonitorFolder)
	}

	if Utils.PathIsExist(htmlMonitorFolder) == false {
		Utils.MakeFolder(htmlMonitorFolder)
	}

	if Utils.PathIsExist(linkMonitorFolder) == false {
		Utils.MakeFolder(linkMonitorFolder)
	}

	if Utils.PathIsExist(fileMonitorFolder) == false {
		Utils.MakeFolder(fileMonitorFolder)
	}

	var linkMonitorFile = filepath.Join(linkMonitorFolder, "Link.liks")
	if Utils.PathIsExist(linkMonitorFile) == false {
		Utils.CreateFile(linkMonitorFile)
	}

	//Copy Template file,for now News and Blank
	var templateFolderPath = filepath.Join(monitorFolderPath, "Templates")
	if Utils.PathIsExist(templateFolderPath) == false {
		Utils.MakeFolder(templateFolderPath)
	}

	var newsMarkdownTemplateFilePath = filepath.Join(templateFolderPath, "News.md")
	_, errCreateNews := Monitor.IPSC_CreateMarkdown(siteFolderPath, siteTitle, newsMarkdownTemplateFilePath, "News")

	if errCreateNews != nil {
		fmt.Println("Cannot copy News Tempalte file (News.md), you can copy it from Resources folder of ipsc")
	}

	var blankMarkdownTemplateFilePath = filepath.Join(templateFolderPath, "Blank.md")
	_, errCreateBlank := Monitor.IPSC_CreateMarkdown(siteFolderPath, siteTitle, blankMarkdownTemplateFilePath, "Blank")

	if errCreateBlank != nil {
		fmt.Println("Cannot copy Blank Tempalte file (Blank.md), you can copy it from Resources folder of ipsc")
	}

	//Save Monitor File
	monitorDefinitionFilePath = filepath.Join(monitorFolderPath, "monitor.sm")

	bSave, errSave := smp.SaveToFile(monitorDefinitionFilePath)
	if errSave != nil {
		return bSave, errSave
	}

	fmt.Println("Create New Monitor Success, Monitor Folder Path " + monitorFolderPath)

	return true, nil
}

func MonitorSite(monitorFolderPath, indexPageSize string, monitorInterval int64) {
	var ch chan int
	//定时任务
	if monitorInterval < 0 {
		return
	}
	fmt.Println(Utils.CurrentTime())
	fmt.Println("Start Monitor, Wait for " + strconv.FormatInt(monitorInterval, 10) + " seconds to start the first monitor")
	fmt.Println()
	ticker := time.NewTicker(time.Second * time.Duration(monitorInterval))
	go func() {
		for range ticker.C {
			_, errMonitor := RunMonitor(monitorFolderPath, indexPageSize)
			if errMonitor != nil {
				Utils.Logger.Println(errMonitor.Error())
			}
		}
		ch <- 1
	}()
	<-ch
}

func CheckIndexPageSize(indexPageSize string) bool {
	if indexPageSize == "" {
		return false
	}

	indexPageSize = strings.ToUpper(indexPageSize)

	if indexPageSize == "NORMAL" || indexPageSize == "SMALL" || indexPageSize == "VERYSMALL" || indexPageSize == "BIG" {
		return true
	}

	return false
}

func RunMonitor(monitorFolderPath, indexPageSize string) (bool, error) {

	fmt.Println("===============")
	fmt.Println(Utils.CurrentTime())
	fmt.Println("Checking and Update monitor folder " + monitorFolderPath)

	if monitorFolderPath == "" {
		return false, errors.New("RunMonitor: Monitor Folder Path is empty")
	}

	if indexPageSize == "" {
		return false, errors.New("RunMonitor: Index Page Size is empty")
	}

	if Utils.PathIsExist(monitorFolderPath) == false {
		return false, errors.New("RunMonitor: Monitor Folder doesn't exist")
	}

	if CheckIndexPageSize(indexPageSize) == false {
		return false, errors.New("RunMonitor: Index Page Size error, should be Normal,Small,VerySmall or Big")
	}

	monitorDefinitionFilePath, _ := Try2FindSmFile(monitorFolderPath)
	if monitorDefinitionFilePath == "" {
		return false, errors.New("RunMonitor: No monitor definition file found in folder  " + monitorFolderPath)
	}

	//Load Monitor
	var smp = Monitor.NewSiteMonitor()

	bLoad, errLoad := smp.LoadFromFile(monitorDefinitionFilePath)

	if errLoad != nil {
		return bLoad, errors.New("RunMonitor: Cannot Load Monitor from file " + monitorDefinitionFilePath)
	}

	//Check files
	var markdownFolderPath = filepath.Join(monitorFolderPath, "Markdown")
	var htmlFolderPath = filepath.Join(monitorFolderPath, "Html")

	var linkFilePath = filepath.Join(monitorFolderPath, "Link", "Link.liks")

	var linkFolderPath = filepath.Join(monitorFolderPath, "Link")

	var normalFileFolderPath = filepath.Join(monitorFolderPath, "Files")

	var mdChanged, htmChanged, linkChanged, fileChanged bool

	if Utils.PathIsExist(markdownFolderPath) {
		fmt.Println("Checking Markdown")
		var addMd, updateMd, deleteMd int
		addMd = 0
		updateMd = 0
		deleteMd = 0
		files, _ := ioutil.ReadDir(markdownFolderPath)

		for _, f := range files {
			extension := filepath.Ext(f.Name())
			if strings.HasPrefix(f.Name(), "_") {
				fmt.Println("Name of " + f.Name() + " start with _ , and it will be treated as temp file to be ignored. If you finished it and want to publish it ,remove _ from beginning of the file name")
				continue
			}
			if extension == ".md" || extension == ".markdown" || extension == ".mmd" || extension == ".mdown" {
				//Markdown
				var fPath = filepath.Join(markdownFolderPath, f.Name())
				var fLastModified = f.ModTime().Format("2006-01-02 15:04:05")
				index := smp.GetMarkdown(fPath)

				if index == -1 {
					fmt.Println(fPath + " is a new file, will add it to ipsc")
					addMd = addMd + 1
					//Add
					mdProperties, _, errReadProperteis := Monitor.ReadMarkdownPageProperties(fPath)

					if errReadProperteis != nil {
						var errMsg = "RunMonitor: Cannot read Markdown properties"
						Utils.Logger.Println(errMsg)
						return false, errReadProperteis
					}

					titleImagePath, errImageName := Utils.GetImageWithSameName(fPath)

					if errImageName != nil {
						Utils.Logger.Println("No Title Image with same name of " + fPath + " found")
						Utils.Logger.Println("Title Image of " + fPath + " will be empty")
					}

					//Check Title Image Size, should smaller than 30KB
					bImageSize, errImageSize := Utils.ImageTooBig(titleImagePath)

					if errImageSize != nil {
						Utils.Logger.Println("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
						Utils.Logger.Println(titleImagePath)
						return false, errors.New("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
					}

					if bImageSize == true {
						var errMsg = "RunMonitor: TitleImage " + titleImagePath + " is bigger than 30KB, please edit it firstly to make it smaller than 30KB, then add again"
						Utils.Logger.Println(errMsg)
						return false, errors.New(errMsg)
					}

					newID, bAdd, errAdd := Monitor.IPSC_AddMarkdown(smp.SiteFolderPath, smp.SiteTitle, fPath, mdProperties.Title, mdProperties.Author, titleImagePath, mdProperties.IsTop)

					if errAdd != nil || bAdd == false {
						var errMsg = "RunMonitor: Cannot Add Markdown file " + fPath
						Utils.Logger.Println(errMsg)
						Utils.Logger.Println(errAdd.Error())
						return false, errors.New(errMsg)
					}

					//Add to monitor file
					var newMdPage Monitor.MarkdownPage
					newMdPage.FilePath = fPath
					newMdPage.ID = newID
					newMdPage.LastModified = Utils.CurrentTime()

					smp.AddMarkdown(newMdPage)
					smp.SaveToFile(monitorDefinitionFilePath)

				} else {
					//Update
					var sourceMarkdown = smp.MarkdownFiles[index]

					//Source file new
					if sourceMarkdown.LastModified < fLastModified {
						fmt.Println(fPath + " has been modified, will update it with ipsc")
						updateMd = updateMd + 1
						mdProperties, _, errReadProperteis := Monitor.ReadMarkdownPageProperties(fPath)

						if errReadProperteis != nil {
							var errMsg = "RunMonitor: Cannot read Markdown properties"
							Utils.Logger.Println(errMsg)
							return false, errReadProperteis
						}

						titleImagePath, errImageName := Utils.GetImageWithSameName(fPath)

						if errImageName != nil {
							Utils.Logger.Println("No Title Image with same name of " + fPath + " found")
							Utils.Logger.Println("Title Image of " + fPath + " will be empty")
						}
						//Check Title Image Size, should smaller than 30KB
						bImageSize, errImageSize := Utils.ImageTooBig(titleImagePath)

						if errImageSize != nil {
							Utils.Logger.Println("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
							Utils.Logger.Println(titleImagePath)
							return false, errors.New("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
						}

						if bImageSize == true {
							var errMsg = "RunMonitor: TitleImage " + titleImagePath + " is bigger than 30KB, please edit it firstly to make it smaller than 30KB, then add again"
							Utils.Logger.Println(errMsg)
							return false, errors.New(errMsg)
						}
						_, errUpdate := Monitor.IPSC_UpdateMarkdownOrHtml(smp.SiteFolderPath, smp.SiteTitle, sourceMarkdown.ID, fPath, mdProperties.Title, mdProperties.Author, titleImagePath, mdProperties.IsTop)

						if errUpdate != nil {
							var errMsg = "RunMonitor: Cannot Update Markdown file " + fPath
							Utils.Logger.Println(errMsg)
							return false, errors.New(errMsg)
						}

						//Update sourceMarkdown
						sourceMarkdown.LastModified = fLastModified
						smp.UpdateMarkdown(sourceMarkdown)
						smp.SaveToFile(monitorDefinitionFilePath)
					}
				}
			}
		}

		//Delete
		var deletedMds []Monitor.MarkdownPage
		for _, sourceMd := range smp.MarkdownFiles {
			if Utils.PathIsExist(sourceMd.FilePath) == false {
				fmt.Println(sourceMd.FilePath + " has been deleted, will delete it from ispc")
				deleteMd = deleteMd + 1
				_, errDelete := Monitor.IPSC_DeletePage(smp.SiteFolderPath, smp.SiteTitle, sourceMd.ID)

				if errDelete != nil {
					var errMsg = "RunMonitor: Cannot Delete file " + sourceMd.FilePath
					Utils.Logger.Println(errMsg)
				}

				imagePath, errImagePath := Utils.GetImageWithSameName(sourceMd.FilePath)

				if errImagePath == nil {
					bDeleteImage := Utils.DeleteFile(imagePath)
					if bDeleteImage == false {
						var errMsg = "RunMonitor: Cannot Delete Image " + imagePath
						Utils.Logger.Println(errMsg)
					}
				}

				//Delete it in smp
				deletedMds = append(deletedMds, sourceMd)
			}
		}

		for _, deletedMd := range deletedMds {
			smp.DeleteMarkdown(deletedMd.FilePath)
			smp.SaveToFile(monitorDefinitionFilePath)
		}

		if addMd == 0 && updateMd == 0 && deleteMd == 0 {
			fmt.Println("Markdown Files not changed, pass")
			mdChanged = false
		} else {
			mdChanged = true
			fmt.Println("Markdown Files")
			fmt.Println("    Add:    " + strconv.Itoa(addMd))
			fmt.Println("    Update: " + strconv.Itoa(updateMd))
			fmt.Println("    Delete: " + strconv.Itoa(deleteMd))
		}
	}

	if Utils.PathIsExist(htmlFolderPath) {
		fmt.Println("Checking Html")
		files, _ := ioutil.ReadDir(htmlFolderPath)
		var addHtm, updateHtm, deleteHtm int
		addHtm = 0
		updateHtm = 0
		deleteHtm = 0

		for _, f := range files {
			extension := filepath.Ext(f.Name())
			if strings.HasPrefix(f.Name(), "_") {
				fmt.Println("Name of " + f.Name() + " start with _ , and it will be treated as temp file to be ignored. If you finished it and want to publish it ,remove _ from beginning of the file name")
				continue
			}
			if extension == ".htm" || extension == ".html" {
				var fPath = filepath.Join(htmlFolderPath, f.Name())
				var fLastModified = f.ModTime().Format("2006-01-02 15:04:05")

				index := smp.GetHtml(fPath)

				if index == -1 {
					//Add
					fmt.Println(fPath + " is a new file, will add it to ipsc")
					addHtm = addHtm + 1
					htmProperties, _, errReadProperteis := Monitor.ReadHtmlProperties(fPath)

					if errReadProperteis != nil {
						var errMsg = "RunMonitor: Cannot read Html properties"
						Utils.Logger.Println(errMsg)
						return false, errReadProperteis
					}

					titleImagePath, errImageName := Utils.GetImageWithSameName(fPath)

					if errImageName != nil {
						Utils.Logger.Println("No Title Image with same name of " + fPath + " found")
						Utils.Logger.Println("Title Image of " + fPath + " will be empty")
					}
					//Check Title Image Size, should smaller than 30KB
					bImageSize, errImageSize := Utils.ImageTooBig(titleImagePath)

					if errImageSize != nil {
						Utils.Logger.Println("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
						Utils.Logger.Println(titleImagePath)
						return false, errors.New("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
					}

					if bImageSize == true {
						var errMsg = "RunMonitor: TitleImage " + titleImagePath + " is bigger than 30KB, please edit it firstly to make it smaller than 30KB, then add again"
						Utils.Logger.Println(errMsg)
						return false, errors.New(errMsg)
					}

					newID, _, errAdd := Monitor.IPSC_AddHtml(smp.SiteFolderPath, smp.SiteTitle, fPath, htmProperties.Title, htmProperties.Author, titleImagePath, htmProperties.IsTop)

					if errAdd != nil {
						var errMsg = "RunMonitor: Cannot Add Html file " + fPath
						Utils.Logger.Println(errMsg)
						Utils.Logger.Println(errAdd.Error())
						return false, errors.New(errMsg)
					}

					var newHtmPage Monitor.HtmlPage
					newHtmPage.FilePath = fPath
					newHtmPage.ID = newID
					newHtmPage.LastModified = Utils.CurrentTime()

					smp.AddHtml(newHtmPage)
					smp.SaveToFile(monitorDefinitionFilePath)
				} else {
					//Update
					var sourceHtml = smp.HtmlFiles[index]

					//Source file new
					if sourceHtml.LastModified < fLastModified {
						fmt.Println(fPath + " has been modified, will update it with ipsc")
						updateHtm = updateHtm + 1
						htmProperties, _, errReadProperteis := Monitor.ReadHtmlProperties(fPath)

						if errReadProperteis != nil {
							var errMsg = "RunMonitor: Cannot read Html properties"
							Utils.Logger.Println(errMsg)
							return false, errReadProperteis
						}

						titleImagePath, errImageName := Utils.GetImageWithSameName(fPath)

						if errImageName != nil {
							Utils.Logger.Println("No Title Image with same name of " + fPath + " found")
							Utils.Logger.Println("Title Image of " + fPath + " will be empty")
						}
						//Check Title Image Size, should smaller than 30KB
						bImageSize, errImageSize := Utils.ImageTooBig(titleImagePath)

						if errImageSize != nil {
							Utils.Logger.Println("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
							Utils.Logger.Println(titleImagePath)
							return false, errors.New("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
						}

						if bImageSize == true {
							var errMsg = "RunMonitor: TitleImage " + titleImagePath + " is bigger than 30KB, please edit it firstly to make it smaller than 30KB, then add again"
							Utils.Logger.Println(errMsg)
							return false, errors.New(errMsg)
						}
						_, errUpdate := Monitor.IPSC_UpdateMarkdownOrHtml(smp.SiteFolderPath, smp.SiteTitle, sourceHtml.ID, fPath, htmProperties.Title, htmProperties.Author, titleImagePath, htmProperties.IsTop)

						if errUpdate != nil {
							var errMsg = "RunMonitor: Cannot Update Html file " + fPath
							Utils.Logger.Println(errMsg)
							return false, errors.New(errMsg)
						}

						sourceHtml.LastModified = fLastModified
						smp.UpdateHtml(sourceHtml)
						smp.SaveToFile(monitorDefinitionFilePath)
					}
				}
			}
		}

		//Delete
		var deletedHtmls []Monitor.HtmlPage
		for _, sourceHtml := range smp.HtmlFiles {
			if Utils.PathIsExist(sourceHtml.FilePath) == false {
				fmt.Println(sourceHtml.FilePath + " has been deleted, will delete it from ispc")
				deleteHtm = deleteHtm + 1
				_, errDelete := Monitor.IPSC_DeletePage(smp.SiteFolderPath, smp.SiteTitle, sourceHtml.ID)

				if errDelete != nil {
					var errMsg = "RunMonitor: Cannot Delete file " + sourceHtml.FilePath
					Utils.Logger.Println(errMsg)
				}

				imagePath, errImagePath := Utils.GetImageWithSameName(sourceHtml.FilePath)

				if errImagePath == nil {
					bDeleteImage := Utils.DeleteFile(imagePath)
					if bDeleteImage == false {
						var errMsg = "RunMonitor: Cannot Delete Image " + imagePath
						Utils.Logger.Println(errMsg)
					}
				}

				deletedHtmls = append(deletedHtmls, sourceHtml)
			}
		}

		for _, deletedHtml := range deletedHtmls {
			smp.DeleteHtml(deletedHtml.FilePath)
			smp.SaveToFile(monitorDefinitionFilePath)
		}

		if addHtm == 0 && updateHtm == 0 && deleteHtm == 0 {
			htmChanged = false
			fmt.Println("Html Files not changed, pass")
		} else {
			htmChanged = true
			fmt.Println("Html Files")
			fmt.Println("    Add:    " + strconv.Itoa(addHtm))
			fmt.Println("    Update: " + strconv.Itoa(updateHtm))
			fmt.Println("    Delete: " + strconv.Itoa(deleteHtm))
		}

	}

	if Utils.PathIsExist(linkFilePath) {
		fmt.Println("Checking Link")
		linkFileInfo, errFileInfo := os.Stat(linkFilePath)
		if errFileInfo == nil {
			var addLink, updateLink, deleteLink int
			addLink = 0
			updateLink = 0
			deleteLink = 0
			fLastModified := linkFileInfo.ModTime().Format("2006-01-02 15:04:05")

			if smp.LinkFileLastModified < fLastModified {
				fLinks, errReadLinks := Monitor.ReadLinksFromFile(linkFilePath)
				if errReadLinks != nil {
					var errMsg = "RunMonitor: Cannot read Links from " + linkFilePath
					Utils.Logger.Println(errMsg)
				} else {

					//Delete
					var deletedLinks []Monitor.LinkPage
					for _, sLink := range smp.LinkFiles {
						if Monitor.FindLink(sLink, fLinks) == false {
							fmt.Println(sLink.Url + " has been deleted, will delete it from ispc")

							_, errDelete := Monitor.IPSC_DeletePage(smp.SiteFolderPath, smp.SiteTitle, sLink.ID)

							if errDelete != nil {
								var errMsg = "RunMonitor: Cannot Delete Link " + sLink.Url
								Utils.Logger.Println(errMsg)
							}

							imagePath, errImagePath := Utils.GetImageWithSameName2(linkFolderPath, sLink.Title)

							if errImagePath == nil {
								bDeleteImage := Utils.DeleteFile(imagePath)
								if bDeleteImage == false {
									var errMsg = "RunMonitor: Cannot Delete Image " + imagePath
									Utils.Logger.Println(errMsg)
								}
							}

							deleteLink = deleteLink + 1
							deletedLinks = append(deletedLinks, sLink)

						}
					}

					for _, deletedLink := range deletedLinks {
						smp.DeleteLink(deletedLink.Url)
						smp.SaveToFile(monitorDefinitionFilePath)
					}

					for _, fLink := range fLinks {
						index := smp.GetLink(fLink.Url)
						if index == -1 {
							fmt.Println(fLink.Url + " is a new Link, will add it to ipsc")

							//Add
							var linkFolderPath = filepath.Join(monitorFolderPath, "Link")
							titleImagePath, errImageName := Utils.GetImageWithSameTitle(linkFolderPath, fLink.Title)

							if errImageName != nil {
								Utils.Logger.Println("No Title Image with same name of " + fLink.Url + " found")
								Utils.Logger.Println("Title Image of " + fLink.Url + " will be empty")
							}
							//Check Title Image Size, should smaller than 30KB
							bImageSize, errImageSize := Utils.ImageTooBig(titleImagePath)

							if errImageSize != nil {
								Utils.Logger.Println("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
								Utils.Logger.Println(titleImagePath)
								return false, errors.New("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
							}

							if bImageSize == true {
								var errMsg = "RunMonitor: TitleImage " + titleImagePath + " is bigger than 30KB, please edit it firstly to make it smaller than 30KB, then add again"
								Utils.Logger.Println(errMsg)
								return false, errors.New(errMsg)
							}
							newID, _, errAddLink := Monitor.IPSC_AddLink(smp.SiteFolderPath, smp.SiteTitle, fLink.Url, fLink.Title, Utils.CurrentUser(), titleImagePath, fLink.IsTop)

							if errAddLink != nil {
								var errMsg = "RunMonitor: Cannot Add Link file " + fLink.Url
								Utils.Logger.Println(errMsg)
								return false, errors.New(errMsg)
							}

							var newLink Monitor.LinkPage
							newLink.Url = fLink.Url
							newLink.ID = newID
							newLink.Title = fLink.Title
							newLink.IsTop = fLink.IsTop

							smp.DeleteLink(fLink.Url)
							smp.AddLink(newLink)
							smp.SaveToFile(monitorDefinitionFilePath)
							addLink = addLink + 1

						} else {
							//Update

							sLink := smp.LinkFiles[index]
							fmt.Println(sLink.Url + " has been modified, will update it with ipsc")

							if sLink.ID == fLink.ID && (sLink.Url != fLink.Url || sLink.Title != fLink.Title || sLink.IsTop != fLink.IsTop) {
								var linkFolderPath = filepath.Join(monitorFolderPath, "Link")
								titleImagePath, errImageName := Utils.GetImageWithSameTitle(linkFolderPath, fLink.Title)

								if errImageName != nil {
									Utils.Logger.Println("No Title Image with same name of " + fLink.Url + " found")
									Utils.Logger.Println("Title Image of " + fLink.Url + " will be empty")
								}

								//Check Title Image Size, should smaller than 30KB
								bImageSize, errImageSize := Utils.ImageTooBig(titleImagePath)

								if errImageSize != nil {
									Utils.Logger.Println("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
									Utils.Logger.Println(titleImagePath)
									return false, errors.New("RunMonitor: Cannot check image size, please make sure the image is smaller than 30KB")
								}

								if bImageSize == true {
									var errMsg = "RunMonitor: TitleImage " + titleImagePath + " is bigger than 30KB, please edit it firstly to make it smaller than 30KB, then add again"
									Utils.Logger.Println(errMsg)
									return false, errors.New(errMsg)
								}

								_, errUpdateLink := Monitor.IPSC_UpdateLink(smp.SiteFolderPath, smp.SiteTitle, fLink.ID, fLink.Url, fLink.Title, Utils.CurrentUser(), titleImagePath, fLink.IsTop)

								if errUpdateLink != nil {
									var errMsg = "RunMonitor: Cannot Update Link file " + fLink.Url
									Utils.Logger.Println(errMsg)
									return false, errors.New(errMsg)
								}

								smp.DeleteLink(sLink.Url)
								smp.AddLink(fLink)
								smp.SaveToFile(monitorDefinitionFilePath)
								updateLink = updateLink + 1
							}
						}
					}

					if addLink != 0 || updateLink != 0 || deleteLink != 0 {
						linkChanged = true
						fmt.Println("Link Files")
						fmt.Println("    Add:    " + strconv.Itoa(addLink))
						fmt.Println("    Update: " + strconv.Itoa(updateLink))
						fmt.Println("    Delete: " + strconv.Itoa(deleteLink))
						//Save Link File
						Monitor.SaveLinksToFile(linkFilePath, smp.LinkFiles)
					} else {
						linkChanged = false
						fmt.Println("Links not changed, pass")
					}

				}
			} else {
				linkChanged = false
				fmt.Println("Links not changed, pass")
			}
		}
	}

	if Utils.PathIsExist(normalFileFolderPath) {
		fmt.Println("Checking Normal File")
		var addNormalFile, updateNormalFile, deleteNormalFile int
		addNormalFile = 0
		updateNormalFile = 0
		deleteNormalFile = 0

		var srcFileList = smp.GetNormalFileList()
		var outputFileList = smp.NormalFiles

		if len(srcFileList) != 0 {
			//Add or Update File
			var handledFolder string
			handledFolder = ""
			for _, srcFile := range srcFileList {

				var iFind = Monitor.GetNormalFile(outputFileList, srcFile.FilePath)
				var srcFullPath = filepath.Join(smp.MonitorFolderPath, srcFile.FilePath)
				if Utils.PathIsDir(srcFullPath) && (strings.HasPrefix(srcFullPath, handledFolder) == false || handledFolder == "") {
					handledFolder = srcFullPath
				}
				if iFind == -1 {
					//New File,add
					if Utils.PathIsDir(srcFullPath) && (strings.HasPrefix(srcFullPath, handledFolder) == false || handledFolder == "") {
						_, errAdd := Monitor.IPSC_AddFile(smp.SiteFolderPath, smp.SiteTitle, srcFullPath)

						if errAdd != nil {
							var errMsg = "RunMonitor: Cannot Add Normal file " + srcFullPath
							Utils.Logger.Println(errMsg)
							Utils.Logger.Println(errAdd.Error())
							return false, errors.New(errMsg)
						}

						smp.NormalFiles = Monitor.AddNormalFile(smp.NormalFiles, srcFile)
						smp.SaveToFile(monitorDefinitionFilePath)
					} else if Utils.PathIsDir(srcFullPath) && strings.HasPrefix(srcFullPath, handledFolder) {
						_, errAdd := Monitor.IPSC_AddFile(smp.SiteFolderPath, smp.SiteTitle, handledFolder)

						if errAdd != nil {
							var errMsg = "RunMonitor: Cannot Add Normal file " + srcFullPath
							Utils.Logger.Println(errMsg)
							Utils.Logger.Println(errAdd.Error())
							return false, errors.New(errMsg)
						}

						smp.NormalFiles = Monitor.AddNormalFile(smp.NormalFiles, srcFile)
						smp.SaveToFile(monitorDefinitionFilePath)
					} else if Utils.PathIsFile(srcFullPath) && strings.HasPrefix(srcFullPath, handledFolder) {
						addNormalFile = addNormalFile + 1
						smp.NormalFiles = Monitor.AddNormalFile(smp.NormalFiles, srcFile)
						smp.SaveToFile(monitorDefinitionFilePath)
						continue
					} else if Utils.PathIsFile(srcFullPath) {
						_, errAdd := Monitor.IPSC_AddFile(smp.SiteFolderPath, smp.SiteTitle, srcFullPath)

						if errAdd != nil {
							var errMsg = "RunMonitor: Cannot Add Normal file " + srcFullPath
							Utils.Logger.Println(errMsg)
							Utils.Logger.Println(errAdd.Error())
							return false, errors.New(errMsg)
						}

						smp.NormalFiles = Monitor.AddNormalFile(smp.NormalFiles, srcFile)
						smp.SaveToFile(monitorDefinitionFilePath)
						addNormalFile = addNormalFile + 1
					}
				} else {
					//Update
					var dstFile = outputFileList[iFind]

					if srcFile.LastModified > dstFile.LastModified {
						var srcFullPath = filepath.Join(smp.MonitorFolderPath, srcFile.FilePath)
						if Utils.PathIsDir(srcFullPath) {
							continue
						}
						_, errDelete := Monitor.IPSC_DeleteFile(smp.SiteFolderPath, smp.SiteTitle, srcFile.FilePath)

						if errDelete != nil {
							var errMsg = "RunMonitor: Cannot Delete Normal file (Update File) " + srcFile.FilePath
							Utils.Logger.Println(errMsg)
							Utils.Logger.Println(errDelete.Error())
							return false, errors.New(errMsg)
						}

						_, errAdd := Monitor.IPSC_AddFile(smp.SiteFolderPath, smp.SiteTitle, srcFullPath)

						if errAdd != nil {
							var errMsg = "RunMonitor: Cannot Add Normal file " + srcFullPath
							Utils.Logger.Println(errMsg)
							Utils.Logger.Println(errAdd.Error())
							return false, errors.New(errMsg)
						}

						smp.NormalFiles = Monitor.UpdateNormalFile(smp.NormalFiles, srcFile)
						smp.SaveToFile(monitorDefinitionFilePath)
						updateNormalFile = updateNormalFile + 1
					}
				}
			}

			//Delete

			var deletedNormalFiles []Monitor.NormalFile
			for _, dstFile := range outputFileList {
				var iFind = Monitor.GetNormalFile(srcFileList, dstFile.FilePath)
				if iFind == -1 {

					deletedNormalFiles = append(deletedNormalFiles, dstFile)

					_, errDelete := Monitor.IPSC_DeleteFile(smp.SiteFolderPath, smp.SiteTitle, "."+dstFile.FilePath)

					if errDelete != nil {
						var errMsg = "RunMonitor: Cannot Delete Normal file (Update File) " + dstFile.FilePath
						Utils.Logger.Println(errMsg)
						Utils.Logger.Println(errDelete.Error())
						return false, errors.New(errMsg)
					}

					deleteNormalFile = deleteNormalFile + 1
				}
			}

			for _, deletedNormalFile := range deletedNormalFiles {
				smp.NormalFiles = Monitor.DeleteNormalFile(smp.NormalFiles, deletedNormalFile)
				smp.SaveToFile(monitorDefinitionFilePath)
			}

			if addNormalFile == 0 && updateNormalFile == 0 && deleteNormalFile == 0 {
				fileChanged = false
				fmt.Println("Normal Files not changed, pass")
			} else {
				htmChanged = true
				fmt.Println("Normal Files")
				fmt.Println("    Add:    " + strconv.Itoa(addNormalFile))
				fmt.Println("    Update: " + strconv.Itoa(updateNormalFile))
				fmt.Println("    Delete: " + strconv.Itoa(deleteNormalFile))
			}
		}
	}

	fmt.Println("Check and Update monitor folder success")

	if mdChanged == true || htmChanged == true || linkChanged == true || fileChanged == true {

		fmt.Println("Now will compile the site")
		watch := MonitorCompile(smp.SiteFolderPath)
		//Compile
		_, errCompile := Monitor.IPSC_Compile(smp.SiteFolderPath, smp.SiteTitle, indexPageSize)

		if errCompile != nil {
			var errMsg = "RunMonitor: Update monitor folder Success, but cannot compile " + errCompile.Error()
			Utils.Logger.Println(errMsg)
			return false, errCompile
		}
		StopMonitor(watch, smp.SiteFolderPath)
		fmt.Println("Compile Success")
	} else {
		fmt.Println("No file changed since previous monitor, pass ")
	}
	//Save Monitor
	smp.SaveToFile(monitorDefinitionFilePath)

	fmt.Println("***************")
	return true, nil
}

func Run() {
	Utils.InitLogger()
	var cp CommandParser
	bParse := cp.ParseCommand()
	if bParse == true {
		_, errRet := Dispatch(cp)
		if errRet != nil {
			Utils.Logger.Println(errRet.Error())
		}
	}
	fmt.Println("")
	fmt.Println("Note:If ipsd failed, read ipsd.log")
}

func main() {
	Run()
	//test()
}
