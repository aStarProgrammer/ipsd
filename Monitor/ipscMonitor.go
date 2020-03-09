package Monitor

import (
	"bytes"
	"errors"
	"fmt"
	"ipsd/Utils"
	"os/exec"
	"strconv"
	"strings"
)

const (
	COMMAND_ADDPAGE           = "AddPage"
	COMMAND_UPDATEPAGE        = "UpdatePage"
	COMMAND_DELETEPAGE        = "DeletePage"
	COMMAND_CREATEMARKDOWN    = "CreateMarkdown"
	COMMAND_COMPILE           = "Compile"
	COMMAND_EXPORTSOURCEPAGES = "ExportSourcePages"
	COMMAND_LISTOUTPUTPAGES   = "ListSourcePages"
	COMMAND_LISTPAGE          = "ListPage"
	COMMAND_ADDFILE           = "AddFile"
	COMMAND_DELETEFILE        = "DeleteFile"
	COMMAND_LISTFILE          = "ListFile"
)

func RunIPSCCommand(ipscCmd *exec.Cmd) (string, error) {
	var stdoutput bytes.Buffer
	var stderr bytes.Buffer

	ipscCmd.Stdout = &stdoutput
	ipscCmd.Stderr = &stderr

	errIPSCCmd := ipscCmd.Run()
	if errIPSCCmd != nil {
		Utils.Logger.Println(fmt.Sprint(errIPSCCmd) + " : " + stderr.String())
		return "", errIPSCCmd
	}
	return string(stdoutput.String()), nil
}

