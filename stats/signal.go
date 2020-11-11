package stats

// Definitions of signals used to communicate to the gosivy agents.

const (
	// SignalMeta reports Go process metadata.
	SignalMeta = byte(0x1)

	// SignalStats reports Go process stats.
	SignalStats = byte(0x2)

	// Delimiter indicates to complete the writing.
	Delimiter = '\n'
)
