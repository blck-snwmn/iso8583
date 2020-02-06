package iso8583

// MTI ...
type MTI struct {
	version         Version
	messageClass    MessageClass
	messageSubClass MessageSubClass
}

// Version ...
type Version string

// MessageClass ...
type MessageClass string

// MessageSubClass ...
type MessageSubClass string
