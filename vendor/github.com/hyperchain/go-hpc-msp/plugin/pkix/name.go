package pkix

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperchain/go-hpc-msp/plugin/common"
	"github.com/meshplus/crypto"
)

//for cert
//  GN=ecert,O=au73jsa03edesad67va,CN=node1,SERIALNUMBER=fd26a860237b461d1baec332
//for ca
//  必须有SERIALNUMBER，CN
//  CN=DummyCA是dummyca的标志，只有无CA模式出现
//  分布式CA下CN必须是hostname
//for pub
//  只有无CA模式出现

const (
	//Version cert organization version
	Version = "version"
	//VP cert organization vp, nvp band node
	VP = "vp"
	//Platform cert organization platform, use flato
	Platform = "platform"
)

//IdentityName identity name
type IdentityName struct {
	//organization，is a base64 encoded JSON dictionary
	O string
	//host name or addr, E.g :node1, 172.16.5.1, www.hyperchain.cn and so on
	CN string
	//cert class, E.g ecert
	GN string
	//serial number, E.g: fd26a860237b461d1baec332
	SerialNumber string
}

//String fmt.string
func (n *IdentityName) String() string {
	r := make([]string, 4)
	for i := 0; i < 4; i++ {
		switch {
		case i == 0 && n.SerialNumber != "":
			r[3] = "SERIALNUMBER=" + n.SerialNumber
		case i == 1 && n.O != "":
			r[0] = "O=" + n.O
		case i == 2 && n.GN != "":
			r[1] = "GN=" + n.GN
		case i == 3 && n.CN != "":
			r[2] = "CN=" + n.CN
		}
	}
	return strings.Join(r, ",")
}

//GetCertType get CertType
func (n *IdentityName) GetCertType() crypto.CertType {
	return common.ParseCertType(n.GN)
}

//GetIdentityNameFromString get IdentityName from string
func GetIdentityNameFromString(s string) *IdentityName {
	if s == common.DummyCAName {
		return &IdentityName{
			CN: common.DummyCAName,
		}
	}
	n := new(IdentityName)
	r := strings.Split(s, ",")
	for i := range r {
		rr := strings.SplitN(r[i], "=", 2)
		if len(rr) != 2 {
			continue
		}
		switch strings.ToUpper(rr[0]) {
		case "SERIALNUMBER":
			n.SerialNumber = rr[1]
		case "GN":
			n.GN = rr[1]
		case "O":
			n.O = rr[1]
		case "CN":
			n.CN = rr[1]
		}
	}
	return n
}

//GetIdentityNameFromPKIXName get IdentityName from PKIXName
func GetIdentityNameFromPKIXName(name Name) *IdentityName {
	//http://tools.ietf.org/html/rfc5280#section-4.1.2.4
	oid := []int{2, 5, 4, 42}
	n := new(IdentityName)
	n.O = name.Organization[0]
	n.CN = name.CommonName
	n.SerialNumber = name.SerialNumber
	//now use OU, Keep GN for compatibility GN, OU string
	if len(name.OrganizationalUnit) > 0 && common.ParseCertType(name.OrganizationalUnit[0]) != crypto.UnknownCertType {
		n.GN = name.OrganizationalUnit[0]
		return n
	}

	var ok bool
	for i := range name.Names {
		if name.Names[i].Type.Equal(oid) {
			n.GN, ok = name.Names[i].Value.(string)
			if ok {
				break
			}
		}
	}
	return n
}

//ParseOrganization get Organization map
//input could be empty
func ParseOrganization(O string) (map[string]string, error) {
	oMap := make(map[string]string)
	//flato: see go-hpc-p2p HTS.GenCRS
	if len(O) == 0 || O == common.DummyID || O == "flato" || O == "hyperchain" {
		return oMap, nil
	}
	o, innerErr := base64.StdEncoding.DecodeString(O)
	if innerErr != nil {
		return nil, fmt.Errorf("base64 decode cert failed, reason: %v", innerErr)
	}

	innerErr = json.Unmarshal(o, &oMap)
	if innerErr != nil {
		return nil, fmt.Errorf("json unmarshal cert failed, reason : %v", innerErr)
	}
	return oMap, nil
}

//GetOrganization get org
func GetOrganization(ext map[string]string) string {
	target := make(map[string]string)
	target[Version] = ext[Version]
	target[Platform] = ext[Platform]
	target[VP] = ext[VP]
	j, _ := json.Marshal(target)
	return base64.StdEncoding.EncodeToString(j)
}
