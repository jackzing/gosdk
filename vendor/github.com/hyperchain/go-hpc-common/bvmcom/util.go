package bvmcom

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"regexp"
	"sort"
	"strings"

	"github.com/hyperchain/go-hpc-common/lservercom"
	"github.com/hyperchain/go-hpc-common/types"
	"github.com/hyperchain/go-hpc-common/types/protos"
	"github.com/hyperchain/go-hpc-common/utils"

	"github.com/fatih/set"
	"github.com/gogo/protobuf/proto"
)

const (
	// configPrefix + key -> config value
	configPrefix = ConfigPath

	// rolePrefix + address -> role weight
	rolePrefix = RolePath

	// rootCAPrefix + caName -> root ca
	rootCAPrefix = RootCAsPath

	// hostsPrefix + role -> hosts
	hostsPrefix = HostsPath

	// srsInfoPrefix + hash -> srs info
	srsInfoPrefix = SRSInfoPath
	// srsListPrefix + hash -> srs info
	srsListPrefix = SRSListPath

	sumOfWeight = "sum_of_weight"

	certPrefix       = CertListPath
	freezeCertPrefix = FreezeCertListPath

	split = "@"

	// defaultLen default len for record length of value
	defaultLen = 4

	//ChangeHashAlgoName execute change algo method name
	ChangeHashAlgoName = "ChangeHashAlgo"
	// ExecuteProposalName execute proposal method name.
	ExecuteProposalName = "Execute"
	// DirectProposalName direct execute proposal operation method name.
	DirectProposalName = "Direct"
	// AddRootCAName add root ca operate name .
	AddRootCAName = "AddRootCA"

	// cNamePrefix + name -> address
	cNamePrefix = CNSPath + split + "contractName"
	// cAddressPrefix + address -> name
	cAddressPrefix = CNSPath + split + "contractAddress"
)

// GetConfigKeyPrefix get config key prefix
func GetConfigKeyPrefix() []byte {
	return []byte(configPrefix)
}

// GetHostsKeyPrefix get hosts key prefix
func GetHostsKeyPrefix() []byte {
	return []byte(hostsPrefix)
}

// CompositeConfigKey composite config key
func CompositeConfigKey(key []byte) []byte {
	//configPrefix + split + key
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(configPrefix)
	buffer.WriteString(split)
	buffer.Write(key)
	return buffer.Bytes()
}

// SplitCompositeConfigKey split composite config key and retrieve real key
func SplitCompositeConfigKey(key []byte) ([]byte, bool) {
	prefix := GetConfigKeyPrefix()
	prefix = append(prefix, []byte(split)...)
	identifierLen := len(prefix)
	if bytes.HasPrefix(key, prefix) {
		return key[identifierLen:], true
	}
	return nil, false
}

// GetRoleKeyPrefix get role key prefix
func GetRoleKeyPrefix() []byte {
	return []byte(rolePrefix)
}

// GetRootCAKeyPrefix get root ca key prefix
func GetRootCAKeyPrefix() []byte {
	return []byte(rootCAPrefix)
}

// GetSplit get split
func GetSplit() []byte {
	return []byte(split)
}

// GetCNamePrefix get contract name prefix
func GetCNamePrefix() []byte {
	return []byte(cNamePrefix)
}

// GetCAddressPrefix get contract address prefix
func GetCAddressPrefix() []byte {
	return []byte(cAddressPrefix)
}

// GetCertPrefix get cert prefix
func GetCertPrefix() []byte {
	return []byte(certPrefix)
}

// CompositeRoleKey composite role key
func CompositeRoleKey(role string, address []byte) []byte {
	//rolePrefix +split + role + split + address
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(rolePrefix)
	buffer.WriteString(split)
	buffer.WriteString(role)
	buffer.WriteString(split)
	buffer.Write(address)
	return buffer.Bytes()
}

// CompositeRootCAKey composite root ca key
func CompositeRootCAKey(caName string) []byte {
	//rootCAPrefix +split + caName
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(rootCAPrefix)
	buffer.WriteString(split)
	buffer.WriteString(caName)
	return buffer.Bytes()
}

