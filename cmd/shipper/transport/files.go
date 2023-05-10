package transport

import (
	"fmt"
	"strings"
	"time"

	"cargoship/internal/logging"

	"github.com/jlaffaye/ftp"
)

//// remote folder/files operation

func checkRemoteFolder(conn *ftp.ServerConn, folderPath string, logger logging.Logger) {
	err := conn.ChangeDir(folderPath)

	if err != nil {
		// folder doesn't exists, create
		logger.LogInfo(fmt.Sprintf("Create remote folder %s\n", folderPath))
		conn.MakeDir(folderPath)
	}
}

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
