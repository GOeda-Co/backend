package stats

//Gateway layer replaces repository layer in the services where other microservices cannot be used
//Inside specific methods i do requests to my other microservices
type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr}
}

