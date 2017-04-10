package broadcaster

type Channel chan Packet

func New() Channel {
	return make(Channel, 1)
}

type Packet struct {
	channel Channel
	value   interface{}
}

func NewPacket(value interface{}) Packet {
	return Packet{
		channel: New(),
		value:   value,
	}
}

func (channel Channel) send(packet Packet) Channel {
	channel <- packet
	return packet.channel
}

func (channel Channel) Dispatch(receive <-chan interface{}) {
	for {
		channel = channel.send(
			NewPacket(<-receive),
		)
	}
}

func (channel Channel) Stream(into chan<- interface{}) {
	for {
		packet := <-channel
		channel = channel.send(packet)
		into <- packet.value
	}
}
