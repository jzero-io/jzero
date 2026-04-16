package progress

// Type represents the type of progress message
type Type int

const (
	// TypeFile indicates a file is being generated
	TypeFile Type = iota

	// TypeDebug indicates debug information
	TypeDebug

	// TypeError indicates an error occurred with a file
	TypeError
)

// Message represents a message sent through the progress channel
type Message struct {
	Type  Type
	Value string
}

// NewFile creates a new file generation progress message
func NewFile(file string) Message {
	return Message{
		Type:  TypeFile,
		Value: file,
	}
}

// NewDebug creates a new debug information progress message
func NewDebug(debug string) Message {
	return Message{
		Type:  TypeDebug,
		Value: debug,
	}
}
