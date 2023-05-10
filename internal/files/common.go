package files

import (
	"os"
	"strings"
	"time"
)

func CheckLocalFolder(folderPath string) {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(folderPath, 0755)
	}
}

func ListLocalDirectory(source string, prefix string, extension string) ([]os.FileInfo, error) {
	var outputList []os.FileInfo

	entries, err := os.ReadDir(source)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			if strings.HasPrefix(entry.Name(), prefix) && strings.HasSuffix(entry.Name(), extension) {
				info, err := entry.Info()
				if err != nil {
					return nil, err
				}
				outputList = append(outputList, info)
			}
		}
	}
	return outputList, nil
}

func DateFilterLocalDirectory(entries []os.FileInfo, lastTime time.Time, maxTime int, limit int) []os.FileInfo {
	var outputList []os.FileInfo

	filesLimit := time.Now().UTC().Add(time.Minute * time.Duration(limit*-1))
	for _, entry := range entries {
		if entry.ModTime().After(lastTime) && entry.ModTime().Before(filesLimit) {
			if len(outputList) == 0 {
				// update file limit with max time
				maxLimit := entry.ModTime().Add(time.Minute * time.Duration(maxTime))
				if maxLimit.Before(filesLimit) {
					filesLimit = maxLimit
				}
			}
			outputList = append(outputList, entry)
		}
		// cut for loop if files (entry) are after the file limit time
		if entry.ModTime().After(filesLimit) {
			break
		}
	}
	return outputList
}
