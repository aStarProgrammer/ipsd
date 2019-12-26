package main

import (
	"bytes"
	"fmt"
	"ipsd/Monitor"
	"os"
	"os/exec"
	"strconv"
)

func testExportSite() {
	var siteFolder, siteTitle, exportFolder string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	exportFolder = "F:\\WatchDogSpace"

	markdownPages, htmlPages, linkPages, bExportSite, errExportSite := Monitor.IPSC_ExportSite(siteFolder, siteTitle, exportFolder)

	if bExportSite && errExportSite == nil {
		fmt.Println("Markdown")
		for index, markdownPage := range markdownPages {
			fmt.Println("  " + strconv.Itoa(index))
			fmt.Println("    FilePath: " + markdownPage.FilePath)
			fmt.Println("    ID: " + markdownPage.ID)
			fmt.Println("    LastModified: " + markdownPage.LastModified)
			fmt.Println("--------------")
		}

		fmt.Println("Html")
		for index, htmlPage := range htmlPages {
			fmt.Println("  " + strconv.Itoa(index))
			fmt.Println("    FilePath: " + htmlPage.FilePath)
			fmt.Println("    ID: " + htmlPage.ID)
			fmt.Println("    LastModified: " + htmlPage.LastModified)
			fmt.Println("--------------")
		}

		fmt.Println("Link")
		for index, linkPage := range linkPages {
			fmt.Println("  " + strconv.Itoa(index))
			fmt.Println("    Url: " + linkPage.Url)
			fmt.Println("    ID: " + linkPage.ID)
			fmt.Println("    Title: " + linkPage.Title)
			fmt.Println("--------------")
		}

	} else {
		fmt.Println(errExportSite.Error())
	}
}

func testCmd() {
	var cmd = "testCmd"
	var arg1 = `-SpaceStr`
	var arg2 = `"A B C"`
	var ipscCmd *exec.Cmd

	ipscCmd = exec.Command(cmd)

	ipscCmd.Args = append(ipscCmd.Args, arg1)
	ipscCmd.Args = append(ipscCmd.Args, arg2)

	var stdoutput bytes.Buffer
	var stderr bytes.Buffer

	ipscCmd.Stdout = &stdoutput
	ipscCmd.Stderr = &stderr

	errIPSCCmd := ipscCmd.Run()
	if errIPSCCmd != nil {
		fmt.Println(fmt.Sprint(errIPSCCmd) + " : " + stderr.String())

	}
	fmt.Println(stdoutput.String())
}

func testLink() {
	var srcFolder = "F:\\SiteOutputFolder"
	var linkFolder = "F:\\A B\\Test Folder"

	errLink := os.Symlink(srcFolder, linkFolder)

	if errLink != nil {
		fmt.Println("MakeSoftLink: " + errLink.Error())
	}
}

func testAddMarkdown() {
	var siteFolder, siteTitle, pagePath, pageTitle, pageAuthor, titleImage string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	pagePath = "F:\\MarkdownWorkspace\\A1.md"
	pageTitle = "Test Markdown Page"
	pageAuthor = "Chao"
	titleImage = "F:\\MarkdownWorkspace\\1.png"

	output, _, errAddMarkdown := Monitor.IPSC_AddMarkdown(siteFolder, siteTitle, pagePath, pageTitle, pageAuthor, titleImage, true)

	if errAddMarkdown != nil {
		fmt.Println(errAddMarkdown.Error())
	}
	fmt.Println(output)

}

func testHtml() {
	var siteFolder, siteTitle, pagePath, pageTitle, pageAuthor, titleImage string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	pagePath = "F:\\MarkdownWorkspace\\H1.html"
	pageTitle = "Test Markdown Page"
	pageAuthor = "Chao"
	titleImage = "F:\\MarkdownWorkspace\\4.png"

	output, _, errAddMarkdown := Monitor.IPSC_AddHtml(siteFolder, siteTitle, pagePath, pageTitle, pageAuthor, titleImage, true)

	if errAddMarkdown != nil {
		fmt.Println(errAddMarkdown.Error())
	}
	fmt.Println(output)
}

func testAddLink() {
	var siteFolder, siteTitle, linkUrl, pageTitle, pageAuthor, titleImage string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	linkUrl = "http://www.sina.com"
	pageTitle = "Test Markdown Page"
	pageAuthor = "Chao"
	titleImage = "F:\\MarkdownWorkspace\\6.png"

	output, _, errAddMarkdown := Monitor.IPSC_AddLink(siteFolder, siteTitle, linkUrl, pageTitle, pageAuthor, titleImage, true)

	if errAddMarkdown != nil {
		fmt.Println(errAddMarkdown.Error())
	}
	fmt.Println(output)
}

func testUpdateMarkdownOrHtml() {
	var siteFolder, siteTitle, pagePath, pageID, pageTitle, pageAuthor, titleImage string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	pagePath = "F:\\MarkdownWorkspace\\A1.md"
	pageTitle = "Test Markdown Page"
	pageAuthor = "Chao"
	titleImage = "F:\\MarkdownWorkspace\\1.png"
	pageID = "c06d5ad55a8c7907109ff84382bf7b15"

	bUpdate, errUpdateMarkdown := Monitor.IPSC_UpdateMarkdownOrHtml(siteFolder, siteTitle, pageID, pagePath, pageTitle, pageAuthor, titleImage, true)

	if errUpdateMarkdown != nil {
		fmt.Println(errUpdateMarkdown.Error())
	}

	if errUpdateMarkdown == nil && bUpdate {
		fmt.Println("Success")
	}
}