// CompositeCNameKey composite contract name key
func CompositeCNameKey(cname string) []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(cNamePrefix)
	buffer.WriteString(split)
	buffer.WriteString(cname)
	return buffer.Bytes()
}

// SplitCNameKey split contract name from key
func SplitCNameKey(key []byte) (string, bool) {
	prefix := GetCNamePrefix()
	prefix = append(prefix, []byte(split)...)
	identifierLen := len(prefix)
	if bytes.HasPrefix(key, prefix) && len(key) > identifierLen {
		return string(key[identifierLen:]), true
	}
	return "", false
}

// CompositeCAddressKey composite contract address key
func CompositeCAddressKey(address []byte) []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(cAddressPrefix)
	buffer.WriteString(split)
	buffer.Write(address)
	return buffer.Bytes()
}

// SplitCAddressKey split contract address from key
func SplitCAddressKey(key []byte) ([]byte, bool) {
	prefix := GetCAddressPrefix()
	prefix = append(prefix, []byte(split)...)
	identifierLen := len(prefix)
	if bytes.HasPrefix(key, prefix) && len(key) == identifierLen+types.AddressLength {
		return key[identifierLen:], true
	}
	return nil, false
}

// CheckCNameFmt check contract name format
func CheckCNameFmt(cnsName string) bool {
	// composite by letter, number, _, -
	nameExpr := "^[0-9a-zA-z_-]{2,40}$"
	compile, _ := regexp.Compile(nameExpr)
	return compile.MatchString(cnsName)
}

// SplitCompositeRoleKey split composite role key and retrieve role, address
func SplitCompositeRoleKey(key []byte) (string, []byte, bool) {
	prefix := GetRoleKeyPrefix()
	prefix = append(prefix, []byte(split)...)
	identifierLen := len(prefix)
	if bytes.HasPrefix(key, prefix) && len(key) > identifierLen+types.AddressLength {
		roleSplitLen := len([]byte(split))
		return string(key[identifierLen : len(key)-roleSplitLen-types.AddressLength]), key[len(key)-types.AddressLength:], true
	}
	return "", nil, false
}

// SplitCompositeRootCAKey split composite root ca key and retrieve ca name
func SplitCompositeRootCAKey(key []byte) (string, bool) {
	prefix := GetRootCAKeyPrefix()
	prefix = append(prefix, []byte(split)...)
	identifierLen := len(prefix)
	if bytes.HasPrefix(key, prefix) && len(key) > identifierLen {
		return string(key[identifierLen:]), true
	}
	return "", false
}

// CompositeSumOfWeightKey composite sum of weight key
func CompositeSumOfWeightKey(role string) []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(sumOfWeight)
	buffer.WriteString(split)
	buffer.WriteString(role)
	return buffer.Bytes()
}

// CompositeHostsKey composite hosts key
func CompositeHostsKey(role string) []byte {
	//HostsPrefix@role
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(hostsPrefix)
	buffer.WriteString(split)
	buffer.WriteString(strings.ToLower(role))
	return buffer.Bytes()
}

// CompositeCertKey composite key of revoked cert
func CompositeCertKey(cert []byte) []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(certPrefix)
	buffer.WriteString(split)
	buffer.Write(cert)
	return buffer.Bytes()
}

// CompositeCertValue composite value of revoked cert
func CompositeCertValue(v interface{}) []byte {
	res, _ := json.Marshal(v)
	return res
}

// DecompositeCertValue decomposite value of revoked cert
func DecompositeCertValue(raw []byte, e interface{}) (r interface{}) {
	_ = json.Unmarshal(raw, e)
	return e
}

// CompositeFreezeCertKey composite key of freeze cert
func CompositeFreezeCertKey(cert []byte) []byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(freezeCertPrefix)
	buffer.WriteString(split)
	buffer.Write(cert)
	return buffer.Bytes()
}

// Operation define the needed Params for call a method
type Operation struct {
	MethodName string
	Params     []string
}

