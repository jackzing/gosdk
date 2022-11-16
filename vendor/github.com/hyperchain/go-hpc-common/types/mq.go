package types

// RegisterMeta register param meta
type RegisterMeta struct {
	// queue related
	QueueName  string   `json:"queueName"`
	EventTypes []string `json:"routingKeys"`
	// self info
	From      string `json:"from"`
	Signature string `json:"signature"`
	// block criteria
	IsVerbose bool `json:"isVerbose"`
	// vm log criteria
	FromBlock string    `json:"fromBlock"`
	ToBlock   string    `json:"toBlock"`
	Addresses []Address `json:"addresses"`
	Topics    [][]Hash  `json:"topics"`
	// strict mode and delay push
	Delay bool `json:"delay"`
	// exception criteria
	//Modules        []string `json:"modules"`
	//ModulesExclude []string `json:"modules_exclude"`
	//SubType        []string `json:"subtypes"`
	//SubTypeExclude []string `json:"subtypes_exclude"`
	//Code           []int    `json:"error_codes"`
	//CodeExclude    []int    `json:"error_codes_exclude"`
}