func testUpdateLink() {
	var siteFolder, siteTitle, linkUrl, pageID, pageTitle, pageAuthor, titleImage string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	linkUrl = "https://www.baidu.com"
	pageTitle = "Test Markdown Page"
	pageAuthor = "Chao"
	titleImage = "F:\\MarkdownWorkspace\\1.png"
	pageID = "5ab08f36ce1c0da22c47ce038e41b95d"

	bUpdate, errUpdate := Monitor.IPSC_UpdateLink(siteFolder, siteTitle, pageID, linkUrl, pageTitle, pageAuthor, titleImage, true)

	if errUpdate != nil {
		fmt.Println(errUpdate.Error())
	}

	if errUpdate == nil && bUpdate {
		fmt.Println("Success")
	}
}

func testDeletePage() {
	var siteFolder, siteTitle, pageID string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	pageID = "1cadbf25e1b192cb5ef0299d018af8de"

	bDelete, errDelete := Monitor.IPSC_DeletePage(siteFolder, siteTitle, pageID)

	if errDelete != nil {
		fmt.Println(errDelete.Error())
	}

	if errDelete == nil && bDelete {
		fmt.Println("Success")
	}
}

func testCreateMarkdown() {
	var siteFolder, siteTitle, pagePath, markdownType string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	pagePath = "F:\\MarkdownWorkspace\\_A1.md"
	markdownType = "News"

	bDelete, errDelete := Monitor.IPSC_CreateMarkdown(siteFolder, siteTitle, pagePath, markdownType)

	if errDelete != nil {
		fmt.Println(errDelete.Error())
	}

	if errDelete == nil && bDelete {
		fmt.Println("Success ")
	}
}

func testCompile() {
	var siteFolder, siteTitle, indexPageSize string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	indexPageSize = "VerySmall"

	bCompile, errCompile := Monitor.IPSC_Compile(siteFolder, siteTitle, indexPageSize)

	if errCompile != nil {
		fmt.Println(errCompile.Error())
	}

	if errCompile == nil && bCompile {
		fmt.Println("Success ")
	}
}

func testReadMarkdownFile() {
	var filePath = "F:\\WatchdogSpace\\Templates\\News.md"

	mdProperties, _, errProperties := Monitor.ReadMarkdownPageProperties(filePath)

	if errProperties != nil {
		var errMsg = "Cannot read properties"
		fmt.Println(errMsg)
		return
	}

	fmt.Println(mdProperties.Title)
	fmt.Println(mdProperties.Author)
	fmt.Println(mdProperties.Description)
	fmt.Println(mdProperties.IsTop)
}

func testReadHtmlFileProperties() {
	var filePath = "F:\\WatchdogSpace\\Html\\H1.html"

	htmProperties, _, errProperties := Monitor.ReadHtmlProperties(filePath)

	if errProperties != nil {
		var errMsg = "Cannot read properties"
		fmt.Println(errMsg)
		return
	}

	fmt.Println(htmProperties.Title)
	fmt.Println(htmProperties.Author)
	fmt.Println(htmProperties.Description)
	fmt.Println(htmProperties.IsTop)
}

func testReadLinksFile() {
	var filePath = "F:\\WatchdogSpace\\Link\\Link.txt"

	links, errProperties := Monitor.ReadLinksFromFile(filePath)

	if errProperties != nil {
		var errMsg = "Cannot read properties"
		fmt.Println(errMsg)
		return
	}

	for _, link := range links {
		fmt.Println(link.Url)
		fmt.Println(link.ID)
		fmt.Println(link.Title)
		fmt.Println("--------")
	}

}

func testNewMonitor() {
	var siteFolder, siteTitle, monitorFolder string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	monitorFolder = "F:\\WatchDogSpace"

	_, errNewMonitor := NewMonitor(monitorFolder, siteFolder, siteTitle)

	if errNewMonitor != nil {
		fmt.Println(errNewMonitor.Error())
	}

}

func testRunMonitor() {
	var monitorFolderPath = "F:\\WatchDogSpace"
	var indexPageSize = "VerySmall"

	_, errRunMonitor := RunMonitor(monitorFolderPath, indexPageSize)

	if errRunMonitor != nil {
		fmt.Println(errRunMonitor.Error())
	}
}

func testAddFile() {
	var siteFolder, siteTitle, filePath string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	filePath = "F:\\MarkdownWorkspace"

	bAdd, errAdd := Monitor.IPSC_AddFile(siteFolder, siteTitle, filePath)

	if errAdd != nil {
		fmt.Println(errAdd.Error())
	}

	if errAdd == nil && bAdd {
		fmt.Println("Success ")
	}
}

func testListFile() {
	var siteFolder, siteTitle string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"

	bList, errList := Monitor.IPSC_ListFile(siteFolder, siteTitle)

	if errList != nil {
		fmt.Println(errList.Error())
	}

	if errList == nil && bList {
		fmt.Println("Success ")
	}
}

func testDeleteFile() {
	var siteFolder, siteTitle, filePath string
	siteFolder = "F:\\TestSite"
	siteTitle = "Test Site"
	filePath = ".\\Files\\MarkdownWorkspace\\A1.md"

	bDelete, errDelete := Monitor.IPSC_DeleteFile(siteFolder, siteTitle, filePath)

	if errDelete != nil {
		fmt.Println(errDelete.Error())
	}

	if errDelete == nil && bDelete {
		fmt.Println("Success ")
	}

}

func test() {
	testRunMonitor()
}