// DecodeOperations Decode operations
func DecodeOperations(in []byte) []*Operation {
	var calls []*Operation
	operateLen := utils.BytesToInt32(in[:defaultLen])
	index := defaultLen
	for i := 0; i < operateLen; i++ {
		methodNameLen := utils.BytesToInt32(in[index : defaultLen+index])
		index += defaultLen
		methodName := string(in[index : index+methodNameLen])
		index += methodNameLen
		paramCount := utils.BytesToInt32(in[index : index+defaultLen])
		index += defaultLen
		params := make([]string, paramCount)
		for i := 0; i < paramCount; i++ {
			paramLen := utils.BytesToInt32(in[index : index+defaultLen])
			param := string(in[index+defaultLen : index+defaultLen+paramLen])
			params[i] = param
			index += defaultLen + paramLen
		}
		calls = append(calls, &Operation{
			MethodName: methodName,
			Params:     params,
		})
	}
	return calls
}

// EncodeOperates encode operations
func EncodeOperates(pas []*Operation) []byte {
	var configBytes []byte
	operationLenBytes := utils.IntToBytes4(len(pas))
	configBytes = append(configBytes, operationLenBytes[:]...)
	for _, pa := range pas {
		result := Encode(pa)
		configBytes = append(configBytes, result[:]...)
	}
	return configBytes
}

// Encode encode MethodName and Params to payload
func Encode(ope *Operation) (result []byte) {
	methodNameLenBytes := utils.IntToBytes4(len(ope.MethodName))
	paramsLenBytes := utils.IntToBytes4(len(ope.Params))
	result = append(result, methodNameLenBytes[:]...)
	result = append(result, []byte(ope.MethodName)...)
	result = append(result, paramsLenBytes[:]...)
	for _, param := range ope.Params {
		paramLenBytes := utils.IntToBytes4(len(param))
		result = append(result, paramLenBytes[:]...)
		result = append(result, []byte(param)...)
	}
	return
}

// Decode decode payload
func Decode(payload []byte) (string, []string, error) {
	if len(payload) < defaultLen {
		return "", nil, errors.New("invalid payload")
	}
	methodNameLen := utils.BytesToInt32(payload[:defaultLen])
	if len(payload) < defaultLen+methodNameLen {
		return "", nil, errors.New("invalid payload")
	}
	methodName := string(payload[defaultLen : defaultLen+methodNameLen])
	if len(payload) < defaultLen+methodNameLen+defaultLen {
		return "", nil, errors.New("invalid payload")
	}
	paramCount := utils.BytesToInt32(payload[defaultLen+methodNameLen : defaultLen+methodNameLen+defaultLen])
	paramsInput := payload[defaultLen+methodNameLen+defaultLen:]
	params := make([]string, paramCount)
	index := 0
	for i := 0; i < paramCount; i++ {
		if len(paramsInput) < index+defaultLen {
			return "", nil, errors.New("invalid payload")
		}
		paramLen := utils.BytesToInt32(paramsInput[index : index+defaultLen])
		if len(paramsInput) < index+defaultLen+paramLen {
			return "", nil, errors.New("invalid payload")
		}
		param := string(paramsInput[index+defaultLen : index+defaultLen+paramLen])
		params[i] = param
		index += defaultLen + paramLen
	}
	return methodName, params, nil
}

// UnmarshalProposal unmarshal proposalData
func UnmarshalProposal(data []byte) (proposal *protos.ProposalData, err error) {
	proposal = &protos.ProposalData{}
	err = proto.Unmarshal(data, proposal)
	return proposal, err
}

// UnmarshalVSet unmarshal vset
func UnmarshalVSet(data []byte) (types.VSet, error) {
	var vset types.VSet
	err := json.Unmarshal(data, &vset)
	return vset, err
}

// UnmarshalHosts unmarshal hosts
func UnmarshalHosts(data []byte) (types.Hosts, error) {
	hosts := make(types.Hosts)
	err := json.Unmarshal(data, &hosts)
	return hosts, err
}

// UnmarshalVPCerts unmarshal vpCerts
func UnmarshalVPCerts(data []byte) (types.VPCerts, error) {
	certs := make(types.VPCerts)
	err := json.Unmarshal(data, &certs)
	return certs, err
}

