package network

// NewUnknown returns a new message of Unknown kind.
func NewUnknown() Message {
	return base{
		Kind: Unknown,
		Data: nil,
	}
}

// NewVoteToStart returns a new message of VoteToStart kind.
func NewVoteToStart() Message {
	return base{
		Kind: VoteToStart,
		Data: nil,
	}
}
