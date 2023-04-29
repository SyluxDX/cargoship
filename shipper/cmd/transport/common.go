package transport

import (
	"fmt"
	"os"

	"cargoship/shipper/cmd/logging"

	"github.com/jlaffaye/ftp"
)

func checkRemoteFolder(conn *ftp.ServerConn, folderPath string, logger logging.Logger) {
	err := conn.ChangeDir(folderPath)

	if err != nil {
		// folder doesn't exists, create
		logger.LogInfo(fmt.Sprintf("Create remote folder %s\n", folderPath))
		conn.MakeDir(folderPath)
	}
}

func checkLocalFolder(folderPath string) {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(folderPath, 0755)
	}
}
