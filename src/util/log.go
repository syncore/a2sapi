package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type logLevel int
type logType int

const (
	lError logLevel = iota
	lInfo
)

const (
	App logType = iota
	Web
)

// Spf (fmt.Sprintf alias) - Some packages importing util for logging wouldn't
// otherwise need fmt
var Spf = fmt.Sprintf

const (
	appLogFilename = "app.log"
	webLogFilename = "web.log"
	logDirectory   = "logs"
)

func getLogPath(lt logType) string {
	if lt == App {
		return path.Join(logDirectory, appLogFilename)
	} else if lt == Web {
		return path.Join(logDirectory, webLogFilename)
	} else {
		return path.Join(logDirectory, appLogFilename)
	}
}

func getLogFilenameFromType(lt logType) string {
	if lt == App {
		return appLogFilename
	} else if lt == Web {
		return webLogFilename
	} else {
		return appLogFilename
	}
}

func createLogDir() error {
	if DirExists(logDirectory) {
		return nil
	}
	if err := os.Mkdir(logDirectory, os.ModeDir); err != nil {
		return err
	}
	return nil
}

func createLogFile(lt logType) error {
	f, err := os.Create(getLogPath(lt))
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func deleteLogs(lt logType) error {
	logfiles, err := getLogFiles(lt)
	if err != nil {
		return err
	}
	for _, f := range logfiles {
		if err := os.Remove(path.Join(logDirectory, f)); err != nil {
			return err
		}
	}
	return nil
}

func isMaxLogSizeExceeded(lt logType, cfg *Config) bool {
	f, err := os.Stat(getLogPath(lt))
	if err != nil {
		return false
	}
	if lt == App {
		return f.Size() > cfg.MaximumAppLogSize*1024
	}
	return f.Size() > cfg.MaximumWebLogSize*1024
}

func logDirNeedsCleaning(lt logType, cfg *Config) bool {
	files, err := ioutil.ReadDir(logDirectory)
	if err != nil {
		return false
	}
	var maxNumLogs int
	if lt == App {
		maxNumLogs = cfg.MaximumAppLogCount
	} else {
		maxNumLogs = cfg.MaximumWebLogCount
	}

	re := regexp.MustCompile(fmt.Sprintf("(?i)%s*.*", getLogFilenameFromType(lt)))
	logCount := 0
	for _, f := range files {
		match := re.FindStringSubmatch(f.Name())
		if f.IsDir() {
			continue
		}
		if len(match) != 1 {
			continue
		}
		logCount++
	}
	return logCount > maxNumLogs
}

func getLogFiles(lt logType) ([]string, error) {
	files, err := ioutil.ReadDir(logDirectory)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(fmt.Sprintf("(?i)(%s.)(\\d+)",
		getLogFilenameFromType(lt)))
	var logfiles []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		result := re.FindStringSubmatch(f.Name())
		if len(result) != 3 {
			continue
		}
		logfiles = append(logfiles, result[0])
	}
	return logfiles, nil
}

func getLatestAndEarliestLogNum(lt logType) (latest int, earliest int,
	err error) {
	logfiles, err := getLogFiles(lt)
	if err != nil {
		return 0, 0, err
	}
	re := regexp.MustCompile(fmt.Sprintf("(?i)(%s.)(\\d+)",
		getLogFilenameFromType(lt)))
	for i, f := range logfiles {
		result := re.FindStringSubmatch(f)
		if len(result) != 3 {
			continue
		}
		num, err := strconv.Atoi(result[2])
		if err != nil {
			continue
		}
		if i == 0 {
			earliest = num
			latest = num
		}
		if num > latest {
			latest = num
		}
		if num < earliest {
			earliest = num
		}
	}
	return latest, earliest, nil
}

func getLatestAndEarliestLog(lt logType) (latestlogfile string,
	earliestlogfile string, err error) {
	latestlognum, earliestlognum, err := getLatestAndEarliestLogNum(lt)
	if err != nil {
		return "", "", err
	}
	if earliestlognum == 0 {
		earliestlogfile = getLogFilenameFromType(lt)
	} else {
		earliestlogfile = fmt.Sprintf("%s.%d", getLogFilenameFromType(lt),
			earliestlognum)
	}
	if latestlognum == 0 {
		latestlogfile = getLogFilenameFromType(lt)
	} else {
		latestlogfile = fmt.Sprintf("%s.%d", getLogFilenameFromType(lt),
			latestlognum)
	}
	return latestlogfile, earliestlogfile, nil

}

