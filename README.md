# Cargoship [In Contruction]
Scripts to extract, compress, sends and cleanup files.

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

### Planning
- Shipper:
    - Write sender, ftp and sftp
- Loader:
    - Compressor
    - Cleaner
- Packager/Processor:
    - Should be in here?
    - Use external awk?
    - Does golang have internal awk?

### ToDo
- Test golang connections to SFTP server
    - Create sftp Extractor script
- Add logging to file:
    - create files if config exists
    - configure loggers: import, export, console
- See how to use goroutines to speed up downloads
- Update README with futher description of Time Windows
- Create Makefile to build scripts to windows and linux
