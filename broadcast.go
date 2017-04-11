package broadcast

type channel chan packet

type packet struct {
	c channel
	v interface{}
}

func Channel() channel {
	return make(channel, 1)
}

func wrap(value interface{}) packet {
	return packet{Channel(), value}
}

func (c *channel) unwrap(p packet) interface{} {
	*c <- p
	*c = p.c
	return p.v
}

func (c channel) From(receive <-chan interface{}) {
	for {
		c.unwrap(wrap(<-receive))
	}
}

func (c channel) Into(send chan<- interface{}) {
	for {
		send <- c.unwrap(<-c)
	}
}
