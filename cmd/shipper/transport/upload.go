package transport

import (
	"fmt"
	"os"

	"cargoship/internal/configurations"
	"cargoship/internal/files"
	"cargoship/internal/logging"

	"github.com/jlaffaye/ftp"
)

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
	service configurations.ShipperService,
	times *[]configurations.ShipperFileTimes,
	scriptLogger logging.Logger,
	filesLogger logging.Logger,
) {

	scriptLogger.LogInfo(fmt.Sprintf("Processing %s: %s\n", service.Mode, service.Name))
	// check folders
	checkRemoteFolder(ftpConn, service.Dst, scriptLogger)
	if service.History != "" {
		files.CheckLocalFolder(service.History)
	}
	// get last file time
	fileTime := configurations.ShipperGetTimes(*times, serverName, service.Mode, service.Name)

	// list files in directory
	entries, err := files.ListLocalDirectory(service.Src, service.Prefix, service.Extension)
	if err != nil {
		scriptLogger.LogWarn(err.Error())
	}

	entries = files.DateFilterLocalDirectory(entries, fileTime, service.MaxTime, service.Window)
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
		configurations.ShipperUpsertTimes(times, serverName, service.Mode, service.Name, lastFileTime)
	}
}