//ipsc -Command "ExportSourcePages" -SiteFolder -SiteTitle -ExportFolder
func IPSC_ExportSite(siteFolder, siteTitle, exportFolder string) ([]MarkdownPage, []HtmlPage, []LinkPage, bool, error) {

	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_ExportSite: siteFolder is empty"
		Utils.Logger.Println(errMsg)
		return nil, nil, nil, false, errors.New(errMsg)
	}

	if Utils.PathIsExist(exportFolder) == false {
		var errMsg = "ipscMonitor.IPSC_ExportSite: exportFolder is empty"
		Utils.Logger.Println(errMsg)
		return nil, nil, nil, false, errors.New(errMsg)
	}

	//ipsc -Command "ExportSourcePages" -SiteFolder -SiteTitle -ExportFolder
	var exportCmd = exec.Command("ipsc")
	exportCmd.Args = append(exportCmd.Args, "-Command")
	exportCmd.Args = append(exportCmd.Args, COMMAND_EXPORTSOURCEPAGES)
	exportCmd.Args = append(exportCmd.Args, "-SiteFolder")
	exportCmd.Args = append(exportCmd.Args, siteFolder)
	exportCmd.Args = append(exportCmd.Args, "-SiteTitle")
	exportCmd.Args = append(exportCmd.Args, siteTitle)
	exportCmd.Args = append(exportCmd.Args, "-ExportFolder")
	exportCmd.Args = append(exportCmd.Args, exportFolder)

	ipscOutput, errRunCmd := RunIPSCCommand(exportCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return nil, nil, nil, false, errRunCmd
	}

	//Get input content
	var markdownPages []MarkdownPage
	var htmlPages []HtmlPage
	var linkPages []LinkPage

	if ipscOutput != "" {
		if strings.Contains(ipscOutput, "`") && strings.HasPrefix(ipscOutput, "Exported:") {
			var exportDatas = strings.Split(ipscOutput, "`")
			if len(exportDatas) < 4 {
				/*At least it will export
				Exported:0`
				Markdown:0`
				Html:0`
				Link:0`
				*/
				var errMsg = "ipscMonitor.IPSC_ExportSite: Export data error"
				Utils.Logger.Println(errMsg)
				return markdownPages, htmlPages, linkPages, false, errors.New(errMsg)
			}

			var sAllCount = exportDatas[0]
			var sMarkdownCount = exportDatas[1]

			allCount, errAllCount := strconv.Atoi(strings.Split(sAllCount, ":")[1])
			if errAllCount != nil {
				var errMsg = "ipscMonitor.IPSC_ExportSite: Cannot get allCount"
				Utils.Logger.Println(errMsg)
				return markdownPages, htmlPages, linkPages, false, errors.New(errMsg)
			}

			markdownCount, errMarkdownCount := strconv.Atoi(strings.Split(sMarkdownCount, ":")[1])
			if errMarkdownCount != nil {
				var errMsg = "ipscMonitor.IPSC_ExportSite: Cannot get MarkdownCount"
				Utils.Logger.Println(errMsg)
				return markdownPages, htmlPages, linkPages, false, errors.New(errMsg)
			}

			var sHtmlCount = exportDatas[2+markdownCount]
			htmlCount, errHtmlCount := strconv.Atoi(strings.Split(sHtmlCount, ":")[1])
			if errHtmlCount != nil {
				var errMsg = "ipscMonitor.IPSC_ExportSite: Cannot get HtmlCount"
				Utils.Logger.Println(errMsg)
				return markdownPages, htmlPages, linkPages, false, errors.New(errMsg)
			}

			var sLinkCount = exportDatas[2+markdownCount+htmlCount+1]
			linkCount, errLinkCount := strconv.Atoi(strings.Split(sLinkCount, ":")[1])
			if errLinkCount != nil {
				var errMsg = "ipscMonitor.IPSC_ExportSite: Cannot get LinkCount"
				Utils.Logger.Println(errMsg)
				return markdownPages, htmlPages, linkPages, false, errors.New(errMsg)
			}

			if allCount != markdownCount+htmlCount+linkCount {
				var errMsg = "ipscMonitor.IPSC_ExportSite: allCount != markdownCount+htmlCount+linkCount"
				Utils.Logger.Println(errMsg)
				return markdownPages, htmlPages, linkPages, false, errors.New(errMsg)
			}

			for index := 2; index < 2+markdownCount; index++ {
				var sMarkdownData = exportDatas[index]
				var markdownDatas = strings.Split(sMarkdownData, "|")
				if len(markdownDatas) != 3 {
					var errMsg = "ipscMonitor.IPSC_ExportSite: Retrive Markdown Data failed " + sMarkdownCount
					Utils.Logger.Println(errMsg)
					return markdownPages, htmlPages, linkPages, false, errors.New(errMsg)
				}

				var mdPage MarkdownPage
				mdPage.FilePath = strings.TrimSpace(markdownDatas[0])
				mdPage.ID = strings.TrimSpace(markdownDatas[1])
				mdPage.LastModified = strings.TrimSpace(markdownDatas[2])

				markdownPages = append(markdownPages, mdPage)

			}

			for index := 2 + markdownCount + 1; index < 2+markdownCount+1+htmlCount; index++ {
				var sHtmlData = exportDatas[index]
				var htmlDatas = strings.Split(sHtmlData, "|")
				if len(htmlDatas) != 3 {
					var errMsg = "ipscMonitor.IPSC_ExportSite: Retrive Html Data failed " + sHtmlData
					Utils.Logger.Println(errMsg)
					return markdownPages, htmlPages, linkPages, false, errors.New(errMsg)
				}

				var htmlPage HtmlPage
				htmlPage.FilePath = strings.TrimSpace(htmlDatas[0])
				htmlPage.ID = strings.TrimSpace(htmlDatas[1])
				htmlPage.LastModified = strings.TrimSpace(htmlDatas[2])

				htmlPages = append(htmlPages, htmlPage)
			}

			for index := 2 + markdownCount + 1 + htmlCount + 1; index < len(exportDatas)-1; index++ {
				var sLinkData = exportDatas[index]
				var linkDatas = strings.Split(sLinkData, "|")
				if len(linkDatas) != 4 {
					var errMsg = "ipscMonitor.IPSC_ExportSite: Retrive Link Data failed " + sLinkData
					Utils.Logger.Println(errMsg)
					return markdownPages, htmlPages, linkPages, false, errors.New(errMsg)
				}

				var linkPage LinkPage
				linkPage.Url = strings.TrimSpace(linkDatas[0])
				linkPage.ID = strings.TrimSpace(linkDatas[1])
				linkPage.Title = strings.TrimSpace(linkDatas[2])
				bIsTop, errParseBool := strconv.ParseBool(linkDatas[3])

				if errParseBool == nil {
					linkPage.IsTop = bIsTop
				} else {
					linkPage.IsTop = false
				}
				linkPages = append(linkPages, linkPage)
			}

		}
	}

	return markdownPages, htmlPages, linkPages, true, nil
}

