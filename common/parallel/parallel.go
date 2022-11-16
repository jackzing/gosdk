package parallel

// ParallelLevel parallel level.
type ParallelLevel int

const (
	// GlobalMutex global mutex.
	GlobalMutex ParallelLevel = 0
	// ContractMutex contract mutex.
	ContractMutex ParallelLevel = 1
	// FieldMutex field mutex.
	FieldMutex ParallelLevel = 2
)
