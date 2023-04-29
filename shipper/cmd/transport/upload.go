package transport

import (
	"fmt"
	"os"
	"strings"
	"time"

	"cargoship/shipper/cmd/configurations"
	"cargoship/shipper/cmd/logging"

	"github.com/jlaffaye/ftp"
)

func listLocalDirectory(source string, prefix string, extension string) ([]os.FileInfo, error) {
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

func dateFilterLocalDirectory(entries []os.FileInfo, lastTime time.Time, maxTime int, limit int) []os.FileInfo {
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

func upload(conn *ftp.ServerConn, source string, entry os.FileInfo, logger logging.Logger) error {
	// local reader
	localReader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer localReader.Close()

	// upload
	err = conn.Stor(entry.Name(), localReader)
	if err != nil {
		return err
	}
	remoteSize, _ := conn.FileSize(entry.Name())
	logger.LogInfo(fmt.Sprintf("Uploaded file %s (size %d), written %d\n", entry.Name(), entry.Size(), remoteSize))

	return nil
}

func UploadFiles(
	serverName string,
	ftpConn *ftp.ServerConn,
	service configurations.ServiceConfig,
	times *[]configurations.FileTimes,
	scriptLogger logging.Logger,
	filesLogger logging.Logger,
) {

	scriptLogger.LogInfo(fmt.Sprintf("Processing %s: %s\n", service.Mode, service.Name))
	// check folders
	checkRemoteFolder(ftpConn, service.Dst, scriptLogger)
	if service.History != "" {
		checkLocalFolder(service.History)
	}
	// get last file time
	fileTime := configurations.GetTimes(*times, serverName, service.Mode, service.Name)

	// list files in directory
	entries, err := listLocalDirectory(service.Src, service.Prefix, service.Extension)
	if err != nil {
		scriptLogger.LogWarn(err.Error())
	}

	entries = dateFilterLocalDirectory(entries, fileTime, service.MaxTime, service.Window)
	// check if there are any files to upload
	if len(entries) == 0 {
		scriptLogger.LogInfo("No files to upload")
		return
	}
	// upload files
	lastFileTime := fileTime
	err = ftpConn.ChangeDir(service.Dst)
	if err != nil {
		scriptLogger.LogWarn(err.Error())
	}
	for _, entry := range entries {
		err := upload(ftpConn, fmt.Sprintf("%s/%s", service.Src, entry.Name()), entry, filesLogger)
		if err != nil {
			scriptLogger.LogWarn(err.Error())
			break
		}
		if service.History != "" {
			os.Rename(
				fmt.Sprintf("%s/%s", service.Src, entry.Name()),
				fmt.Sprintf("%s/%s", service.History, entry.Name()),
			)
			scriptLogger.LogInfo(fmt.Sprintf("Moved file %s to history folder %s\n", entry.Name(), service.History))
		}
		// update
		lastFileTime = entry.ModTime()
	}
	if lastFileTime != fileTime {
		configurations.UpsertTimes(times, serverName, service.Mode, service.Name, lastFileTime)
	}
}
