package hvm

import "encoding/json"

// GenAbi gen abi from abiJSON.
func GenAbi(abiJSON string) (Abi, error) {
	var abi Abi
	err := json.Unmarshal([]byte(abiJSON), &abi)
	if err != nil {
		return nil, err
	}
	return abi, nil
}
