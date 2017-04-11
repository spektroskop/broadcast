package broadcast

type Channel chan packet

func New() Channel {
	return make(Channel, 1)
}

type packet struct {
	c Channel
	v interface{}
}

func wrap(value interface{}) packet {
	return packet{New(), value}
}

func (c Channel) send(p packet) Channel {
	c <- p
	return p.c
}

func (c Channel) Dispatch(value <-chan interface{}) {
	for {
		c = c.send(
			wrap(<-value),
		)
	}
}

func (c Channel) Stream(into chan<- interface{}) {
	for {
		p := <-c
		c = c.send(p)
		into <- p.v
	}
}
