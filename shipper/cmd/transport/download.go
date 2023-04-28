package transport

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"cargoship/shipper/cmd/configurations"

	"github.com/jlaffaye/ftp"
)

func listRemoteDirectory(conn *ftp.ServerConn, source string, prefix string, extension string) ([]*ftp.Entry, error) {
	var outputList []*ftp.Entry

	entries, err := conn.List(source)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name, prefix) && strings.HasSuffix(entry.Name, extension) {
			outputList = append(outputList, entry)
		}
	}
	return outputList, nil
}

func dateFilterRemoteDirectory(entries []*ftp.Entry, lastTime time.Time, maxTime int, limit int) []*ftp.Entry {
	var outputList []*ftp.Entry

	filesLimit := time.Now().UTC().Add(time.Minute * time.Duration(limit*-1))

	for _, entry := range entries {
		if entry.Time.After(lastTime) && entry.Time.Before(filesLimit) {
			if len(outputList) == 0 {
				// update file limit with max time
				maxLimit := entry.Time.Add(time.Minute * time.Duration(maxTime))
				if maxLimit.Before(filesLimit) {
					filesLimit = maxLimit
				}
			}
			outputList = append(outputList, entry)
		}
		// cut for loop if files (entry) are after the file limit time
		if entry.Time.After(filesLimit) {
			break
		}
	}
	return outputList
}

func download(conn *ftp.ServerConn, destination string, entry *ftp.Entry) error {
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
		log.Panic(err)
	}
	log.Printf("Donwloaded file %s (size %d), written %d\n", entry.Name, entry.Size, sizeWritten)

	return nil
}

func DownloadFiles(serverName string, ftpConn *ftp.ServerConn, service configurations.ServiceConfig, times *[]configurations.FileTimes) {

	log.Printf("Processing %s: %s\n", service.Mode, service.Name)
	// check folder
	checkLocalFolder(service.Dst)
	// check remote hostory folder
	if service.History != "" {
		checkRemoteFolder(ftpConn, service.History)
	}

	// get last file time
	fileTime := configurations.GetTimes(*times, serverName, service.Mode, service.Name)
	// list files in directory
	entries, err := listRemoteDirectory(ftpConn, service.Src, service.Prefix, service.Extension)
	if err != nil {
		log.Panic(err)
	}

	entries = dateFilterRemoteDirectory(entries, fileTime, service.MaxTime, service.Window)
	// check if there are any files to download
	if len(entries) == 0 {
		log.Println("No files to download")
		return
	}
	// donwload files
	lastFileTime := fileTime
	err = ftpConn.ChangeDir(service.Src)
	if err != nil {
		log.Printf("[ERROR] %s\n", err)
	}
	for _, entry := range entries {
		err := download(ftpConn, fmt.Sprintf("%s/%s", service.Dst, entry.Name), entry)
		if err != nil {
			log.Printf("[ERROR] %s\n", err)
			break
		}
		if service.History != "" {
			err := ftpConn.Rename(entry.Name, fmt.Sprintf("%s/%s", service.History, entry.Name))
			if err != nil {
				log.Printf("[ERROR] %s\n", err)
			}
		}
		// update
		lastFileTime = entry.Time
	}
	// update last downloaded time
	if lastFileTime != fileTime {
		configurations.UpsertTimes(times, serverName, service.Mode, service.Name, lastFileTime)
	}
}
