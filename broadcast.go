package broadcast

import (
	"context"
)

type (
	Receiver struct {
		ctx         context.Context
		destination chan interface{}
		source      Channel
	}

	Channel struct {
		packets chan Packet
		listen  chan Receiver
	}

	Packet struct {
		packets chan Packet
		value   interface{}
	}
)

func open() chan Packet {
	return make(chan Packet, 1)
}

func wrap(packets chan Packet, value interface{}) Packet {
	return Packet{packets, value}
}

func (c *Channel) unwrap(p Packet) interface{} {
	c.packets <- p
	c.packets = p.packets
	return p.value
}

func (c Channel) receive(r Receiver) {
	for {
		select {
		case <-r.ctx.Done():
			return
		case p := <-c.packets:
			r.destination <- c.unwrap(p)
		}
	}
}

func New() Channel {
	return Channel{
		packets: open(),
		listen:  make(chan Receiver),
	}
}

func (c Channel) From(ctx context.Context, source <-chan interface{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case value := <-source:
			c.unwrap(wrap(open(), value))
		case receiver := <-c.listen:
			go c.receive(receiver)
		}
	}
}

func (c Channel) Into(ctx context.Context, d chan interface{}) {
	c.listen <- Receiver{ctx: ctx, destination: d}
}

func (c Channel) Receiver(d chan interface{}) Receiver {
	return Receiver{
		ctx:         context.Background(),
		destination: d,
		source:      c,
	}
}

func (r Receiver) Listen(ctx context.Context) {
	r.ctx = ctx
	r.source.listen <- r
}

func (r Receiver) Destination() <-chan interface{} {
	return r.destination
}