// UnmarshalAlgoMap Unmarshal algo map
func UnmarshalAlgoMap(data []byte) (types.BlockAlgoSets, error) {
	var algoSet types.BlockAlgoSets
	err := json.Unmarshal(data, &algoSet)
	return algoSet, err
}

// MarshalAlgoMap Marshal algo map
func MarshalAlgoMap(set types.BlockAlgoSets) ([]byte, error) {
	return json.Marshal(set)
}

// MarshalCerts marshal Certs
func MarshalCerts(s types.Certs) ([]byte, error) {
	data := []byte{}
	name := [][types.SignLen]byte{}
	for k := range s {
		name = append(name, k)
	}
	sort.Slice(name,
		func(i, j int) bool {
			for ii := 0; ii < types.SignLen; ii++ {
				if name[i][ii] != name[j][ii] {
					return name[i][ii] > name[j][ii]
				}
			}
			return true
		})

	index := 0
	for _, k := range name {
		td, _ := json.Marshal(s[k])
		ilen := types.SignLen + len(td) + 4
		tmp := int2Byte(index + ilen)
		data = append(data, tmp[:]...)
		data = append(data, k[:]...)
		data = append(data, td...)
		index = index + ilen
	}
	return data, nil
}

func int2Byte(data int) (ret [4]byte) {
	tmp := 0xff
	for index := uint(0); index < uint(4); index++ {
		ret[index] = byte((tmp << (index * 8) & data) >> (index * 8))
	}
	return ret
}

func byte2Int(data []byte) int {
	var ret int
	blen := len(data)
	for i := uint(0); i < uint(blen); i++ {
		ret = ret | (int(data[i]) << (i * 8))
	}
	return ret
}

// UnmarshalCerts unmarshal Certs
func UnmarshalCerts(data []byte, exist bool) (types.Certs, error) {
	sss := types.Certs{}
	dlen := len(data)
	if exist && dlen != 0 {
		for index := 0; index+4 < dlen; {
			ll := byte2Int(data[index : index+4])
			tt := [types.SignLen]byte{}
			ttt := types.CertInfo{}
			if index+4+types.SignLen <= dlen {
				copy(tt[:], data[index+4:index+4+types.SignLen])
				_ = json.Unmarshal(data[index+types.SignLen+4:ll], &ttt)
			}
			sss[tt] = ttt
			index = ll
		}
	}
	return sss, nil
}

// UnmarshalInputReplace unmarshal input replace
func UnmarshalInputReplace(data []byte) (types.InputReplace, error) {
	input := new(types.InputReplace)
	err := json.Unmarshal(data, &input)
	return *input, err
}

// UnmarshalInputNVPRevoke unmarshal input NVP revoke
func UnmarshalInputNVPRevoke(data []byte) (types.InputNVPRevoke, error) {
	input := new(types.InputNVPRevoke)
	err := json.Unmarshal(data, &input)
	return *input, err
}

// HasDeployPermission return if the deployer has permission to deploy contract with given license type.
// if has not permission to deploy contract, return false.
// if has permission to deploy contract, return true.
func HasDeployPermission(licenseTpy lservercom.LicenseType, deployer types.Address) bool {
	// check when license type is subLocal, if the deploy account is the specify account
	return licenseTpy != lservercom.SubLocal || deployer.Hex() == deployAddr.Hex()
}

// CompositeSRSKey composite srs key
func CompositeSRSKey(hash string) []byte {
	//SRSInfoPrefix@hash
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(srsInfoPrefix)
	buffer.WriteString(split)
	buffer.WriteString(hash)
	return buffer.Bytes()
}

// CompositeSRSListKey composite srs list key
func CompositeSRSListKey(ct string) []byte {
	//SRSInfoPrefix@hash
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(srsListPrefix)
	buffer.WriteString(split)
	buffer.WriteString(ct)
	return buffer.Bytes()
}

