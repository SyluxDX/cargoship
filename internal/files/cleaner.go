package files

import (
	"fmt"
	"os"

	"cargoship/internal/configurations"
	"cargoship/internal/logging"
)

func CleanFiles(
	service configurations.ServiceConfig,
	times *configurations.FileTimes,
	scriptLogger logging.Logger,
	filesLogger logging.Logger,
) {
	scriptLogger.LogInfo(fmt.Sprintf("Processing %s: %s\n", service.Mode, service.Name))
	files, err := listLocalDirectory(service.Src, service.Prefix, service.Extension)
	if err != nil {
		scriptLogger.LogWarn(err.Error())
	}

	// get last file time
	fileTime := times.GetTimes(service.Name, service.Mode)
	// filter local files
	filterd := dateFilterLocalDirectory(files, fileTime, service.MaxTime, service.Window)

	lastFileTime := fileTime
	for _, file := range filterd {
		filesLogger.LogInfo(fmt.Sprintf("Deleted file %s", file.Name()))
		// delete file
		err := os.Remove(fmt.Sprintf("%s/%s", service.Src, file.Name()))
		if err != nil {
			scriptLogger.LogError(err.Error())
		}
		// update
		lastFileTime = file.ModTime()
	}
	if lastFileTime != fileTime {
		times.UpsertTimes(service.Name, service.Mode, lastFileTime.UTC())
	}
	scriptLogger.LogInfo(fmt.Sprintf("Deleted %d file(s)", len(filterd)))
}
