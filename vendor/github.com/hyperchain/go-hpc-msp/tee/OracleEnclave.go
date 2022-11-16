package tee

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperchain/go-hpc-msp/config"
)

//HTTPS_POST HTTPS_GET HTTPS_PUT HTTPS_DELETE support https header as nil
const (
	HTTPS_POST   = "POST"
	HTTPS_GET    = "GET"
	HTTPS_PUT    = "PUT"
	HTTPS_DELETE = "DELETE"

	HTTPS_PATCH    = "PATCH"
	HTTPS_COPY     = "COPY"
	HTTPS_HEAD     = "HEAD"
	HTTPS_OPTIONS  = "OPTIONS"
	HTTPS_LINK     = "LINK"
	HTTPS_UNLINK   = "UNLINK"
	HTTPS_PURGE    = "PURGE"
	HTTPS_LOCK     = "LOCK"
	HTTPS_UNLOCK   = "UNLOCK"
	HTTPS_PROPFIND = "PROPFIND"
	HTTPS_VIES     = "VIEW"
)

//NewOracleEnclave is an instance of tee oracle
func NewOracleEnclave(config config.EncryptionConfigInterface) (Oracle, error) {
	err := loadEnclave(config.EncryptionEngineCrypt())
	if err != nil {
		return nil, err
	}

	e := &OracleEnclaveImpl{}
	return e, nil
}

//OracleEnclaveImpl oracle enclave
type OracleEnclaveImpl struct{}

//Response https response
type Response struct {
	Status int
	Body   string
	Header map[string]string
}

//Close sgx close
func (os *OracleEnclaveImpl) Close() {
	if eid == 0 {
		return
	}
	close(eid)
}

//IsSGX return true
func (os *OracleEnclaveImpl) IsSGX() bool {
	return true
}

//Fetch https
func (os *OracleEnclaveImpl) Fetch(url string, method string, body string, head map[string]string) (Response, error) {

	reponse := Response{}
	if !strings.HasPrefix(url, "https") {
		return reponse, errors.New("url should have prefix, eg. https://xxx")
	}

	var host string
	s := strings.Split(url, "/")
	if len(s) < 3 {
		return reponse, errors.New("parse url failed")
	}
	host = s[2]

	if len(head) == 0 {
		switch method {
		case HTTPS_GET:
			head = make(map[string]string)
			head["Host"] = host
			head["Accept"] = "*/*"
		case HTTPS_POST, HTTPS_DELETE, HTTPS_PUT:
			head = make(map[string]string)
			head["Host"] = host
			head["Accept"] = "*/*"
			head["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
			head["Content-Length"] = strconv.Itoa(len(body))
		default:
			return reponse, fmt.Errorf("current method %v not support request header is nil", method)
		}
	}

	return getHTTPSResponse(eid, method, url, body, head)
}

//StatusReturn return message by response status
func StatusReturn(s int) string {
	switch s {
	case 200:
		return "OK"
	case 301:
		return "Moved Permanently"
	case 302:
		return "Found"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 410:
		return "Gone"
	case 500:
		return "Internal Server Error"
	default:
		return "status invalid, see other reason"
	}
}