//ipsc -Command "AddPage" -SiteFolder -SiteTitle -PagePath -LinkUrl -PageType -PageTitle -PageAuthor -TitleImage -IsTop
func IPSC_AddMarkdown(siteFolder, siteTitle, pagePath, pageTitle, pageAuthor, titleImage string, isTop bool) (string, bool, error) {
	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_AddMarkdown: siteFolder is empty"
		Utils.Logger.Println(errMsg)
		return "", false, errors.New(errMsg)
	}

	var addCmd = exec.Command("ipsc")
	addCmd.Args = append(addCmd.Args, "-Command")
	addCmd.Args = append(addCmd.Args, COMMAND_ADDPAGE)
	addCmd.Args = append(addCmd.Args, "-SiteFolder")
	addCmd.Args = append(addCmd.Args, siteFolder)
	addCmd.Args = append(addCmd.Args, "-SiteTitle")
	addCmd.Args = append(addCmd.Args, siteTitle)
	addCmd.Args = append(addCmd.Args, "-PagePath")
	addCmd.Args = append(addCmd.Args, pagePath)
	addCmd.Args = append(addCmd.Args, "-PageTitle")
	addCmd.Args = append(addCmd.Args, pageTitle)
	addCmd.Args = append(addCmd.Args, "-PageType")
	addCmd.Args = append(addCmd.Args, "MARKDOWN")
	addCmd.Args = append(addCmd.Args, "-PageAuthor")
	addCmd.Args = append(addCmd.Args, pageAuthor)
	addCmd.Args = append(addCmd.Args, "-TitleImage")
	addCmd.Args = append(addCmd.Args, titleImage)

	addCmd.Args = append(addCmd.Args, "-IsTop")
	if isTop {
		addCmd.Args = append(addCmd.Args, "true")
	} else {
		addCmd.Args = append(addCmd.Args, "false")
	}

	ipscOutput, errRunCmd := RunIPSCCommand(addCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return "", false, errRunCmd
	}

	//Get input content

	if ipscOutput != "" {
		if strings.HasPrefix(ipscOutput, "Add Success,") {
			//fmt.Println(ipscOutput)
			ipscOutput = ipscOutput[strings.LastIndex(ipscOutput, " ")+1:]
			ipscOutput = strings.TrimSpace(ipscOutput)
			return ipscOutput, true, nil
		} else {
			return "", false, errors.New(ipscOutput)
		}
	}

	return "", false, errors.New(ipscOutput)
}

