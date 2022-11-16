package msp

//CA,Cert
const (
	RootCA  = "root"
	ECert   = "ecert"
	RCert   = "rcert"
	SDKCert = "sdkcert"
	SDKCA   = "sdkca"

	FLATO       = "flato"
	Classic     = "classic"
	Distributed = "distributed"

	NoneKeyStore = "without key store"

	//Default used for the second parameter of GetIdentities, to get the default identity
	Default = "default"
)

//LengthOfSessionKey The Length Of the msp.SessionKey
const LengthOfSessionKey = 32

//EmptyExt empty ext
const EmptyExt = "e30="

//DummyIDName dummy ca name
const DummyIDName = "O=dummyID"
