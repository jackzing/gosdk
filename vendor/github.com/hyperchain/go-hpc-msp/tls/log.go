//go:build tlslog
// +build tlslog

package tls

import (
	"fmt"
	"log"
	"os"
)

var info, _ = os.OpenFile("tls_all.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
var erro, _ = os.OpenFile("tls_err.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
var warn, _ = os.OpenFile("tls_err.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

var (
	infoLog    = log.New(info, ">[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLog = log.New(warn, ">[WARN]: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog   = log.New(erro, ">[ERRO]: ", log.Ldate|log.Ltime|log.Lshortfile)
)

//Infof Infof
func Infof(s string, v ...interface{}) {
	_ = infoLog.Output(2, fmt.Sprintf(s, v...))
}

//Warningf Warningf
func Warningf(s string, v ...interface{}) {
	_ = warningLog.Output(2, fmt.Sprintf(s, v...))
}

//Errorf Errorf
func Errorf(s string, v ...interface{}) {
	_ = errorLog.Output(2, fmt.Sprintf(s, v...))
}
