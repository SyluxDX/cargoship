package transport

import (
	"fmt"
	"io"
	"log"
	"os"

	"cargoship/internal/configurations"
	"cargoship/internal/files"
	"cargoship/internal/logging"

	"github.com/jlaffaye/ftp"
)

func download(conn *ftp.ServerConn, destination string, entry *ftp.Entry, logger logging.Logger) error {
	remoteReader, err := conn.Retr(entry.Name)
	if err != nil {
		return err
	}
	defer remoteReader.Close()

	// create local writer
	localWriter, err := os.OpenFile(
		destination,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return err
	}
	defer localWriter.Close()

	sizeWritten, err := io.Copy(localWriter, remoteReader)
	if err != nil {
		return err
	}
	logger.LogInfo(fmt.Sprintf("Donwloaded file %s (size %d), written %d\n", entry.Name, entry.Size, sizeWritten))

	return nil
}

func DownloadFiles(
	serverName string,
	ftpConn *ftp.ServerConn,
	service configurations.ShipperService,
	times *[]configurations.ShipperFileTimes,
	scriptLogger logging.Logger,
	filesLogger logging.Logger,
) {

	// check folder
	files.CheckLocalFolder(service.Dst)
	// check remote hostory folder
	if service.History != "" {
		checkRemoteFolder(ftpConn, service.History, scriptLogger)
	}

	// get last file time
	fileTime := configurations.ShipperGetTimes(*times, serverName, service.Mode, service.Name)
	// list files in directory
	entries, err := listRemoteDirectory(ftpConn, service.Src, service.Prefix, service.Extension)
	if err != nil {
		log.Panic(err)
	}

	entries = dateFilterRemoteDirectory(entries, fileTime, service.MaxTime, service.Window)
	// check if there are any files to download
	if len(entries) == 0 {
		scriptLogger.LogInfo("No files to download")
		return
	}
	// donwload files
	lastFileTime := fileTime
	err = ftpConn.ChangeDir(service.Src)
	if err != nil {
		scriptLogger.LogWarn(err.Error())
	}
	for _, entry := range entries {
		err := download(ftpConn, fmt.Sprintf("%s/%s", service.Dst, entry.Name), entry, filesLogger)
		if err != nil {
			scriptLogger.LogWarn(err.Error())
			break
		}
		if service.History != "" {
			err := ftpConn.Rename(entry.Name, fmt.Sprintf("%s/%s", service.History, entry.Name))
			if err != nil {
				scriptLogger.LogWarn(err.Error())
			}
			scriptLogger.LogInfo(fmt.Sprintf("Moved file %s to history folder %s\n", entry.Name, service.History))
		}
		// update
		lastFileTime = entry.Time
	}
	// update last downloaded time
	if lastFileTime != fileTime {
		configurations.ShipperUpsertTimes(times, serverName, service.Mode, service.Name, lastFileTime)
	}
}
