package toggle

type Toggle struct {
	C chan struct{}
}

func New(ch chan *Toggle) *Toggle {
	tg := &Toggle{
		C: make(chan struct{}),
	}
	ch <- tg
	return tg
}

func (t *Toggle) Toggle() {
	t.C <- struct{}{}
}

func (t *Toggle) Wait() {
	<-t.C
}
