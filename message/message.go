package message

import (
	"encoding/json"
	"net"
	"time"
)

type Message struct {
	Data        []byte
	Destination string
	ClientIP    net.IP
	Timestamp   time.Time
}

func (m *Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func Deserialize(data []byte) (*Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
