package packet

// LayerType is a non-standard enumeration of layer instances such as Ethernet,
// or IPv4.
type LayerType uint8

const (
	LayerTypeUnknown  LayerType = 0
	LayerTypeEthernet LayerType = 1
	LayerTypeIPv4     LayerType = 2
	LayerTypeARP      LayerType = 3
)

// Layer contains a decoded layer instance, such as an Ethernet frame, or an IP
// packet.
type Layer interface {
	// Type returns the layer type.
	Type() LayerType

	// GetContents() returns the entire layer contents in bytes.
	GetContents() []byte

	// GetPayload() returns the layer payload. That is Contents()[hdrSize:]
	GetPayload() []byte
}
