# Shipper
Scripts to extract and sends files to ftps and sftp servers

## Configuration
| Field Name                | Type          | Description                                                                                                                                          |
|---------------------------|---------------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| log2console               | boolean       | Flag indicating if the logging should be duplicated to the console                                                                                   |
| timesFilePath             | string        | File Path to the Times File                                                                                                                          |
| logging                   | object        |                                                                                                                                                      |
|        script             | string        | Path to Script Logging,  which can have a dynamic timestamp, see [Dynamic Timestamp](../../README.md#dynamic-timestamp) for more information         |
|        files              | string        | Path to processed files Logging, which can have a dynamic timestamp, see [Dynamic Timestamp](../../README.md#dynamic-timestamp) for more information |
| ftps                      | array         |                                                                                                                                                      |
|     name                  | string        | Server identifier name                                                                                                                               |
|     hostname              | string        | Server Hostname                                                                                                                                      |
|     port                  | int           | Server port                                                                                                                                          |
|     username              | string        | Username to authenticate on the server                                                                                                               |
|     password              | string        | Password to authenticate on the server                                                                                                               |
|     protocol              | string        | Server trasferer protocol: ftp, sftp                                                                                                                 |
| services                  | array         |                                                                                                                                                      |
|         name              | string        | Service identifier name                                                                                                                              |
|         ftpConfig         | string array  | List of servers to run the services againts                                                                                                          |
|         sourceFolder      | string        | Source Folder                                                                                                                                        |
|         destinationFolder | string        | Destination Folder                                                                                                                                   |
|         filePrefix        | string        | File prefix to filter source files                                                                                                                   |
|         fileExtension     | string        | File extention to filter source files                                                                                                                |
|         historyFolder     | string        | Folder to move/archive process files                                                                                                                 |
|         maxTime           | int           | Max time (minutes) windows of files to process, see [Time Windows](../../README.md#time-windows) for more information                                |
|         windowLimit       | int           | Limit (minutes) in relation to NOW where newer files won't be process, see [Time Windows](../../README.md#time-windows) for more information         |
