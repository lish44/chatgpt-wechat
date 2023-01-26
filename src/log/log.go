package log

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

type logLevel uint32

type Log log.BufferPool

const (
	Panic logLevel = iota
	Fatal
	Error
	Warn
	Info
	Debug
	Trace
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	//log.AddHook(&myHook{})

	//setLogOutPut()
}

func setLogOutPut() {

	f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err)
	}
	log.SetOutput(f)
}

type myHook struct {
	recordNum int64
}

func (h *myHook) Fire(entry *log.Entry) error {
	h.recordNum++
	entry.Data["record_num "] = h.recordNum
	//entry.Data["time"] = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (h *myHook) Levels() []log.Level {
	return log.AllLevels
}

var fieldsMapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{})
	},
}

func Record(level logLevel, message string, fields ...interface{}) {
	var lvl log.Level
	switch level {
	case Panic:
		lvl = log.PanicLevel
	case Fatal:
		lvl = log.FatalLevel
	case Error:
		lvl = log.ErrorLevel
	case Warn:
		lvl = log.WarnLevel
	case Info:
		lvl = log.InfoLevel
	case Debug:
		lvl = log.DebugLevel
	case Trace:
		lvl = log.TraceLevel
	}
	if fields != nil {
		fieldsMap := fieldsMapPool.Get().(map[string]interface{})
		defer fieldsMapPool.Put(fieldsMap)

		for i := 0; i < len(fields); i += 2 {
			fieldsMap[fields[i].(string)] = fields[i+1]
		}

		log.WithFields(log.Fields(fieldsMap)).Log(lvl, message)
		return
	}

	log.WithFields(nil).Log(lvl, message)
}
