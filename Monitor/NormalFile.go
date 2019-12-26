package Monitor

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

func AddNormalFile(files []NormalFile, file NormalFile) []NormalFile {
	files = append(files, file)
	return files
}

func DeleteNormalFile(files []NormalFile, file NormalFile) []NormalFile {

	var iFind = GetNormalFile(files, file.FilePath)

	if iFind != -1 {
		files = append(files[:iFind], files[iFind+1:]...)
	}
	return files
}

func UpdateNormalFile(files []NormalFile, file NormalFile) []NormalFile {
	var iFind = GetNormalFile(files, file.FilePath)

	if iFind != -1 {
		files[iFind].FileName = file.FileName
		files[iFind].FilePath = file.FilePath
		files[iFind].LastModified = file.LastModified
	}
	return files
}

type NormalFileSlice []NormalFile

func (s NormalFileSlice) Len() int {
	return len(s)
}

func (s NormalFileSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s NormalFileSlice) Less(i, j int) bool {
	return s[i].FilePath < s[j].FilePath
}
