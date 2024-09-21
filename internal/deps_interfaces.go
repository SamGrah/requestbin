package app

type Db interface {
	Connect() error
}

type Services interface{}

type Server interface {
	Start() error
}
