package gologging

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// mainLogger is the main logger that stores file and directory information
var mainLogger Logger

var errorExists bool

// Logger is a struct that contains the log file name, the directory location and the full path of the log file
type Logger struct
{
    LogFileName string
    DirLocation string
    fullPath string
}

// isWindows checks if the operating system is Windows
func isWindows() bool {
    return runtime.GOOS == "windows"
}

// isValidDir checks if the directory exists
func isValidDir(location string) bool {
    _, err := os.Stat(location)
    if err != nil {
        return false
    }
    return true
}

// expandDirectory expands the directory location if it starts with the ~ directory
func expandDirectory(location string) (string, error) {
    usr, err := user.Current()
    if err != nil {
        return "", err
    }
    dir := usr.HomeDir
    if location == "~" {
        // In case of "~", which won't be caught by the "else if"
        location = dir
    } else if strings.HasPrefix(location, "~/") {
        // Use strings.HasPrefix so we don't match paths like
        // "/something/~/something/"
        location = filepath.Join(dir, location[2:])
    }
    return location, nil
}

// makeLoggingDir creates the log directory
func makeLoggingDir() error { 
    exists := isValidDir(mainLogger.DirLocation) 
    if (exists) {
        return nil;
    }
    err := os.MkdirAll(mainLogger.DirLocation, 0755)

    return err
}

// makeLoggingFile creates the log file
func makeLoggingFile() error {
    fileLocation := filepath.Join(mainLogger.DirLocation, mainLogger.LogFileName)
    mainLogger.fullPath = fileLocation
    _, err := os.Create(fileLocation)
    return err 
}

// HasError checks if an error has occurred in printing
func HasError() bool {
    return errorExists
}

// Print prints a string to the log file without adding a newline character
func Print(str string) {
    file, err := os.OpenFile(mainLogger.fullPath, os.O_APPEND|os.O_WRONLY, 0644)
    if (err != nil) {
        errorExists = true
    }
    file.WriteString(str)
    file.Close()
}

// PrintLn prints a string to the log file and adds a newline character
func PrintLn(str string) {
    file, err := os.OpenFile(mainLogger.fullPath, os.O_APPEND|os.O_WRONLY, 0644)
    if (err != nil) {
        errorExists = true
    }
    file.WriteString(str + "\n")
    file.Close()
}

// Init initializes the logger
// Uses the Logger struct to set the directory location and the log file name
func Init(logger Logger) error {
    errorExists = false
    // Expanding the directory location if it starts with the ~ directory
    // $HOME is not supported and ~ is only supported in the beginning of the string
    newLocation, err := expandDirectory(logger.DirLocation)
    if err != nil { 
        return err
    }
    // Setting the new directory location
    logger.DirLocation = newLocation
    // Setting the mainLogger
    mainLogger = logger
    // Making the directory that should be logged to
    err = makeLoggingDir()
    if err != nil {
        return err
    }
    // Making the file that should be logged to
    err = makeLoggingFile()
    if err != nil {
        return err
    }
    return nil
}

