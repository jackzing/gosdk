package types

// API describes the set of methods offered over the RPC interface.
type API struct {
	Svcname string      // service name
	Version string      // api version
	Service interface{} // api service instance which holds the methods
	Public  bool        // indication if the methods must be considered safe for public use
}

// OuterAPIs defines a method to return API structure which is
// defined by outer service to offer outer service.
type OuterAPIs interface {
	APIs() *API
}
