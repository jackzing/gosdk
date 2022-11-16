package logger

import (
	"fmt"
	"log"
	"os"
)

type logger log.Logger

var std = log.New(os.Stderr, "", log.LstdFlags)

func (*logger) Debug(v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[DEBU] ", fmt.Sprint(v...)))
}
func (*logger) Debugf(format string, v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[DEBU] ", fmt.Sprintf(format, v...)))
}
func (*logger) Info(v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[INFO] ", fmt.Sprint(v...)))
}
func (*logger) Infof(format string, v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[INFO] ", fmt.Sprintf(format, v...)))
}
func (*logger) Notice(v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[NOTI] ", fmt.Sprint(v...)))
}
func (*logger) Noticef(format string, v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[NOTI] ", fmt.Sprintf(format, v...)))
}
func (*logger) Warning(v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[WARN] ", fmt.Sprint(v...)))
}
func (*logger) Warningf(format string, v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[WARN] ", fmt.Sprintf(format, v...)))
}
func (*logger) Error(v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[ERRO] ", fmt.Sprint(v...)))
}
func (*logger) Errorf(format string, v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[ERRO] ", fmt.Sprintf(format, v...)))
}
func (*logger) Critical(v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[CRIT] ", fmt.Sprint(v...)))
}
func (*logger) Criticalf(format string, v ...interface{}) {
	_ = std.Output(2, fmt.Sprint("[CRIT] ", fmt.Sprintf(format, v...)))
}

//MSPLoggerSingleCase implement of Logger interface
var MSPLoggerSingleCase = (*logger)(log.New(os.Stdout, "msp", 5))

type fakeLogger int

func (fakeLogger) Debug(v ...interface{})                    {}
func (fakeLogger) Debugf(format string, v ...interface{})    {}
func (fakeLogger) Info(v ...interface{})                     {}
func (fakeLogger) Infof(format string, v ...interface{})     {}
func (fakeLogger) Notice(v ...interface{})                   {}
func (fakeLogger) Noticef(format string, v ...interface{})   {}
func (fakeLogger) Warning(v ...interface{})                  {}
func (fakeLogger) Warningf(format string, v ...interface{})  {}
func (fakeLogger) Error(v ...interface{})                    {}
func (fakeLogger) Errorf(format string, v ...interface{})    {}
func (fakeLogger) Critical(v ...interface{})                 {}
func (fakeLogger) Criticalf(format string, v ...interface{}) {}

//MSPFakeLogger no logger
var MSPFakeLogger fakeLogger
