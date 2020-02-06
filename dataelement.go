package iso8583

// Field is dataelement field identifier
type Field int

const (
	_ Field = iota
	existNextBitmap
	// TODO all field
)

// DataElement ...
type DataElement struct {
}

// DataElements ...
type DataElements map[Field]DataElement
