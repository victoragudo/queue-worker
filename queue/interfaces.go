package queue

type IProcessor interface {
	Process(item any) error
}
