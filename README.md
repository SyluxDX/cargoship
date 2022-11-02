# Cargoship [In Contruction]
scripts to extract and sends files to ftps and sftp servers

## Configuration
| Field Name                | Type          | Description                                                                                                           |
|---------------------------|---------------|-----------------------------------------------------------------------------------------------------------------------|
| log2console               | boolean       | Flag indicating if the logging should be duplicated to the console                                                    |
| timesFilePath             | string        | File Path to the Times File                                                                                           |
| logging                   | object        |                                                                                                                       |
|        folder             | string        | Logging Folder                                                                                                        |
|        filename           | string        | Logging filename which can have a dynamic timestamp, see [Dynamic Timestamp](#dynamic-timestamp) for more information |
| ftps                      | array         |                                                                                                                       |
|     name                  | string        | Server identifier name                                                                                                |
|     hostname              | string        | Server Hostname                                                                                                       |
|     port                  | int           | Server port                                                                                                           |
|     username              | string        | Username to authenticate on the server                                                                                |
|     password              | string        | Password to authenticate on the server                                                                                |
|     protocol              | string        | Server trasferer protocol: ftp, sftp                                                                                  |
| services                  | array         |                                                                                                                       |
|         name              | string        | Service identifier name                                                                                               |
|         ftpConfig         | string array  | List of servers to run the services againts                                                                           |
|         sourceFolder      | string        | Source Folder                                                                                                         |
|         destinationFolder | string        | Destination Folder                                                                                                    |
|         filePrefix        | string        | File prefix to filter source files                                                                                    |
|         fileExtension     | string        | File extention to filter source files                                                                                 |
|         historyFolder     | string        | Folder to move/archive process files                                                                                  |
|         maxTime           | int           | Max time windows of files to process, see [Time Windows](#time-windows) for more information                          |
|         windowLimit       | int           | Limit in relation to NOW where newer files won't be process, see [Time Windows](#time-windows) for more information   |

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

### ToDo
- Create Makefile to build scripts to windows and linux
- Makefile to clean ftdata and run server?
- Test golang connections to SFTP server
- Create Extractor script (ftp done, code sftp)
- Create Sender Script
- Ponder if extractor and send should be one or two scripts
- See how to use goroutines to speed up downloads
