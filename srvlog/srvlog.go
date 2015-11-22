package srvlog

import (
	"fmt"
	"github.com/mgutz/ansi"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	logToConsol = false
	logToFile   = false
	filePath    = "/home/mladen/Desktop/test.log"

	//log queue
	logChan = make(chan logItem, 1024)

	//colour functions
	startupColourFunc = ansi.ColorFunc("green+b:black")
	fatalColourFunc   = ansi.ColorFunc("red+b:black")
	generalColourFunc = ansi.ColorFunc("blue+b:black")
	warningColourFunc = ansi.ColorFunc("yellow+b:black")
)

const (
	startupColour int = iota
	fatalColour
	generalColour
	warningColour
)

func Init(consolLog, fileLog bool, pathToFile string) {
	logToConsol = consolLog
	logToFile = fileLog
	filePath = pathToFile

	if logToConsol || logToFile {
		go listen()
	}
}

func Startup(args ...interface{}) {
	logChan <- logItem{
		longTime:         true,
		preset:           "STARTUP:",
		presetColourFunc: startupColour,
		content:          stringFromArgs(args...),
	}
}

func Fatal(args ...interface{}) {
	logChan <- logItem{
		longTime:         false,
		preset:           "FATAL:",
		presetColourFunc: fatalColour,
		content:          stringFromArgs(args...),
	}
}

func General(args ...interface{}) {
	logChan <- logItem{
		longTime:         false,
		preset:           "GENERAL:",
		presetColourFunc: generalColour,
		content:          stringFromArgs(args...),
	}
}

func Warning(args ...interface{}) {
	logChan <- logItem{
		longTime:         false,
		preset:           "WARNING:",
		presetColourFunc: warningColour,
		content:          stringFromArgs(args...),
	}
}

func stringFromArgs(args ...interface{}) string {
	var formatStr string
	for i := range args {
		if i != 0 {
			formatStr += " "
		}
		formatStr += "%v"
	}
	return fmt.Sprintf(formatStr, args...)
}

type logItem struct {
	longTime         bool
	preset           string
	presetColourFunc int
	content          string
}

func listen() {
	for {
		select {
		case item := <-logChan:
			writeToConsole(item)
			writeToFile(item)
			if item.presetColourFunc == fatalColour {
				os.Exit(1)
			}
		}
	}
}

func writeToFile(item logItem) {
	if logToFile {
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalln(err)
		}

		line := getTimestamp(item.longTime) + " " + item.preset + " " + item.content + "\n"
		_, err = file.WriteString(line)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func writeToConsole(item logItem) {
	if logToConsol {
		fmt.Println(getTimestamp(item.longTime), colourInText(item.preset, item.presetColourFunc), item.content)
	}
}

func colourInText(text string, colourFunc int) string {
	switch colourFunc {
	case startupColour:
		return startupColourFunc(text)
	case fatalColour:
		return fatalColourFunc(text)
	case generalColour:
		return generalColourFunc(text)
	case warningColour:
		return warningColourFunc(text)
	default:
		return text
	}
}

func getTimestamp(long bool) string {
	now := time.Now().Local()
	hour, min, sec := now.Clock()
	clockString := strconv.Itoa(hour) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec)
	if long {
		year, month, day := now.Date()
		dateString := strconv.Itoa(year) + "/" + month.String() + "/" + strconv.Itoa(day)
		return dateString + " " + clockString
	}
	return clockString
}
