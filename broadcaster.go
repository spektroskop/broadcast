package broadcast

type Channel chan packet

func New() Channel {
	return make(Channel, 1)
}

type packet struct {
	channel Channel
	value   interface{}
}

func wrap(value interface{}) packet {
	return packet{
		channel: New(),
		value:   value,
	}
}

func (channel Channel) send(pack packet) Channel {
	channel <- pack
	return pack.channel
}

func (channel Channel) Dispatch(value <-chan interface{}) {
	for {
		channel = channel.send(
			wrap(<-value),
		)
	}
}

func (channel Channel) Stream(into chan<- interface{}) {
	for {
		packet := <-channel
		into <- packet.value
		channel = channel.send(packet)
	}
}
