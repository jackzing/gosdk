package types

// GenesisPropertyAccountParam used to create genesis property accounts
type GenesisPropertyAccountParam struct {
	EntityID       []byte      `json:"entityID"`
	CreateTime     uint64      `json:"createTime"`
	MetaData       string      `json:"metaData"`
	PropertyStatus int         `json:"status"`
	Owner          *HPCAccount `json:"owner"`
}
