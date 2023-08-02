package network

// NewUnknown returns a new message of Unknown kind.
func NewUnknown() Message {
	return base{
		Kind: Unknown,
		Data: nil,
	}
}
