package iso8583

// Message ...
type Message struct {
	mti          *MTI
	bitmaps      []*Bitmap
	dataElements DataElements
}
