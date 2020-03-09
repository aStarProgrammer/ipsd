package Monitor

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"ipsd/Utils"
	"os"
	"strconv"
	"strings"
)

type LinkPage struct {
	Url   string
	ID    string
	Title string
	IsTop bool
}

func LinkPage2String(linkPage LinkPage) string {
	var str string
	str = linkPage.Url
	str = str + "|"
	str = str + linkPage.ID
	str = str + "|"
	str = str + linkPage.Title
	str = str + "|"
	str = str + strconv.FormatBool(linkPage.IsTop)
	return str
}

func ReadLinksFromFile(filePath string) ([]LinkPage, error) {
	if Utils.PathIsExist(filePath) == false {
		return nil, errors.New("ReadLinksFromFile: FilePath not exist " + filePath)
	}

	bFileContent, errReadFile := ioutil.ReadFile(filePath)

	if errReadFile != nil {
		var errMsg string
		errMsg = "ReadLinksFromFile: : Cannot read Markdown file " + filePath
		Utils.Logger.Println(errMsg)

		return nil, errors.New(errMsg)
	}

	var strLinks []string
	errUnmarshal := json.Unmarshal(bFileContent, &strLinks)
	if errUnmarshal != nil {
		Utils.Logger.Println(errUnmarshal.Error())
		return nil, errUnmarshal
	}

	var links []LinkPage
	for _, sLink := range strLinks {
		var strVariables = strings.Split(sLink, "|")
		if len(strVariables) != 4 {
			var errMsg string
			errMsg = "ReadLinksFromFile: : Cannot get Link Information from  " + sLink
			Utils.Logger.Println(errMsg)

			return nil, errors.New(errMsg)
		}

		var link LinkPage
		link.Url = strVariables[0]
		link.ID = strVariables[1]
		link.Title = strVariables[2]
		bIsTop, errParseBool := strconv.ParseBool(strVariables[3])

		if errParseBool == nil {
			link.IsTop = bIsTop
		} else {
			link.IsTop = false
		}

		links = append(links, link)
	}

	return links, nil
}

func FindLink(sLink LinkPage, links []LinkPage) bool {
	for _, fLink := range links {
		if sLink.Url == fLink.Url && sLink.ID == fLink.ID && sLink.IsTop == fLink.IsTop && sLink.Title == fLink.Title {
			return true
		}
	}

	return false
}

func Links2Json(links []LinkPage) (string, error) {
	if links == nil {
		return "", errors.New("LinksPage.Links2Json: links is nil")
	}

	if len(links) == 0 {
		return "", errors.New("LinksPage.Links2Json: links is empty")
	}

	var strLinks []string

	/*
		LinkPage1 LinkPage2 ==> "https://www.baidu.com|cbfcda20e936ca4e07eeb5b75960163e|baidu.com|true","https://www.google.com|cbfcda20e936ca4e07eeb5b75960163e|google|false"

	*/
	for _, link := range links {
		linkStr := LinkPage2String(link)
		if "" != linkStr {
			strLinks = append(strLinks, linkStr)
		}
	}

	// "https://www.baidu.com|cbfcda20e936ca4e07eeb5b75960163e|baidu.com|true","https://www.google.com|cbfcda20e936ca4e07eeb5b75960163e|google|false"
	//====>
	// ["https://www.baidu.com|cbfcda20e936ca4e07eeb5b75960163e|baidu.com|true","https://www.google.com|cbfcda20e936ca4e07eeb5b75960163e|google|false"]

	var _jsonbyte []byte
	var errJson error
	if len(strLinks) != 0 {
		_jsonbyte, errJson = json.Marshal(strLinks)
	}

	return string(_jsonbyte), errJson
}

func SaveLinksToFile(filePath string, links []LinkPage) (bool, error) {
	if "" == filePath {
		return false, errors.New("SaveLinksToFile:FilePath is empty")
	}

	linkStr, errConvert := Links2Json(links)

	if errConvert != nil {
		return false, errors.New("LinkPage.SaveLinksToFile:cannot convert link Pages to json string")
	}

	var errFilePath error
	if !Utils.PathIsExist(filePath) {
		filePath, errFilePath = Utils.MakePath(filePath)
		if errFilePath != nil {
			return false, errors.New("LinkPage.SaveLinksToFile:Path nor exist and create parent folder failed")
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

	_, errWriteFile := fp.WriteString(linkStr)

	if errWriteFile != nil {
		return false, errors.New("LInkPage,SaveLinksToFile:Write json to file failed")
	}

	return true, nil
}
