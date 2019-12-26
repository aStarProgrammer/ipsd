package Utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/shamsher31/goimgtype"
)

//
func PathIsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GUID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func CurrentTime() string {
	t := time.Now()
	str := t.Format("2006-01-02 15:04:05")

	return str
}

func PathIsMarkdown(filePath string) bool {
	if filePath == "" {
		return false
	}

	ext := filepath.Ext(filePath)

	if ext == ".md" || ext == ".markdown" || ext == ".mdown" || ext == ".mmd" {
		return true
	}
	return false
}

func MakeFolder(sPath string) (bool, error) {
	if sPath == "" {
		return false, errors.New("MakeFolder: sPath is empty")
	}

	sFolderPath, errFolderPath := MakePath(sPath)

	if errFolderPath != nil {
		return false, errFolderPath
	}

	errFolderPath = os.Mkdir(sFolderPath, os.ModePerm)

	if errFolderPath != nil {
		return false, errFolderPath
	}

	return true, nil
}

func SaveBase64AsImage(imageContent, targetPath string) (bool, error) {
	if imageContent == "" {
		return false, errors.New("SaveBase64AsImage : image content is empty")
	}

	if targetPath == "" {
		return false, errors.New("SaveBase64AsImage : target file path is empty")
	}

	if PathIsExist(targetPath) {
		bDelete := DeleteFile(targetPath)
		if bDelete == false {
			return false, errors.New("SaveBase64AsImage : target Path already exist and cannot delete")
		}
	}

	if strings.Contains(imageContent, "data:") == false || strings.Contains(imageContent, ";base64,") == false {
		return false, errors.New("SaveBase64AsImage : Image Content Format Error")
	}

	var base64Index = strings.Index(imageContent, ";base64,")
	var base64Image = imageContent[base64Index+8:]

	decodedImage, errDecode := base64.StdEncoding.DecodeString(base64Image)
	if errDecode != nil {
		return false, errors.New("SaveBase64AsImage : Cannot Decode Base64 Image")
	}
	err2 := ioutil.WriteFile(targetPath, decodedImage, 0666)

	if err2 != nil {
		return false, errors.New("SaveBase64AsImage : Cannot Save image")
	}

	return true, nil
}

func ReadImageAsBase64(imagePath string) (string, error) {

	var retImage string
	retImage = ""

	image, errRead := ioutil.ReadFile(imagePath)

	if errRead != nil {
		return "", errors.New("ReadImageAsBase64: Read Fail")
	}

	imageBase64 := base64.StdEncoding.EncodeToString(image)

	datatype, err2 := imgtype.Get(imagePath)
	if err2 != nil {
		return "", errors.New("ReadImageAsBase64: Cannot get image type")
	} else {
		retImage = "data:" + datatype + ";base64," + imageBase64
	}

	return retImage, nil
}

func PathIsImage(filePath string) bool {

	if filePath == "" {
		return false
	}

	_, err2 := imgtype.Get(filePath)
	if err2 != nil {
		return false
	}
	return true
}

func GetImageType(base64Image string) (string, error) {
	if base64Image == "" {
		return "", errors.New("Get Image Type: base64Image is empty")
	}

	var datatypeParts = strings.Split(base64Image, ";") //Get data:image/png
	if len(datatypeParts) > 1 {
		var datatypePart = datatypeParts[0]
		var datatypes = strings.Split(datatypePart, ":") //Get image/png
		if len(datatypes) == 2 {
			var datatype = datatypes[1]
			var subTypes = strings.Split(datatype, "/") //Get png
			if len(subTypes) == 2 {
				return subTypes[1], nil
			} else {
				return "", errors.New("Get Image Type : Cannot get image type")
			}
		} else {
			return "", errors.New("Get Image Type : Cannot get image type")
		}
	}

	return "", errors.New("Get Image Type : Cannot get image type")
}

func MakePath(sPath string) (string, error) {
	sfolder, sfile := filepath.Split(sPath)

	if sfolder == "" || sfile == "" {
		return "", errors.New("MakePath: folder or file name is empty folder " + sfolder + " file " + sfile)
	}

	sfolder = filepath.Clean(sfolder)

	if !PathIsExist(sfolder) {
		os.MkdirAll(sfolder, os.ModePerm)
	}

	return filepath.Join(sfolder, sfile), nil

}

