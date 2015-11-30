package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
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
	lDebug logLevel = iota
	lError
	lInfo
)

const (
	App logType = iota
	Debug
	Steam
	Web
)

const (
	appLogFilename   = "app.log"
	steamLogFilename = "steam.log"
	webLogFilename   = "web.log"
	logDirectory     = "logs"
)

func getLogPath(lt logType) string {
	switch lt {
	case App:
		return path.Join(logDirectory, appLogFilename)
	case Steam:
		return path.Join(logDirectory, steamLogFilename)
	case Web:
		return path.Join(logDirectory, webLogFilename)
	default:
		return path.Join(logDirectory, appLogFilename)
	}
}

func getLogFilenameFromType(lt logType) string {
	switch lt {
	case App:
		return appLogFilename
	case Steam:
		return steamLogFilename
	case Web:
		return webLogFilename
	default:
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
	return f.Size() > cfg.MaximumLogSize
}

func logDirNeedsCleaning(lt logType, cfg *Config) bool {
	files, err := ioutil.ReadDir(logDirectory)
	if err != nil {
		return false
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
	return logCount > cfg.MaximumLogCount
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
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
			// re-create default app|web.log
			if err := createLogFile(lt); err != nil {
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
		} else {
			// other logs exist, rename current as app|web.log.largestnum+1 & re-create
			latestlognum, _, err := getLatestAndEarliestLogNum(lt)
			if err != nil {
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
			if err := os.Rename(getLogPath(lt), path.Join(logDirectory,
				fmt.Sprintf("%s.%d", getLogFilenameFromType(lt),
					latestlognum+1))); err != nil {
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
			// re-create
			if err := createLogFile(lt); err != nil {
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
		}
	}
	return nil
}

func writeLogEntry(lt logType, loglevel logLevel, msg string,
	text ...interface{}) error {
	cfg, err := ReadConfig()
	if err != nil {
		return err
	}
	if lt == App && !cfg.EnableAppLogging {
		return nil
	} else if lt == Debug && !cfg.EnableDebugMessages {
		return nil
	} else if lt == Steam && !cfg.EnableSteamLogging {
		return nil
	} else if lt == Web && !cfg.EnableWebLogging {
		return nil
	}

	if lt == Debug {
		fmt.Printf(fmt.Sprintf("[%s] %s - %s\n",
			loglevel, time.Now().Format("Mon Jan _2 15:04:05 2006 EST"),
			fmt.Sprintf(msg, text...)))
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
		loglevel, time.Now().Format("Mon Jan _2 15:04:05 2006 EST"),
		fmt.Sprintf(msg, text...)))
	if err != nil {
		return fmt.Errorf("Unable to write log entry '%s': %s\n", text, err)
	}
	f.Sync()
	return nil
}

func (l logLevel) String() string {
	switch l {
	case lDebug:
		return "Debug"
	case lError:
		return "Error"
	case lInfo:
		return "Info"
	default:
		return ""
	}
}

func WriteDebug(msg string, input ...interface{}) {
	if err := writeLogEntry(Debug, lDebug, msg, input...); err != nil {
		fmt.Print(err)
	}
}

func LogAppError(e error, input ...interface{}) error {
	_ = writeLogEntry(App, lError, e.Error(), input...)
	return fmt.Errorf(e.Error(), input...)
}

func LogAppErrorf(msg string, input ...interface{}) error {
	_ = writeLogEntry(App, lError, msg, input...)
	return fmt.Errorf(msg, input...)
}
func LogAppInfo(msg string, input ...interface{}) {
	_ = writeLogEntry(App, lInfo, msg, input...)
}

func LogSteamInfo(msg string, input ...interface{}) {
	_ = writeLogEntry(Steam, lInfo, msg, input...)
}

func LogSteamError(e error, input ...interface{}) error {
	_ = writeLogEntry(Steam, lError, e.Error(), input...)
	return fmt.Errorf(e.Error(), input...)
}

func LogSteamErrorf(msg string, input ...interface{}) error {
	_ = writeLogEntry(Steam, lError, msg, input...)
	return fmt.Errorf(msg, input...)
}

func LogWebError(e error, input ...interface{}) error {
	_ = writeLogEntry(Web, lError, e.Error(), input...)
	return fmt.Errorf(e.Error(), input...)
}

func LogWebErrorf(msg string, input ...interface{}) error {
	if err := writeLogEntry(Web, lError, msg, input...); err != nil {
		fmt.Print(err)
	}
	return fmt.Errorf(msg, input...)
}

func LogWebRequest(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inner.ServeHTTP(w, r)

		if err := writeLogEntry(Web, lInfo, fmt.Sprintf("%s\t%s\t%s\t%s",
			r.Method, r.RequestURI, r.RemoteAddr, name)); err != nil {
			fmt.Print(err)
		}
	})
}
