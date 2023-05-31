# Cargoship [In Contruction]
Scripts to extract, compress, sends and cleanup files.

# Scripts/Modules
## Shipper
Send and Retries files from a ftp and sftp server.

Configurable for each server:
- folders
- files prefix and extensions
- time window to retrieve/send

For more details check [module documentation](cmd/shipper/README.md)

## Loader
Clean and compress local files for sending

Configurable "services":
- folders
- files prefix and extensions
- time window to retrieve/send

For more details check [module documentation](cmd/loader/README.md)

## Packager/Processor (ToDo?)
Should be in here?

Use external awk?

Does golang have internal awk?

## Dynamic Timestamp

### Golang Timestamp Format
| Time Part | Value |
|-----------|-------|
| year      | 2006  |
| month     | 01    |
| day       | 02    |
| hours     | 15    |
| minutes   | 04    |
| seconds   | 05    |

## Time Windows

### maxTime

_template_

Time limit calculating by using the first valid file to download and add minutes equal to maxTime value

_add more info_

### windowLimit

_template_

Time limit calculated by substratcing minutes equal to windowLimit value to current date

_add more info_

### ToDo
- Test golang connections to SFTP server
    - Create sftp Extractor script
- See how to use goroutines to speed up downloads
- Update README with futher description of Time Windows
- Create Makefile to build scripts to windows and linux
- Think if want/need to change project structure

### golang possible structure
#### /cmd

This folder contains the main application entry point files for the project, with the directory name matching the name for the binary. So for example `cmd/simple-service` meaning that the binary we publish will be `simple-service`.

#### /internal

This package holds the private library code used in your service, it is specific to the function of the service and not shared with other services. One thing to note is this privacy is enforced by the compiler itself, see the Go 1.4 release notes for more details.

#### /pkg

This folder contains code which is OK for other services to consume, this may include API clients, or utility functions which may be handy for other projects but don’t justify their own project. Personally I prefer to use this over `internal`, mainly as I like to keep things open for reuse in most of projects.

```
cmd/
    server/
        main.go
    cli/
        main.go
```
