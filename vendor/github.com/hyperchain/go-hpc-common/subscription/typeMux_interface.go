package subscription

// TypeMuxInterface is the interface of the `TypeMux` type
type TypeMuxInterface interface {
	Post(ev interface{}) error
	Subscribe(types ...interface{}) Subscription
}
