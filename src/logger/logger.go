package logger

import (
	"a2sapi/src/config"
	"a2sapi/src/constants"
	"a2sapi/src/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type logLevel int

const (
	lDebug logLevel = iota
	lError
	lInfo
)

func getLogPath(lt constants.LogType) string {
	switch lt {
	case constants.LTypeApp:
		return constants.AppLogFilePath
	case constants.LTypeSteam:
		return constants.SteamLogFilePath
	case constants.LTypeWeb:
		return constants.WebLogFilePath
	default:
		return constants.AppLogFilePath
	}
}

func getLogFilenameFromType(lt constants.LogType) string {
	switch lt {
	case constants.LTypeApp:
		return constants.AppLogFilename
	case constants.LTypeSteam:
		return constants.SteamLogFilename
	case constants.LTypeWeb:
		return constants.WebLogFilename
	default:
		return constants.AppLogFilename
	}
}

func deleteLogs(lt constants.LogType) error {
	logfiles, err := getLogFiles(lt)
	if err != nil {
		return err
	}
	for _, f := range logfiles {
		if err := os.Remove(path.Join(constants.LogDirectory, f)); err != nil {
			return err
		}
	}
	return nil
}

func isMaxLogSizeExceeded(lt constants.LogType, cfg *config.Config) bool {
	f, err := os.Stat(getLogPath(lt))
	if err != nil {
		return false
	}
	return f.Size() > cfg.LogConfig.MaximumLogSize*1024
}

func logDirNeedsCleaning(lt constants.LogType, cfg *config.Config) bool {
	files, err := ioutil.ReadDir(constants.LogDirectory)
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
	return logCount > cfg.LogConfig.MaximumLogCount
}

func getLogFiles(lt constants.LogType) ([]string, error) {
	files, err := ioutil.ReadDir(constants.LogDirectory)
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

func getLatestAndEarliestLogNum(lt constants.LogType) (latest int, earliest int,
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

func getLatestAndEarliestLog(lt constants.LogType) (latestlogfile string,
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

func verifyLogPaths(lt constants.LogType) error {
	if err := util.CreateDirectory(constants.LogDirectory); err != nil {
		return err
	}
	if !util.FileExists(getLogPath(lt)) {
		if err := util.CreateEmptyFile(getLogPath(lt), false); err != nil {
			return fmt.Errorf("verifyLogPaths error: %s", err)
		}
	}
	return nil
}

func verifyLogSettings(lt constants.LogType, cfg *config.Config) error {
	// too many stale logs
	if logDirNeedsCleaning(lt, cfg) {
		// delete stale logs
		if err := deleteLogs(lt); err != nil {
			return fmt.Errorf("verifyLogSettings error: %s\n", err)
		}
		// re-create default logfile
		if err := util.CreateEmptyFile(getLogPath(lt), true); err != nil {
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
			if err := os.Rename(getLogPath(lt), path.Join(constants.LogDirectory,
				fmt.Sprintf("%s.1", getLogFilenameFromType(lt)))); err != nil {
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
			// re-create default app|web.log
			if err := util.CreateEmptyFile(getLogPath(lt), true); err != nil {
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
		} else {
			// other logs exist, rename current as app|web.log.largestnum+1 & re-create
			latestlognum, _, err := getLatestAndEarliestLogNum(lt)
			if err != nil {
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
			if err := os.Rename(getLogPath(lt), path.Join(constants.LogDirectory,
				fmt.Sprintf("%s.%d", getLogFilenameFromType(lt),
					latestlognum+1))); err != nil {
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
			// re-create
			if err := util.CreateEmptyFile(getLogPath(lt), true); err != nil {
				return fmt.Errorf("verifyLogSettings error: %s\n", err)
			}
		}
	}
	return nil
}

func writeLogEntry(lt constants.LogType, loglevel logLevel, msg string,
	text ...interface{}) error {
	cfg := config.ReadConfig()

	if lt == constants.LTypeApp && !cfg.LogConfig.EnableAppLogging {
		return nil
	} else if lt == constants.LTypeDebug && !cfg.DebugConfig.EnableDebugMessages {
		return nil
	} else if lt == constants.LTypeSteam && !cfg.LogConfig.EnableSteamLogging {
		return nil
	} else if lt == constants.LTypeWeb && !cfg.LogConfig.EnableWebLogging {
		return nil
	}

	if lt == constants.LTypeDebug {
		fmt.Printf(fmt.Sprintf("[%s] %s - %s\n",
			loglevel, time.Now().Format("Mon Jan 2 15:04:05 2006 EST"),
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

// WriteDebug writes the specified debug message to to stdout, if enabled.
func WriteDebug(msg string, input ...interface{}) {
	if err := writeLogEntry(constants.LTypeDebug, lDebug, msg, input...); err != nil {
		fmt.Print(err)
	}
}

// LogAppError logs application-related errors, if enabled, to the app log file.
func LogAppError(e error, input ...interface{}) error {
	_ = writeLogEntry(constants.LTypeApp, lError, e.Error(), input...)
	return fmt.Errorf(e.Error(), input...)
}

// LogAppErrorf logs formatted application-related errors, if enabled, to the app log file.
func LogAppErrorf(msg string, input ...interface{}) error {
	_ = writeLogEntry(constants.LTypeApp, lError, msg, input...)
	return fmt.Errorf(msg, input...)
}

// LogAppInfo logs application-related info messages, if enabled, to the app log file.
func LogAppInfo(msg string, input ...interface{}) {
	_ = writeLogEntry(constants.LTypeApp, lInfo, msg, input...)
}

// LogSteamInfo logs Steam-related info messages, if enabled, to the Steam log file.
func LogSteamInfo(msg string, input ...interface{}) {
	_ = writeLogEntry(constants.LTypeSteam, lInfo, msg, input...)
}

// LogSteamError logs Steam-related errors, if enabled, to the Steam log file.
func LogSteamError(e error, input ...interface{}) error {
	_ = writeLogEntry(constants.LTypeSteam, lError, e.Error(), input...)
	return fmt.Errorf(e.Error(), input...)
}

// LogSteamErrorf logs formatted Steam-related errors, if enabled, to the Steam log file.
func LogSteamErrorf(msg string, input ...interface{}) error {
	_ = writeLogEntry(constants.LTypeSteam, lError, msg, input...)
	return fmt.Errorf(msg, input...)
}

// LogWebError logs API-related web errors, if enabled, to the web log file.
func LogWebError(e error, input ...interface{}) error {
	_ = writeLogEntry(constants.LTypeWeb, lError, e.Error(), input...)
	return fmt.Errorf(e.Error(), input...)
}

// LogWebErrorf logs formatted API-related web errors, if enabled, to the web log file.
func LogWebErrorf(msg string, input ...interface{}) error {
	if err := writeLogEntry(constants.LTypeWeb, lError, msg, input...); err != nil {
		fmt.Print(err)
	}
	return fmt.Errorf(msg, input...)
}

// LogWebRequest logs web requests as debug messages, if enabled to stdout, as well
// as logging web requests as info messages to the web log file.
func LogWebRequest(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inner.ServeHTTP(w, r)
		u, err := url.QueryUnescape(r.URL.String())
		if err != nil {
			u = fmt.Sprintf("Invalid URL [missing 2 chars after percent sign]: %s",
				r.URL.String())
		}
		if err := writeLogEntry(constants.LTypeDebug, lDebug, fmt.Sprintf(
			"URL: %s\tPATH: %s\tQUERY:%v", u, r.URL.Path,
			r.URL.Query())); err != nil {
			fmt.Print(err)
		}

		if err := writeLogEntry(constants.LTypeWeb, lInfo, fmt.Sprintf("%s\t%s\t%s\t%s",
			r.Method, u, r.RemoteAddr, name)); err != nil {
			fmt.Print(err)
		}
	})
}
