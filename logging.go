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
    fullPath string
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
    exists := isValidDir(mainLogger.DirLocation) 
    if (exists) {
        return nil;
    }
    err := os.MkdirAll(mainLogger.DirLocation, 0755)

    return err
}

func makeLoggingFile() error {
    fileLocation := filepath.Join(mainLogger.DirLocation, mainLogger.LogFileName)
    mainLogger.fullPath = fileLocation
    _, err := os.Create(fileLocation)
    return err 
}

func Print(str string) {
    file, err := os.OpenFile(mainLogger.fullPath, os.O_APPEND|os.O_WRONLY, 0644)
    if (err != nil) {
        panic(err)
    }
    file.WriteString(str)
    file.Close()
}

func PrintLn(str string) {
    file, err := os.OpenFile(mainLogger.fullPath, os.O_APPEND|os.O_WRONLY, 0644)
    if (err != nil) {
        panic(err)
    }
    file.WriteString(str + "\n")
    file.Close()
}

func Init(logger Logger) {
    newLocation, err := expandDirectory(logger.DirLocation)
    if err != nil {
        panic(err)
    }
    logger.DirLocation = newLocation
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
    logger := Logger { LogFileName: "TestLogg.txt", DirLocation: "~/.logging/"} 
    Init(logger)
}