// CompositeAnchorNsKey compose ns anchor key
func CompositeAnchorNsKey(key string) []byte {
	return append(AnchorNsPrefix, []byte(key)...)
}

// CompositeAnchorResultKey compose anchor result key
func CompositeAnchorResultKey(txHash []byte) []byte {
	return append(AnchorExecResultPrefix, txHash...)
}

// EncodeAnchorValue construct ledger anchor node value: status list: 0 -> ToRegisgter; 1 -> Registered; 2 -> ToUnregister; 3 -> replaced
func EncodeAnchorValue(hostname string, status uint16) []byte {
	result := make([]byte, len(hostname)+4)
	binary.BigEndian.PutUint16(result[0:2], uint16(len(hostname)))
	copy(result[2:2+len(hostname)], hostname)
	binary.BigEndian.PutUint16(result[2+len(hostname):4+len(hostname)], status) // status is solidly 2 bytes len
	return result
}

// DecodeAnchorValue decode value into hostname and status
func DecodeAnchorValue(value []byte) (string, uint16) {
	hostanmeLen := binary.BigEndian.Uint16(value[0:2])
	hostnameBytes := make([]byte, hostanmeLen)
	copy(hostnameBytes, value[2:2+hostanmeLen])
	return string(hostnameBytes), binary.BigEndian.Uint16(value[2+hostanmeLen : 4+hostanmeLen])
}

// StringVersionArr implements interface sort.Sort
type StringVersionArr []string

func (sa StringVersionArr) Len() int { return len(sa) }

func (sa StringVersionArr) Swap(i, j int) { sa[i], sa[j] = sa[j], sa[i] }

func (sa StringVersionArr) Less(i, j int) bool {
	result := types.CompareVersion(sa[i], sa[j])
	if result >= 0 {
		return false
	}
	return true
}

// ArrayIntersection used to calculate the intersection of multiple version arrays.
func ArrayIntersection(ss ...[]string) []string {
	if len(ss) < 2 {
		return nil
	}

	sets := make([]set.Interface, 0, len(ss))
	for _, vc := range ss {
		aset := set.New(set.NonThreadSafe)
		for _, v := range vc {
			aset.Add(v)
		}
		sets = append(sets, aset)
	}

	all := set.Union(sets[0], sets[1], sets[1:]...)
	shared := make([]string, 0)

	all.Each(func(item interface{}) bool {
		for _, s := range sets {
			if !s.Has(item) {
				return true
			}
		}

		shared = append(shared, item.(string))

		return true
	})

	sortedShared := StringVersionArr(shared)
	sort.Sort(sortedShared)

	return sortedShared
}

// UnmarshalSetSupportedVersionInput used to unmarshal bytes to type SetSupportedVersionInput.
func UnmarshalSetSupportedVersionInput(data []byte) (types.SetSupportedVersionInput, error) {
	var svi types.SetSupportedVersionInput
	err := json.Unmarshal(data, &svi)
	return svi, err
}

// UnmarshalSystemUpgradeInput used to unmarshal bytes to type RunningVersion.
func UnmarshalSystemUpgradeInput(data []byte) (types.RunningVersion, error) {
	var rv types.RunningVersion
	err := json.Unmarshal(data, &rv)
	return rv, err
}

// UnmarshalSupportedVersions used to unmarshal bytes to type SupportedVersions.
func UnmarshalSupportedVersions(data []byte) (types.SupportedVersions, error) {
	var supportedVersions types.SupportedVersions
	err := json.Unmarshal(data, &supportedVersions)
	return supportedVersions, err
}

// UnmarshalAvailableVersion used to unmarshal bytes to type AvailableVersion.
func UnmarshalAvailableVersion(data []byte) (types.AvailableVersion, error) {
	var availableVersions types.AvailableVersion
	err := json.Unmarshal(data, &availableVersions)
	return availableVersions, err
}

// UnmarshalRunningVersion used to unmarshal bytes to type RunningVersion.
func UnmarshalRunningVersion(data []byte) (types.RunningVersion, error) {
	var runningVersions types.RunningVersion
	err := json.Unmarshal(data, &runningVersions)
	return runningVersions, err
}