func MakeSoftLink4Folder(srcFolder, linkFolder string) (bool, error) {
	srcExist := PathIsExist(srcFolder)

	if !srcExist {
		var errMsg = "Make Soft Link 4 Folder: SrcFolder Not Exist " + srcFolder
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	targetExist := PathIsExist(linkFolder)

	if targetExist {
		var errMsg = "Make Soft Link 4 Folder:linkFolder Already Exist " + linkFolder
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	errLink := os.Symlink(srcFolder, linkFolder)

	if errLink != nil {
		Logger.Println("MakeSoftLink: " + errLink.Error())
		return false, errLink
	}
	return true, nil
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		var errMsg = "CopyFile " + src + "is not a regular file"
		return 0, errors.New(errMsg)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func MoveFile(src, dst string) (int64, error) {
	iCopy, errCopy := CopyFile(src, dst)

	if errCopy != nil {
		return 0, errCopy
	}

	errRemove := os.Remove(src)

	if errRemove != nil {
		return 0, errRemove
	}

	return iCopy, nil
}

func DeleteFile(filePath string) bool {
	errRemove := os.Remove(filePath)

	if errRemove != nil {
		return false
	}
	return true
}

func CreateFile(filePath string) bool {
	if PathIsExist(filePath) == false {
		filePath, _ = MakePath(filePath)
	}
	file2, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		Logger.Println(err.Error())
		return false
	}
	file2.Close()

	return true
}

func GetImageWithSameName(filePath string) (string, error) {
	if PathIsExist(filePath) == false {
		return "", errors.New("GetImageWithSameName: filePath not exist " + filePath)
	}

	sFolder, sFile := filepath.Split(filePath)
	sExt := filepath.Ext(filePath)
	sShortName := strings.Replace(sFile, sExt, "", -1)

	files, errReadDir := ioutil.ReadDir(sFolder)

	if errReadDir != nil {
		return "", errReadDir
	}

	for _, file := range files {
		var fName = file.Name()
		var fPath = filepath.Join(sFolder, fName)
		var fExt = filepath.Ext(fPath)
		if PathIsImage(fPath) {
			fShortName := strings.Replace(fName, fExt, "", -1)
			if sShortName == fShortName {
				return fPath, nil
			}
		}
	}

	return "", errors.New("GetImageWithSameName: not find image with same name of " + filePath)
}

func GetImageWithSameTitle(fileFolder, fileTitle string) (string, error) {
	if PathIsExist(fileFolder) == false {
		return "", errors.New("GetImageWithSameTitle: File Folder not exist " + fileFolder)
	}

	files, errReadDir := ioutil.ReadDir(fileFolder)

	if errReadDir != nil {
		return "", errReadDir
	}

	for _, file := range files {
		var fName = file.Name()
		var fPath = filepath.Join(fileFolder, fName)
		var fExt = filepath.Ext(fPath)
		if PathIsImage(fPath) {
			fShortName := strings.Replace(fName, fExt, "", -1)
			if fileTitle == fShortName {
				return fPath, nil
			}
		}
	}

	return "", errors.New("GetImageWithSameTitle: not find image with same name of " + fileTitle)
}

func CurrentUser() string {
	u, err := user.Current()
	if err != nil {
		return "Chao"
	}

	return u.Username
}

const MAXTITLEIMAGESIZE int64 = 30720

func ImageTooBig(titleImagePath string) (bool, error) {
	fileInfoTitleImage, errFileInfoTitleImage := os.Stat(titleImagePath)

	if errFileInfoTitleImage != nil {
		var errMsg = "Utils.ImageTooBig: Cannot get file size of titleImage"
		Logger.Println(errMsg)
		return false, errors.New(errMsg)
	}

	titleImageSize := fileInfoTitleImage.Size()

	if titleImageSize > MAXTITLEIMAGESIZE {
		return true, nil
	}

	return false, nil
}

var Logger *log.Logger

func InitLogger() {
	file, err := os.Create("ipsd.log")
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}

	Logger = log.New(file, "", log.Llongfile)
}

func PathIsFile(filePath string) bool {
	if filePath == "" {
		return false
	}

	if PathIsExist(filePath) == false {
		return false
	}
	f, err := os.Stat(filePath)

	if err != nil {
		Logger.Println(err.Error())
		return false
	}

	if f.IsDir() == false {
		return true
	}

	return false
}

func PathIsDir(filePath string) bool {

	if filePath == "" {
		return false
	}

	if PathIsExist(filePath) == false {
		return false
	}

	f, err := os.Stat(filePath)

	if err != nil {
		Logger.Println(err.Error())
		return false
	}

	return f.IsDir()
}
