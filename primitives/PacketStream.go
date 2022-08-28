package primitives

import (
	"encoding/binary"
	"io"
	"log"
)

// packetStream is a wrapper for a socket connection for easier uses.
type packetStream struct {
	stream  io.ReadWriteCloser
	packets chan *[]byte
}

// newPacketStream is the constructor.
func newPacketStream(stream io.ReadWriteCloser) *packetStream {
	wrapper := packetStream{stream, make(chan *[]byte)}
	wrapper.readPackets()

	return &wrapper
}

// Continually processes events from the stream.
func (w *packetStream) readPackets() {
	var length uint32

	go func() {
		for {

			err := binary.Read(w.stream, binary.BigEndian, &length)
			if err != nil {
				log.Printf("Failed to read packet length: %s", err)
				return
			}

			if length > 0 {
				packet := make([]byte, length)

				i, err := w.stream.Read(packet)
				if err != nil {
					log.Printf("Failed to read packet: %s", err)
					return
				}

				if i != int(length) {
					log.Printf("Invalid packet size. Wanted: %d Read: %d", length, i)
					return
				}
				w.packets <- &packet
			}

		}
	}()
}

func (w *packetStream) read() *[]byte {
	return <-w.packets
}

// Sends events to the stream to be read.
func (w *packetStream) write(data *[]byte) (int, error) {

	err := binary.Write(w.stream, binary.BigEndian, uint32(len(*data)))

	if err != nil {
		log.Printf("Failed to write packet length %d. error:%s", len(*data), err)
		return 0, err
	}

	return w.stream.Write(*data)
}