func IPSC_AddHtml(siteFolder, siteTitle, pagePath, pageTitle, pageAuthor, titleImage string, isTop bool) (string, bool, error) {
	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_AddHtml: siteFolder is empty"
		Utils.Logger.Println(errMsg)
		return "", false, errors.New(errMsg)
	}

	var addCmd = exec.Command("ipsc")
	addCmd.Args = append(addCmd.Args, "-Command")
	addCmd.Args = append(addCmd.Args, COMMAND_ADDPAGE)
	addCmd.Args = append(addCmd.Args, "-SiteFolder")
	addCmd.Args = append(addCmd.Args, siteFolder)
	addCmd.Args = append(addCmd.Args, "-SiteTitle")
	addCmd.Args = append(addCmd.Args, siteTitle)
	addCmd.Args = append(addCmd.Args, "-PagePath")
	addCmd.Args = append(addCmd.Args, pagePath)
	addCmd.Args = append(addCmd.Args, "-PageTitle")
	addCmd.Args = append(addCmd.Args, pageTitle)
	addCmd.Args = append(addCmd.Args, "-PageType")
	addCmd.Args = append(addCmd.Args, "HTML")
	addCmd.Args = append(addCmd.Args, "-PageAuthor")
	addCmd.Args = append(addCmd.Args, pageAuthor)
	addCmd.Args = append(addCmd.Args, "-TitleImage")
	addCmd.Args = append(addCmd.Args, titleImage)

	addCmd.Args = append(addCmd.Args, "-IsTop")
	if isTop {
		addCmd.Args = append(addCmd.Args, "true")
	} else {
		addCmd.Args = append(addCmd.Args, "false")
	}

	ipscOutput, errRunCmd := RunIPSCCommand(addCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return "", false, errRunCmd
	}

	//Get input content

	if ipscOutput != "" {
		if strings.HasPrefix(ipscOutput, "Add Success,") {
			if strings.HasSuffix(ipscOutput, "Done") {
				ipscOutput = ipscOutput[:strings.Index(ipscOutput, "Done")]
			}
			ipscOutput = ipscOutput[strings.LastIndex(ipscOutput, " ")+1:]
			ipscOutput = strings.TrimSpace(ipscOutput)
			return ipscOutput, true, nil
		} else {
			return "", false, errors.New(ipscOutput)
		}
	}

	return "", false, errors.New(ipscOutput)
}

func IPSC_AddLink(siteFolder, siteTitle, linkUrl, pageTitle, pageAuthor, titleImage string, isTop bool) (string, bool, error) {
	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_AddLink: siteFolder is empty"
		Utils.Logger.Println(errMsg)
		return "", false, errors.New(errMsg)
	}

	var addCmd = exec.Command("ipsc")
	addCmd.Args = append(addCmd.Args, "-Command")
	addCmd.Args = append(addCmd.Args, COMMAND_ADDPAGE)
	addCmd.Args = append(addCmd.Args, "-SiteFolder")
	addCmd.Args = append(addCmd.Args, siteFolder)
	addCmd.Args = append(addCmd.Args, "-SiteTitle")
	addCmd.Args = append(addCmd.Args, siteTitle)
	addCmd.Args = append(addCmd.Args, "-LinkUrl")
	addCmd.Args = append(addCmd.Args, linkUrl)
	addCmd.Args = append(addCmd.Args, "-PageTitle")
	addCmd.Args = append(addCmd.Args, pageTitle)
	addCmd.Args = append(addCmd.Args, "-PageType")
	addCmd.Args = append(addCmd.Args, "LINK")
	addCmd.Args = append(addCmd.Args, "-PageAuthor")
	addCmd.Args = append(addCmd.Args, pageAuthor)
	addCmd.Args = append(addCmd.Args, "-TitleImage")
	addCmd.Args = append(addCmd.Args, titleImage)

	addCmd.Args = append(addCmd.Args, "-IsTop")
	if isTop {
		addCmd.Args = append(addCmd.Args, "true")
	} else {
		addCmd.Args = append(addCmd.Args, "false")
	}

	ipscOutput, errRunCmd := RunIPSCCommand(addCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return "", false, errRunCmd
	}

	//Get input content

	if ipscOutput != "" {
		if strings.HasPrefix(ipscOutput, "Add Success,") {
			//fmt.Println(ipscOutput)
			if strings.HasSuffix(ipscOutput, "Done") {
				ipscOutput = ipscOutput[:strings.Index(ipscOutput, "Done")]
			}
			ipscOutput = ipscOutput[strings.LastIndex(ipscOutput, " ")+1:]
			ipscOutput = strings.TrimSpace(ipscOutput)

			if Utils.IsGuid(ipscOutput) {
				return ipscOutput, true, nil
			} else {
				return "", false, errors.New(ipscOutput)
			}

		} else {
			Utils.Logger.Println(ipscOutput)
			return "", false, errors.New(ipscOutput)
		}
	}

	return ipscOutput, false, nil
}

