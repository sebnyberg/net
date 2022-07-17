package packet

type TCP struct {
	Payload []byte
}

func (e *TCP) Unmarshal(data []byte) error {
	return nil
}
