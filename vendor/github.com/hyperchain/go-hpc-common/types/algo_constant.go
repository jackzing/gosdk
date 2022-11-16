package types

import "github.com/meshplus/crypto"

const (
	//DefaultAlgo default algo
	DefaultAlgo = 0xffffffff
	//HashMethod use to check algo
	HashMethod = "hash"
	//EncryptMethod use to check algo
	EncryptMethod = "encrypt"
	//FileLogEnc use to GetAlgoFromEncrypt
	FileLogEnc = "filelog"
	//MultiCacheEnc use to Multicache
	MultiCacheEnc = "multicache"
	//MultiHash multi hash encoding
	MultiHash = "MULTIHASH"
	//MultiCipher multi cipher encoding
	MultiCipher = "MULTICIPHER"
	//HexCode hex encoding
	HexCode = "HEX"
	//RawCode raw encoding
	RawCode = "RAW"
	//encVerLen length of encode version
	encVerLen = 4
)

const (
	//SHA2_224 hash algo
	SHA2_224 = "SHA2_224"
	//SHA2_256 hash algo
	SHA2_256 = "SHA2_256"
	//SHA2_384 hash algo
	SHA2_384 = "SHA2_384"
	//SHA2_512 hash algo
	SHA2_512 = "SHA2_512"
	//SHA3_224 hash algo
	SHA3_224 = "SHA3_224"
	//SHA3_256 hash algo
	SHA3_256 = "SHA3_256"
	//SHA3_384 hash algo
	SHA3_384 = "SHA3_384"
	//SHA3_512 hash algo
	SHA3_512 = "SHA3_512"
	//KECCAK_224 hash algo
	KECCAK_224 = "KECCAK_224"
	//KECCAK_256 hash algo
	KECCAK_256 = "KECCAK_256"
	//KECCAK_384 hash algo
	KECCAK_384 = "KECCAK_384"
	//KECCAK_512 hash algo
	KECCAK_512 = "KECCAK_512"
	//SM3 hash algo
	SM3 = "SM3"
	//SDHash algo
	SDHash = "SELF_DEFINED_HASH"

	//SM4_CBC encrypt algo
	SM4_CBC = "SM4_CBC"
	//AES_CBC encrypt algo
	AES_CBC = "AES_CBC"
	//DES3_CBC encrypt algo
	DES3_CBC = "3DES_CBC"
	//TEE encrypt algo
	TEE = "TEE"
	//SDCrypt algo
	SDCrypt = "SELF_DEFINED_CRYPTO"
)

//AlgoSet algo set
type AlgoSet struct {
	HashAlgo    string `json:"hash_algo,omitempty"`
	HashText    string `json:"hash_text,omitempty"`
	EncryptAlgo string `json:"encrypt_algo,omitempty"`
	EncryptText string `json:"encrypt_text,omitempty"`
}

//BlockAlgo block number--> algo
type BlockAlgo struct {
	BlockNumber int64
	AlgoSet     AlgoSet
}

//BlockAlgoSets array of block algo
type BlockAlgoSets []BlockAlgo

//algoVersion Compatible downward
type algoVersion struct {
	algo    string
	version TXVersion
}

//DefaultAlgoSet defaultAlgoSet
var DefaultAlgoSet = AlgoSet{
	HashAlgo:    KECCAK_256,
	HashText:    HexCode,
	EncryptAlgo: AES_CBC,
	EncryptText: RawCode,
}

//hashVersionSet algo --> tx version
var hashVersionSet = []algoVersion{
	{
		algo:    SHA2_224,
		version: TxVersion35,
	},
	{
		algo:    SHA2_256,
		version: TxVersion35,
	},
	{
		algo:    SHA2_384,
		version: TxVersion35,
	},
	{
		algo:    SHA2_512,
		version: TxVersion35,
	},
	{
		algo:    SHA3_224,
		version: TxVersion35,
	},
	{
		algo:    SHA3_256,
		version: TxVersion35,
	},
	{
		algo:    SHA3_384,
		version: TxVersion35,
	},
	{
		algo:    SHA3_512,
		version: TxVersion35,
	},
	{
		algo:    KECCAK_224,
		version: TxVersion35,
	},
	{
		algo:    KECCAK_256,
		version: TxVersion35,
	},
	{
		algo:    KECCAK_384,
		version: TxVersion35,
	},
	{
		algo:    KECCAK_512,
		version: TxVersion35,
	},
	{
		algo:    SM3,
		version: TxVersion35,
	},
	{
		algo:    SDHash,
		version: TxVersion41,
	},
}

//encryptVersionSet algo --> tx version
var encryptVersionSet = []algoVersion{
	{
		algo:    SM4_CBC,
		version: TxVersion35,
	},
	{
		algo:    AES_CBC,
		version: TxVersion35,
	},
	{
		algo:    DES3_CBC,
		version: TxVersion35,
	},
	{
		algo:    TEE,
		version: TxVersion40,
	},
	{
		algo:    SDCrypt,
		version: TxVersion41,
	},
}

//algoToInt algo(string) convert int
var algoToInt = map[string]int{
	SHA2_224:   crypto.SHA2_224,
	SHA2_256:   crypto.SHA2_256,
	SHA2_384:   crypto.SHA2_384,
	SHA2_512:   crypto.SHA2_512,
	SHA3_224:   crypto.SHA3_224,
	SHA3_256:   crypto.SHA3_256,
	SHA3_384:   crypto.SHA3_384,
	SHA3_512:   crypto.SHA3_512,
	KECCAK_224: crypto.KECCAK_224,
	KECCAK_256: crypto.KECCAK_256,
	KECCAK_384: crypto.KECCAK_384,
	KECCAK_512: crypto.KECCAK_512,
	SM3:        crypto.SM3,
	SDHash:     crypto.SelfDefinedHash,

	SM4_CBC:  crypto.Sm4 | crypto.CBC,
	AES_CBC:  crypto.Aes | crypto.CBC,
	DES3_CBC: crypto.Des3 | crypto.CBC,
	TEE:      crypto.TEE,
	SDCrypt:  crypto.SelfDefinedCrypt,
}

var encMagic = []byte{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
	0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
	0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
}
