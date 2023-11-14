package queue

type Processor struct {
	DoWork func(item any) error
}

func (p *Processor) Process(item any) error {
	err := p.DoWork(item)
	if err != nil {
		return err
	}
	return nil
}

func NewProcessor(callback func(item any) error) *Processor {
	return &Processor{
		DoWork: callback,
	}
}
