package transport

import (
	"log"
	"os"

	"github.com/jlaffaye/ftp"
)

func checkRemoteFolder(conn *ftp.ServerConn, folderPath string) {
	err := conn.ChangeDir(folderPath)
	if err != nil {
		// folder doesn't exists, create
		log.Printf("Create remote folder %s\n", folderPath)
		conn.MakeDir(folderPath)
	}
}

func checkLocalFolder(folderPath string) {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(folderPath, 0755)
	}
}
