package common

import (
	"os"
	"path"
	"sync"
	"time"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

const defaultLoggerFormat = "%{color}%{time:2006-01-02 15:04:05.000} [%{level:.5s}][%{module}][%{shortfile}] %{message} %{color:reset}"

var (
	format    = logging.MustStringFormatter(defaultLoggerFormat)
	loggers   = make(map[string]*SdkLogger)
	logger    = logging.MustGetLogger("common")
	backend   logging.LeveledBackend
	conf      *viper.Viper
	once      sync.Once
	lock      sync.RWMutex
	autoReset = true
	fileName  = "gosdk" + time.Now().Format("-2006-01-02-15:04:05 PM") + ".log"
)

type SdkLogger struct {
	*logging.Logger
	Conf *logConf
}

type logConf struct {
	FileName string
}

func newConsoleBackend(vip *viper.Viper) logging.LeveledBackend {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)

	// set level
	logLevel := vip.GetString(LogOutputLevel)

	logger.Debugf("[CONFIG]: %s = %v", LogOutputLevel, logLevel)

	level, _ := logging.LogLevel(logLevel)

	backendLeveled.SetLevel(level, "")

	return backendLeveled
}

func newFileBackend(vip *viper.Viper) logging.LeveledBackend {
	dir := vip.GetString(LogDir)
	logger.Debugf("[CONFIG]: %s = %v", LogDir, dir)

	filePath := path.Join(dir, fileName)
	_ = os.MkdirAll(dir, 0777)
	file, err := os.Create(filePath)
	if err != nil {
		logger.Errorf("create %s failed", filePath)
	}
	fileBackend := logging.NewLogBackend(file, "", 0)
	fileBackendFormatter := logging.NewBackendFormatter(fileBackend, format)
	fileBackendLeveled := logging.AddModuleLevel(fileBackendFormatter)

	// set level
	logLevel := vip.GetString(LogOutputLevel)

	logger.Debugf("[CONFIG]: %s = %v", LogOutputLevel, logLevel)

	level, _ := logging.LogLevel(logLevel)

	fileBackendLeveled.SetLevel(level, "")

	return fileBackendLeveled
}

func Reset() {
	autoReset = true
}

func updateBackend() {
	t := time.NewTicker(24 * time.Hour)
	if autoReset {
		for range t.C {
			SetBackend()
		}
	}
}

func InitLog(vip *viper.Viper) {
	once.Do(func() {
		conf = vip
		SetBackend()
		go updateBackend()
	})
}

func SetBackend() {
	consoleBackendLeveled := newConsoleBackend(conf)
	fileBackendLeveled := newFileBackend(conf)

	backend = logging.MultiLogger(consoleBackendLeveled, fileBackendLeveled)
	lock.Lock()
	for _, logger := range loggers {
		logger.SetBackend(backend)
	}
	lock.Unlock()
}

func SetCustomBackend(customBackend logging.LeveledBackend) {
	lock.Lock()
	defer lock.Unlock()
	autoReset = false
	for _, logger := range loggers {
		logger.SetBackend(customBackend)
	}
}

func GetLogger(module string) *SdkLogger {
	var logger *logging.Logger
	var logConf = new(logConf)
	var sdkLogger = new(SdkLogger)

	lock.Lock()
	defer lock.Unlock()
	if loggers[module] != nil {
		return loggers[module]
	} else {
		logger = logging.MustGetLogger(module)
		if backend != nil {
			logger.SetBackend(backend)
		}
		logConf.FileName = fileName
		sdkLogger.Logger = logger
		sdkLogger.Conf = logConf
		loggers[module] = sdkLogger
	}
	return sdkLogger
}