func verifyLogPaths(lt logType) error {
	if err := createLogDir(); err != nil {
		return err
	}
	if !FileExists(getLogPath(lt)) {
		if err := createLogFile(lt); err != nil {
			return fmt.Errorf("verifyLogPaths error: %s", err)
		}
	}
	return nil
}

func verifyLogSettings(lt logType, cfg *Config) error {
	// too many stale logs
	if logDirNeedsCleaning(lt, cfg) {
		// delete stale logs
		if err := deleteLogs(lt); err != nil {
			return fmt.Errorf("verifyLogSettings error: %s\n", err)
		}
		// re-create default logfile
		if err := createLogFile(lt); err != nil {
			return fmt.Errorf("verifyLogSettings error: %s\n", err)
		}
	}
	// log file to be written to is too large
	if isMaxLogSizeExceeded(lt, cfg) {
		latestlog, _, err := getLatestAndEarliestLog(lt)
		if err != nil {
			return fmt.Errorf("verifyLogSettings error: %s\n", err)
		}
		// only app|web.log exists, rename it to app|web.log.1 & re-create
		if strings.EqualFold(latestlog, getLogFilenameFromType(lt)) {
			if err := os.Rename(getLogPath(lt), path.Join(logDirectory,
				fmt.Sprintf("%s.1", getLogFilenameFromType(lt)))); err != nil {
				return fmt.Errorf("verfyLogSettings error: %s\n", err)
			}
			// re-create default app|web.log
			if err := createLogFile(lt); err != nil {
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
		} else {
			// other logs exist, rename current as app|web.log.largestnum+1 & re-create
			latestlognum, _, err := getLatestAndEarliestLogNum(lt)
			if err != nil {
				return fmt.Errorf("verfyLogSettings error: %s\n", err)
			}
			if err := os.Rename(getLogPath(lt), path.Join(logDirectory,
				fmt.Sprintf("%s.%d", getLogFilenameFromType(lt),
					latestlognum+1))); err != nil {
				return fmt.Errorf("verfyLogSettings error: %s\n", err)
			}
			// re-create
			if err := createLogFile(lt); err != nil {
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
		}
	}
	return nil
}

func writeLogEntry(lt logType, loglevel logLevel, text ...interface{}) error {
	cfg, err := ReadConfig()
	if err != nil {
		return err
	}
	if lt == App && !cfg.EnableAppLogging {
		return nil
	} else if lt == Web && !cfg.EnableWebLogging {
		return nil
	}
	if err := verifyLogPaths(lt); err != nil {
		return err
	}
	if err := verifyLogSettings(lt, cfg); err != nil {
		return err
	}
	f, err := os.OpenFile(getLogPath(lt), os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("Log entry write error, unable to open log file '%s': %s\n",
			getLogPath(lt), err)
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("[%s] %s - %s\n",
		loglevel, time.Now().Format("Mon Jan _2 15:04:05 2006 EST"), fmt.Sprintf(
			"%v", text...)))
	if err != nil {
		return fmt.Errorf("Unable to write log entry '%s': %s\n", text, err)
	}
	f.Sync()
	return nil
}

func (l logLevel) String() string {
	switch l {
	case lError:
		return "Error"
	case lInfo:
		return "Info"
	default:
		return ""
	}
}

func LogAppError(input ...interface{}) error {
	if err := writeLogEntry(App, lError, input...); err != nil {
		fmt.Print(err)
	}
	return fmt.Errorf("%s\n", input...)
}

func LogAppInfo(input ...interface{}) {
	if err := writeLogEntry(App, lInfo, input...); err != nil {
		fmt.Print(err)
	}
}

func LogWebError(input ...interface{}) error {
	if err := writeLogEntry(Web, lError, input...); err != nil {
		fmt.Print(err)
	}
	return fmt.Errorf("%s\n", input...)
}

func LogWebInfo(input ...interface{}) {
	if err := writeLogEntry(Web, lInfo, input...); err != nil {
		fmt.Print(err)
	}
}
