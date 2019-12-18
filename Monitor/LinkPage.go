package Monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"ipsd/Utils"
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
		fmt.Println(errMsg)

		return nil, errors.New(errMsg)
	}

	var strLinks []string
	errUnmarshal := json.Unmarshal(bFileContent, &strLinks)
	if errUnmarshal != nil {
		fmt.Println(errUnmarshal.Error())
		return nil, errUnmarshal
	}

	var links []LinkPage
	for _, sLink := range strLinks {
		var strVariables = strings.Split(sLink, "|")
		if len(strVariables) != 4 {
			var errMsg string
			errMsg = "ReadLinksFromFile: : Cannot get Link Information from  " + sLink
			fmt.Println(errMsg)

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
