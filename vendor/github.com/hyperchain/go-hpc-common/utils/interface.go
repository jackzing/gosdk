package utils

// Logger is the statedb logger interface which managers logger output.
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Notice(v ...interface{})
	Noticef(format string, v ...interface{})

	Warning(v ...interface{})
	Warningf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Critical(v ...interface{})
	Criticalf(format string, v ...interface{})
}

// PendingQueue is a sequential waiter interface
type PendingQueue interface {
	Start(uint64) error
	Stop()
	Wait(seqNo uint64)
	Done(seqNo uint64)
	Reset(uint64) error
}
