package Page

type NormalFile struct {
	FileName     string
	FilePath     string
	LastModified string
}

func GetNormalFile(files []NormalFile, filePath string) int {
	if nil == files {
		return -1
	}

	if filePath == "" {
		return -1
	}

	if len(files) == 0 {
		return -1
	}

	for index, file := range files {
		if file.FilePath == filePath {
			return index
		}
	}

	return -1
}
