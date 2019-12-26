package Monitor

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"ipsd/Utils"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type SiteMonitor struct {
	MonitorFolderPath    string
	SiteFolderPath       string
	SiteTitle            string
	LinkFileLastModified string

	MarkdownFiles []MarkdownPage
	HtmlFiles     []HtmlPage
	LinkFiles     []LinkPage
	NormalFiles   []NormalFile
}

func NewSiteMonitor() *SiteMonitor {
	var sm SiteMonitor
	var smp *SiteMonitor
	smp = &sm

	return smp
}

func NewSiteMonitor_WithArgs(monitorFolderPath, siteFolderPath, siteTitle, linkLastModified string) *SiteMonitor {
	var sm SiteMonitor
	var smp *SiteMonitor
	smp = &sm

	smp.MonitorFolderPath = monitorFolderPath
	smp.SiteFolderPath = siteFolderPath
	smp.SiteTitle = siteTitle
	smp.LinkFileLastModified = linkLastModified

	return smp
}

func (smp *SiteMonitor) FromJson(_jsonString string) (bool, error) {
	if "" == _jsonString {
		return false, errors.New("Argument jsonString is null")
	}

	errUnmarshal := json.Unmarshal([]byte(_jsonString), smp)
	if errUnmarshal != nil {
		return false, errUnmarshal
	}
	return true, nil
}

func IsSiteMonitorEmpty(sm SiteMonitor) bool {
	if sm.MonitorFolderPath == "" {
		return true
	}
	return false
}

func (smp *SiteMonitor) ToJson() (string, error) {
	var _jsonbyte []byte

	if smp == nil {
		return "", errors.New("Pointer smp is nil")
	}

	if IsSiteMonitorEmpty(*smp) {
		return "", errors.New("Site Monitor is empty")
	}

	_jsonbyte, err := json.Marshal(*smp)

	return string(_jsonbyte), err
}

func (smp *SiteMonitor) LoadFromFile(filePath string) (bool, error) {
	if "" == filePath {
		return false, errors.New("FilePath is empty")
	}

	bFileExist := Utils.PathIsExist(filePath)

	if false == bFileExist {
		return false, errors.New("File not exist")
	}

	_json, errRead := ioutil.ReadFile(filePath)

	if errRead != nil {
		return false, errors.New("Read File Fail")
	}

	_jsonString := string(_json)

	if "" == _jsonString {
		return false, errors.New("File is empty")
	}

	bUnMarshal, errUnMarshal := smp.FromJson(_jsonString)

	return bUnMarshal, errUnMarshal
}

func (smp *SiteMonitor) SaveToFile(filePath string) (bool, error) {
	if "" == filePath {
		return false, errors.New("FilePath is empty")
	}

	if IsSiteMonitorEmpty(*smp) {
		return false, errors.New("Site Monitor is empty")
	}

	json, errMarshal := smp.ToJson()

	if errMarshal != nil {
		return false, errMarshal
	}

	var errFilePath error
	if !Utils.PathIsExist(filePath) {
		filePath, errFilePath = Utils.MakePath(filePath)
		if errFilePath != nil {
			return false, errors.New("Path nor exist and create parent folder failed")
		}
	}
	//路径分为绝对路径和相对路径
	//create，文件存在则会覆盖原始内容（其实就相当于清空），不存在则创建
	fp, error := os.Create(filePath)
	if error != nil {
		return false, error
	}
	//延迟调用，关闭文件
	defer fp.Close()

	_, errWriteFile := fp.WriteString(json)

	if errWriteFile != nil {
		return false, errors.New("Write json to file failed")
	}

	return true, nil
}

