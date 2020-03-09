package Page

import (
	"encoding/json"
	"errors"
	"ipsd/Utils"
	"strconv"
)

type PageOutputFile struct {
	ID          string
	Title       string
	Description string
	TitleImage  string
	Type        string
	Author      string
	CreateTime  string
	IsTop       bool
	FilePath    string
}

func NewPageOutputFile() PageOutputFile {
	var po PageOutputFile
	po.ID = Utils.GUID()
	po.CreateTime = Utils.CurrentTime()
	return po
}

func NewPageOutputFileP() *PageOutputFile {
	var po PageOutputFile
	var pop *PageOutputFile
	pop = &po

	pop.ID = Utils.GUID()
	pop.CreateTime = Utils.CurrentTime()

	return pop
}

func ResetPageOutputFile(po PageOutputFile) {
	po.ID = ""
}

func IsPageOutputFileEmpty(po PageOutputFile) bool {
	if po.ID == "" {
		return true
	}
	return false
}

func (po *PageOutputFile) ToJson() (string, error) {
	var _jsonbyte []byte

	if po == nil {
		var errMsg = "PageOutputFile->ToJson: Pointer po is nil"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	if IsPageOutputFileEmpty(*po) {
		var errMsg = "PageOutputFile->ToJson:Page Output File is empty"
		Utils.Logger.Println(errMsg)
		return "", errors.New(errMsg)
	}

	_jsonbyte, err := json.Marshal(*po)

	return string(_jsonbyte), err
}

func (po *PageOutputFile) ToString() string {
	var properties string
	properties = "ID: " + po.ID
	properties += "|Title: " + po.Title
	properties += "|Author: " + po.Author
	properties += "|Type: " + po.Type
	properties += "|CreateTime: " + po.CreateTime
	properties += "|IsTop: " + strconv.FormatBool(po.IsTop)
	properties += "|FilePath: " + po.FilePath

	return properties
}

type PageOutputFileSlice []PageOutputFile

func (s PageOutputFileSlice) Len() int {
	return len(s)
}

func (s PageOutputFileSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s PageOutputFileSlice) Less(i, j int) bool {
	return s[i].CreateTime < s[j].CreateTime
}