func IPSC_UpdateMarkdownOrHtml(siteFolder, siteTitle, pageID, pagePath, pageTitle, pageAuthor, titleImage string, isTop bool) (bool, error) {
	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_UpdateMarkdownOrHtml: siteFolder is empty"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var updateCmd = exec.Command("ipsc")
	updateCmd.Args = append(updateCmd.Args, "-Command")
	updateCmd.Args = append(updateCmd.Args, COMMAND_UPDATEPAGE)
	updateCmd.Args = append(updateCmd.Args, "-SiteFolder")
	updateCmd.Args = append(updateCmd.Args, siteFolder)
	updateCmd.Args = append(updateCmd.Args, "-SiteTitle")
	updateCmd.Args = append(updateCmd.Args, siteTitle)
	updateCmd.Args = append(updateCmd.Args, "-PageID")
	updateCmd.Args = append(updateCmd.Args, pageID)
	updateCmd.Args = append(updateCmd.Args, "-PagePath")
	updateCmd.Args = append(updateCmd.Args, pagePath)
	updateCmd.Args = append(updateCmd.Args, "-PageTitle")
	updateCmd.Args = append(updateCmd.Args, pageTitle)
	updateCmd.Args = append(updateCmd.Args, "-PageAuthor")
	updateCmd.Args = append(updateCmd.Args, pageAuthor)
	updateCmd.Args = append(updateCmd.Args, "-TitleImage")
	updateCmd.Args = append(updateCmd.Args, titleImage)

	updateCmd.Args = append(updateCmd.Args, "-IsTop")
	if isTop {
		updateCmd.Args = append(updateCmd.Args, "true")
	} else {
		updateCmd.Args = append(updateCmd.Args, "false")
	}

	ipscOutput, errRunCmd := RunIPSCCommand(updateCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return false, errRunCmd
	}

	//Get input content

	if ipscOutput != "" {
		if strings.HasPrefix(ipscOutput, "Update Success") {
			return true, nil
		}
	}

	return false, nil
}

func IPSC_UpdateLink(siteFolder, siteTitle, pageID, linkUrl, pageTitle, pageAuthor, titleImage string, isTop bool) (bool, error) {
	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_UpdateLink: siteFolder is empty"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var updateCmd = exec.Command("ipsc")
	updateCmd.Args = append(updateCmd.Args, "-Command")
	updateCmd.Args = append(updateCmd.Args, COMMAND_UPDATEPAGE)
	updateCmd.Args = append(updateCmd.Args, "-SiteFolder")
	updateCmd.Args = append(updateCmd.Args, siteFolder)
	updateCmd.Args = append(updateCmd.Args, "-SiteTitle")
	updateCmd.Args = append(updateCmd.Args, siteTitle)
	updateCmd.Args = append(updateCmd.Args, "-PageID")
	updateCmd.Args = append(updateCmd.Args, pageID)
	updateCmd.Args = append(updateCmd.Args, "-LinkUrl")
	updateCmd.Args = append(updateCmd.Args, linkUrl)
	updateCmd.Args = append(updateCmd.Args, "-PageTitle")
	updateCmd.Args = append(updateCmd.Args, pageTitle)
	updateCmd.Args = append(updateCmd.Args, "-PageAuthor")
	updateCmd.Args = append(updateCmd.Args, pageAuthor)
	updateCmd.Args = append(updateCmd.Args, "-TitleImage")
	updateCmd.Args = append(updateCmd.Args, titleImage)

	updateCmd.Args = append(updateCmd.Args, "-IsTop")
	if isTop {
		updateCmd.Args = append(updateCmd.Args, "true")
	} else {
		updateCmd.Args = append(updateCmd.Args, "false")
	}

	ipscOutput, errRunCmd := RunIPSCCommand(updateCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return false, errRunCmd
	}

	//Get input content

	if ipscOutput != "" {
		if strings.HasPrefix(ipscOutput, "Update Success") {
			return true, nil
		} else {
			Utils.Logger.Println(ipscOutput)
		}
	}

	return false, nil
}

