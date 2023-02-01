# Loader
Scripts to compress files and clean local folders

## Configuration
| Field Name                | Type          | Description                                                                                                                     |
|---------------------------|---------------|---------------------------------------------------------------------------------------------------------------------------------|
| log2console               | boolean       | Flag indicating if the logging should be duplicated to the console                                                              |
| timesFilePath             | string        | File Path to the Times File                                                                                                     |
| logging                   | object        |                                                                                                                                 |
|        folder             | string        | Logging Folder                                                                                                                  |
|        filename           | string        | Logging filename which can have a dynamic timestamp, see [Dynamic Timestamp](#dynamic-timestamp) for more information           |
| services                  | array         |                                                                                                                                 |
|         name              | string        | Service identifier name                                                                                                         |
|         sourceFolder      | string        | Source Folder                                                                                                                   |
|         destinationFolder | string        | Destination Folder                                                                                                              |
|         filePrefix        | string        | File prefix to filter source files                                                                                              |
|         fileExtension     | string        | File extention to filter source files                                                                                           |
|         maxTime           | int           | Max time (minutes) windows of files to process, see [Time Windows](#time-windows) for more information                          |
|         windowLimit       | int           | Limit (minutes) in relation to NOW where newer files won't be process, see [Time Windows](#time-windows) for more information   |
