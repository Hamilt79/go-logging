package gologging

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

var mainLogger Logger

type Logger struct
{
    LogFileName string
    DirLocation string
}

func isWindows() bool {
    return runtime.GOOS == "windows"
}

func isValidDir(location string) bool {
    _, err := os.Stat(location)
    if err != nil {
        return false
    }
    return true
}

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

func makeLoggingDir() error { 
    exists := isValidDir(mainLogger.dirLocation) 
    if (exists) {
        return nil;
    }
    err := os.MkdirAll(mainLogger.dirLocation, 0755)

    return err
}

func makeLoggingFile() error {
    fileLocation := filepath.Join(mainLogger.dirLocation, mainLogger.logFileName)
    _, err := os.Create(fileLocation)
    return err 
}

func PrintLn() {
    
}

func Init(logger Logger) {
    newLocation, err := expandDirectory(logger.dirLocation)
    if err != nil {
        panic(err)
    }
    logger.dirLocation = newLocation
    mainLogger = logger
    err = makeLoggingDir()
    if err != nil {
        panic(err)
    }
    err = makeLoggingFile()
    if err != nil {
        panic(err)
    }
}

func main() {
    logger := Logger { logFileName: "TestLogg.txt", dirLocation: "~/.logging/"} 
    Init(logger)
}
