package Monitor

import (
	"errors"
	"io/ioutil"
	"ipsd/Utils"
	"strings"
)

type HtmlPage struct {
	FilePath     string
	ID           string
	LastModified string
}

func ReadHtmlProperties(filePath string) (*HtmlProperties, bool, error) {
	if Utils.PathIsExist(filePath) == false {
		var errMsg = "Monitor.ReadHtmlPageProperties: Html File not exist " + filePath
		Utils.Logger.Println(errMsg)
		return nil, false, errors.New(errMsg)
	}

	bFileContent, errReadFile := ioutil.ReadFile(filePath)

	if errReadFile != nil {
		var errMsg string
		errMsg = "Monitor.ReadHtmlPageProperties: Cannot read Html file " + filePath
		Utils.Logger.Println(errMsg)

		return nil, false, errors.New(errMsg)
	}

	fileContent := string(bFileContent)

	if strings.HasPrefix(fileContent, "<!--") == false || strings.Contains(fileContent, "[//]: #") == false {
		var errMsg = "You should edit html file before add it to the folder, for more information, read the ReadMe"
		Utils.Logger.Println(errMsg)
		return nil, false, errors.New(errMsg)
	}

	metadataInfo := fileContent[:strings.Index(fileContent, "-->")]
	metadataInfo = metadataInfo[strings.Index(metadataInfo, "[//]: #"):]

	var tvalues = strings.Split(metadataInfo, "[//]: #")

	var values []string

	for _, v := range tvalues {
		var t = strings.TrimSpace(v)
		if t != "" {
			values = append(values, t)
		}
	}

	if len(values) != 4 {
		var errMsg string
		errMsg = "Monitor.ReadHtmlPageProperties: Cannot read metadata properties of Html file " + filePath
		Utils.Logger.Println(errMsg)

		return nil, false, errors.New(errMsg)
	}

	var title, author, description, vIsTop string
	var isTop bool
	title = strings.TrimSpace(values[0])
	title = values[0][strings.Index(title, ":")+1:]
	title = title[:len(title)-1]
	title = strings.TrimSpace(title)

	author = strings.TrimSpace(values[1])
	author = values[1][strings.Index(author, ":")+1:]
	author = author[:len(author)-1]
	author = strings.TrimSpace(author)

	description = strings.TrimSpace(values[2])
	description = values[2][strings.Index(description, ":")+1:]
	description = description[:len(description)-1]
	description = strings.TrimSpace(description)

	vIsTop = strings.TrimSpace(values[3])
	vIsTop = values[3][strings.Index(vIsTop, ":")+1:]
	vIsTop = vIsTop[:len(vIsTop)-1]
	vIsTop = strings.TrimSpace(vIsTop)
	if strings.Contains(vIsTop, "true") {
		isTop = true
	} else {
		isTop = false
	}

	var htmProperties HtmlProperties
	var htmP *HtmlProperties
	htmP = &htmProperties

	htmP.Title = title
	htmP.Author = author
	htmP.Description = description
	htmP.IsTop = isTop

	return htmP, true, nil
}

type HtmlProperties struct {
	Title       string
	Author      string
	Description string
	IsTop       bool
}
