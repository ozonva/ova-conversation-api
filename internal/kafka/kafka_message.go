package kafka

type MessageType int

const (
	Create MessageType = iota
	Update
	Remove
)

func (mt MessageType) String() string {
	switch mt {
	case Create:
		return "Created"
	case Update:
		return "Updated"
	case Remove:
		return "Removed"
	default:
		return "Unknown message type"
	}
}

type Message struct {
	Type MessageType
	Body map[string]interface{}
}
