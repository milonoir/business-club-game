package message

const (
	// separator is a delimiter character to separate fields in a message.
	separator = ":"
)

// Message defines the interface that all kinds of messages have to implement.
type Message interface {
	Type() Kind
	Payload() any
}

// NewUnknown returns a new message of Unknown kind.
func NewUnknown() Message {
	return base{
		Kind: Unknown,
		Data: nil,
	}
}