//ipsc -Command "DeletePage"  -SiteFolder "F:\TestSite" -SiteTitle "StarSite" -PageID "fc0f8d635ebb04d1c9393a722e8fc185" -RestorePage
func IPSC_DeletePage(siteFolder, siteTitle, pageID string) (bool, error) {
	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_DeletePage: siteFolder is empty"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var deleteCmd = exec.Command("ipsc")
	deleteCmd.Args = append(deleteCmd.Args, "-Command")
	deleteCmd.Args = append(deleteCmd.Args, COMMAND_DELETEPAGE)
	deleteCmd.Args = append(deleteCmd.Args, "-SiteFolder")
	deleteCmd.Args = append(deleteCmd.Args, siteFolder)
	deleteCmd.Args = append(deleteCmd.Args, "-SiteTitle")
	deleteCmd.Args = append(deleteCmd.Args, siteTitle)
	deleteCmd.Args = append(deleteCmd.Args, "-PageID")
	deleteCmd.Args = append(deleteCmd.Args, pageID)
	deleteCmd.Args = append(deleteCmd.Args, "-RestorePage")
	deleteCmd.Args = append(deleteCmd.Args, "false")

	ipscOutput, errRunCmd := RunIPSCCommand(deleteCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return false, errRunCmd
	}

	//Get input content

	if ipscOutput != "" {
		if strings.HasPrefix(ipscOutput, "Delete Success") {
			return true, nil
		}
	}

	return false, nil
}

//ipsc -Command "CreateMarkdown" -SiteFolder "F:\TestSite" -SiteTitle "StarSite" -PagePath "F:\MarkdownWorkspace\_A1.md" -MarkdownType "News"
func IPSC_CreateMarkdown(siteFolder, siteTitle, pagePath, markdownType string) (bool, error) {
	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_CreateMarkdown: siteFolder is empty"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var createMdCmd = exec.Command("ipsc")
	createMdCmd.Args = append(createMdCmd.Args, "-Command")
	createMdCmd.Args = append(createMdCmd.Args, COMMAND_CREATEMARKDOWN)
	createMdCmd.Args = append(createMdCmd.Args, "-SiteFolder")
	createMdCmd.Args = append(createMdCmd.Args, siteFolder)
	createMdCmd.Args = append(createMdCmd.Args, "-SiteTitle")
	createMdCmd.Args = append(createMdCmd.Args, siteTitle)
	createMdCmd.Args = append(createMdCmd.Args, "-PagePath")
	createMdCmd.Args = append(createMdCmd.Args, pagePath)
	createMdCmd.Args = append(createMdCmd.Args, "-MarkdownType")
	createMdCmd.Args = append(createMdCmd.Args, markdownType)

	ipscOutput, errRunCmd := RunIPSCCommand(createMdCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return false, errRunCmd
	}

	//Get input content

	if ipscOutput != "" {
		if strings.HasPrefix(ipscOutput, "CreateMarkdown Success") {
			return true, nil
		}
	}

	return false, nil
}

//ipsc -Command "Compile" -SiteFolder "F:\TestSite" -SiteTitle "StarSite" -IndexPageSize "Normal"
func IPSC_Compile(siteFolder, siteTitle, indexPageSize string) (bool, error) {
	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_Compile: siteFolder is empty"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var compileCmd = exec.Command("ipsc")
	compileCmd.Args = append(compileCmd.Args, "-Command")
	compileCmd.Args = append(compileCmd.Args, COMMAND_COMPILE)
	compileCmd.Args = append(compileCmd.Args, "-SiteFolder")
	compileCmd.Args = append(compileCmd.Args, siteFolder)
	compileCmd.Args = append(compileCmd.Args, "-SiteTitle")
	compileCmd.Args = append(compileCmd.Args, siteTitle)
	compileCmd.Args = append(compileCmd.Args, "-IndexPageSize")
	compileCmd.Args = append(compileCmd.Args, indexPageSize)

	ipscOutput, errRunCmd := RunIPSCCommand(compileCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return false, errRunCmd
	}

	//Get input content

	if ipscOutput != "" {
		if strings.HasPrefix(ipscOutput, "Compile Summary") {
			return true, nil
		}
	}

	return false, nil
}