func (smp *SiteMonitor) AddMarkdown_Args(filePath, ID, lastModified string) (bool, error) {
	if smp.GetMarkdown(filePath) != -1 {
		var errMsg = "SiteMonitor.AddMarkdown_Args: Markdown already added"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var mdPage MarkdownPage
	mdPage.FilePath = filePath
	mdPage.ID = ID
	mdPage.LastModified = lastModified

	smp.MarkdownFiles = append(smp.MarkdownFiles, mdPage)

	return true, nil
}

func (smp *SiteMonitor) AddMarkdown(mdPage MarkdownPage) (bool, error) {
	if smp.GetMarkdown(mdPage.FilePath) != -1 {
		var errMsg = "SiteMonitor.AddMarkdown: Markdown already added"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	smp.MarkdownFiles = append(smp.MarkdownFiles, mdPage)

	return true, nil
}

func (smp *SiteMonitor) GetMarkdown(filePath string) int {
	if filePath == "" {
		return -1
	}

	for index, mdPage := range smp.MarkdownFiles {
		var sFilePath = strings.ToLower(mdPage.FilePath)
		var fFilePath = strings.ToLower(filePath)
		if sFilePath == fFilePath {
			return index
		}
	}

	return -1
}

func (smp *SiteMonitor) UpdateMarkdown_Args(filePath, ID, lastModified string) (bool, error) {
	var index = smp.GetMarkdown(filePath)

	if index == -1 {
		var errMsg = "SiteMonitor.UpdateMarkdown_Args: Cannot find markdown with file path " + filePath
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}
	smp.MarkdownFiles[index].FilePath = filePath
	smp.MarkdownFiles[index].ID = ID
	smp.MarkdownFiles[index].LastModified = lastModified
	return true, nil
}

func (smp *SiteMonitor) UpdateMarkdown(mdPage MarkdownPage) (bool, error) {
	var index = smp.GetMarkdown(mdPage.FilePath)

	if index == -1 {
		var errMsg = "SiteMonitor.UpdateMarkdown: Cannot find markdown with file path " + mdPage.FilePath
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}
	smp.MarkdownFiles[index].FilePath = mdPage.FilePath
	smp.MarkdownFiles[index].ID = mdPage.ID
	smp.MarkdownFiles[index].LastModified = mdPage.LastModified
	return true, nil
}

func (smp *SiteMonitor) DeleteMarkdown(filePath string) (bool, error) {
	var index = smp.GetMarkdown(filePath)
	if index == -1 {
		var errMsg = "SiteMonitor.DeleteMarkdown: Cannot find markdown with file path " + filePath
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	smp.MarkdownFiles = append(smp.MarkdownFiles[:index], smp.MarkdownFiles[index+1:]...)
	return true, nil
}

func (smp *SiteMonitor) AddHtml_Args(filePath, ID, lastModified string) (bool, error) {
	if smp.GetHtml(filePath) != -1 {
		var errMsg = "SiteMonitor.AddHtml_Args: Html already added"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var htmPage HtmlPage
	htmPage.FilePath = filePath
	htmPage.ID = ID
	htmPage.LastModified = lastModified

	smp.HtmlFiles = append(smp.HtmlFiles, htmPage)

	return true, nil
}

func (smp *SiteMonitor) AddHtml(htmPage HtmlPage) (bool, error) {
	if smp.GetHtml(htmPage.FilePath) != -1 {
		var errMsg = "SiteMonitor.AddHtml: Html already added"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	smp.HtmlFiles = append(smp.HtmlFiles, htmPage)

	return true, nil
}

func (smp *SiteMonitor) GetHtml(filePath string) int {
	if filePath == "" {
		return -1
	}

	for index, htmPage := range smp.HtmlFiles {
		var sFilePath = strings.ToLower(htmPage.FilePath)
		var fFilePath = strings.ToLower(filePath)
		if sFilePath == fFilePath {
			return index
		}
	}

	return -1
}

func (smp *SiteMonitor) UpdateHtml_Args(filePath, ID, lastModified string) (bool, error) {
	var index = smp.GetHtml(filePath)

	if index == -1 {
		var errMsg = "SiteMonitor.UpdateHtml_Args: Cannot find html with file path " + filePath
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	smp.HtmlFiles[index].FilePath = filePath
	smp.HtmlFiles[index].ID = ID
	smp.HtmlFiles[index].LastModified = lastModified

	return true, nil
}

func (smp *SiteMonitor) UpdateHtml(htmPage HtmlPage) (bool, error) {
	var index = smp.GetHtml(htmPage.FilePath)

	if index == -1 {
		var errMsg = "SiteMonitor.UpdateHtml: Cannot find html with file path " + htmPage.FilePath
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	smp.HtmlFiles[index].FilePath = htmPage.FilePath
	smp.HtmlFiles[index].ID = htmPage.ID
	smp.HtmlFiles[index].LastModified = htmPage.LastModified

	return true, nil
}

func (smp *SiteMonitor) DeleteHtml(filePath string) (bool, error) {
	var index = smp.GetHtml(filePath)
	if index == -1 {
		var errMsg = "SiteMonitor.DeleteHtml: Cannot find html with file path " + filePath
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	smp.HtmlFiles = append(smp.HtmlFiles[:index], smp.HtmlFiles[index+1:]...)
	return true, nil
}

func (smp *SiteMonitor) AddLink_Args(linkUrl, ID, linkTitle string, isTop bool) (bool, error) {
	if smp.GetLink(linkUrl) != -1 {
		var errMsg = "SiteMonitor.AddLink_Args: Link already added"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	var linkPage LinkPage
	linkPage.Url = linkUrl
	linkPage.ID = ID
	linkPage.Title = linkTitle
	linkPage.IsTop = isTop

	smp.LinkFiles = append(smp.LinkFiles, linkPage)

	return true, nil
}

func (smp *SiteMonitor) AddLink(linkPage LinkPage) (bool, error) {
	if smp.GetLink(linkPage.Url) != -1 {
		var errMsg = "SiteMonitor.AddLink: Link already added"
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	smp.LinkFiles = append(smp.LinkFiles, linkPage)

	return true, nil
}

func (smp *SiteMonitor) GetLink(linkUrl string) int {
	if linkUrl == "" {
		return -1
	}

	for index, linkPage := range smp.LinkFiles {
		var sLink = strings.ToLower(linkPage.Url)
		var fLink = strings.ToLower(linkUrl)
		if sLink == fLink {
			return index
		}
	}

	return -1
}

func (smp *SiteMonitor) UpdateLink_Args(linkUrl, ID, linkTitle string, isTop bool) (bool, error) {
	var index = smp.GetLink(linkUrl)

	if index == -1 {
		var errMsg = "SiteMonitor.UpdateLink_Args: Cannot find link with url " + linkUrl
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	smp.LinkFiles[index].Url = linkUrl
	smp.LinkFiles[index].ID = ID
	smp.LinkFiles[index].Title = linkTitle
	smp.LinkFiles[index].IsTop = isTop

	return true, nil
}

func (smp *SiteMonitor) UpdateLink(linkPage LinkPage) (bool, error) {
	var index = smp.GetLink(linkPage.Url)

	if index == -1 {
		var errMsg = "SiteMonitor.UpdateLink: Cannot find link with url " + linkPage.Url
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	smp.LinkFiles[index].Url = linkPage.Url
	smp.LinkFiles[index].ID = linkPage.ID
	smp.LinkFiles[index].Title = linkPage.Title

	return true, nil
}

func (smp *SiteMonitor) DeleteLink(linkUrl string) (bool, error) {
	var index = smp.GetLink(linkUrl)
	if index == -1 {
		var errMsg = "SiteMonitor.DeleteLink: Cannot find link with url " + linkUrl
		Utils.Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	smp.LinkFiles = append(smp.LinkFiles[:index], smp.LinkFiles[index+1:]...)
	return true, nil
}

func (smp *SiteMonitor) GetMonitorFilesFolderPath() string {
	return filepath.Join(smp.MonitorFolderPath, "Files")
}

func (smp *SiteMonitor) GetNormalFileList() NormalFileSlice {
	var monitorFilesFolder = smp.GetMonitorFilesFolderPath()

	var filesList NormalFileSlice

	filepath.Walk(monitorFilesFolder, func(path string, info os.FileInfo, err error) error {
		var fileName = info.Name()
		var relativePath = path[len(smp.MonitorFolderPath):]
		var lastModified = info.ModTime().Format("2006-01-02 15:04:05")

		if fileName != "Files" {
			var normalFile NormalFile
			normalFile.FileName = fileName
			normalFile.FilePath = relativePath
			normalFile.LastModified = lastModified

			filesList = append(filesList, normalFile)
		}

		return nil
	})

	sort.Sort(filesList)
	return filesList
}
