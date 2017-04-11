package broadcast

type channel chan packet

type packet struct {
	c channel
	v interface{}
}

func Channel() channel {
	return make(channel, 1)
}

func (c channel) From(receive <-chan interface{}) {
	for {
		p := packet{Channel(), <-receive}
		c <- p
		c = p.c
	}
}

func (c channel) Into(send chan<- interface{}) {
	for {
		p := <-c
		c <- p
		c = p.c
		send <- p.v
	}
}
