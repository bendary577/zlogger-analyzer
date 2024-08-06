package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "io/ioutil"
    "github.com/gin-gonic/gin"
)

type scanResult struct {
    Title             string `json:"title"`
    Type              int    `json:"type"` //TODO : use enum (exception, info, searchTerm)
    NumberOfOccurance int    `json:"numberOfOccurance"`
}

func check(e error) {
    if e != nil {
        log.Fatal(e)
        panic(e)
    }
}

func newScanResult(Type int, Title string) *scanResult {
    s := scanResult{Title: Title}
    s.Type = Type
    s.NumberOfOccurance = 0
    return &s
}

func processLogFile(sptr *scanResult, file *os.File, searchTerm string) {
    scanner := bufio.NewScanner(file)
    lineNumber := 0

    for scanner.Scan() {
        lineNumber++
        line := scanner.Text()
        if strings.Contains(line, searchTerm) {
            sptr.NumberOfOccurance++
            fmt.Printf("Found '%s' in line %d: %s\n", searchTerm, lineNumber, line)
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func LogAnalyzerHandler(c *gin.Context) {
    //TODO : replace with logging
    fmt.Println("started zlog analyzer!")

    // Initialize the scanResultPtr
    var scanResultPtr *scanResult

    //TODO : remove log files location to a configured value
    folderPath := "./logFiles"
    fileCount := 0
    /*
       how the zlogger analyzer get the logs data that will get analyzed
       0 : read logs from path - default mode
       1 : read logs from client app call
       2 : read logs from a webhook event on configured log path
       //TODO : use enums
    */
    zlogAnalyzerProcessingMode := 0

    if zlogAnalyzerProcessingMode == 0 {
        dir, err := os.Open(folderPath)
        //TODO : better exception handling
        check(err)
        defer dir.Close()

        files, err := ioutil.ReadDir(folderPath)
        check(err)

        // Check number of files
        for _, file := range files {
            if !file.IsDir() {
                fileCount++
            }
        }

        fmt.Printf("Total number of files in the folder: %d\n", fileCount)

        if fileCount > 0 {
            // Initialize the scanResultPtr
            scanResultPtr = newScanResult(0, "Error Result")
            for _, file := range files {
                if !file.IsDir() {
                    filePointer, err := os.Open(filepath.Join(folderPath, file.Name()))
                    check(err)
                    processLogFile(scanResultPtr, filePointer, "ERR")
                    filePointer.Close()
                }
            }
            fmt.Printf("Found '%d' number of occurrences in the scan result\n", scanResultPtr.NumberOfOccurance)
        } else {
            fmt.Println("No files found in the directory.")
        }
    } else {
        fmt.Println("Invalid processing mode.")
    }

    // Check if scanResultPtr is nil
    if scanResultPtr == nil {
        // TODO : handle
        fmt.Println("Not able to scan and analyze the file")
        c.JSON(500, gin.H{"error": "Not able to scan and analyze the file"})
        return
    }

    fmt.Printf("Returning scan result: %+v\n", scanResultPtr) // Debug statement
    c.JSON(200, scanResultPtr)
}

func main() {
    router := gin.Default()
    router.GET("/getScanResult", LogAnalyzerHandler)
    router.Run("localhost:8080")
}