func IPSC_AddFile(siteFolder, siteTitle, filePath string) (bool, error) {

	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_AddFile: siteFolder not exist"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	if Utils.PathIsExist(filePath) == false {
		var errMsg = "ipscMonitor.IPSC_AddFile: " + filePath + " not exist"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var ipscCmd = exec.Command("ipsc")
	ipscCmd.Args = append(ipscCmd.Args, "-Command")
	ipscCmd.Args = append(ipscCmd.Args, COMMAND_ADDFILE)
	ipscCmd.Args = append(ipscCmd.Args, "-SiteFolder")
	ipscCmd.Args = append(ipscCmd.Args, siteFolder)
	ipscCmd.Args = append(ipscCmd.Args, "-SiteTitle")
	ipscCmd.Args = append(ipscCmd.Args, siteTitle)
	ipscCmd.Args = append(ipscCmd.Args, "-FilePath")
	ipscCmd.Args = append(ipscCmd.Args, filePath)

	ipscOutput, errRunCmd := RunIPSCCommand(ipscCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return false, errRunCmd
	}

	//Get input content

	if ipscOutput != "" {
		if strings.HasPrefix(ipscOutput, "Add File Success") {
			return true, nil
		}
	}

	return true, nil
}

func IPSC_DeleteFile(siteFolder, siteTitle, filePath string) (bool, error) {
	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_DeleteFile: siteFolder is empty"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var ipscCmd = exec.Command("ipsc")
	ipscCmd.Args = append(ipscCmd.Args, "-Command")
	ipscCmd.Args = append(ipscCmd.Args, COMMAND_DELETEFILE)
	ipscCmd.Args = append(ipscCmd.Args, "-SiteFolder")
	ipscCmd.Args = append(ipscCmd.Args, siteFolder)
	ipscCmd.Args = append(ipscCmd.Args, "-SiteTitle")
	ipscCmd.Args = append(ipscCmd.Args, siteTitle)
	ipscCmd.Args = append(ipscCmd.Args, "-FilePath")
	ipscCmd.Args = append(ipscCmd.Args, filePath)

	ipscOutput, errRunCmd := RunIPSCCommand(ipscCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return false, errRunCmd
	}

	//Get input content

	if ipscOutput != "" {
		fmt.Println(ipscOutput)
		if strings.HasPrefix(ipscOutput, "Delete Success") {
			return true, nil
		}
	}

	return true, nil
}

func IPSC_ListFile(siteFolder, siteTitle string) (bool, error) {
	if Utils.PathIsExist(siteFolder) == false {
		var errMsg = "ipscMonitor.IPSC_ListFile: siteFolder is empty"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var ipscCmd = exec.Command("ipsc")
	ipscCmd.Args = append(ipscCmd.Args, "-Command")
	ipscCmd.Args = append(ipscCmd.Args, COMMAND_LISTFILE)
	ipscCmd.Args = append(ipscCmd.Args, "-SiteFolder")
	ipscCmd.Args = append(ipscCmd.Args, siteFolder)
	ipscCmd.Args = append(ipscCmd.Args, "-SiteTitle")
	ipscCmd.Args = append(ipscCmd.Args, siteTitle)

	ipscOutput, errRunCmd := RunIPSCCommand(ipscCmd)

	if errRunCmd != nil {
		Utils.Logger.Println(errRunCmd.Error())
		return false, errRunCmd
	}

	//Get input content

	if ipscOutput != "" {
		fmt.Println(ipscOutput)
	}

	return true, nil
}
