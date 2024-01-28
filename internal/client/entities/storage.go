package entities

type IStorage interface {
	Init() error
	Shutdown()
}
