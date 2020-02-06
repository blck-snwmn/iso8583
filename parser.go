package iso8583

const mtiLen = 2
const bitmapLen = 8 // read 64bit = 8byte
const bitLenPerByte = 8

type dataElementFn func() (DataElement, error)

// Parser ...
type Parser struct {
	input               []byte
	position            int
	parseDataElementFns map[Field]dataElementFn
}

func (p *Parser) parseMti() (*MTI, error) {
	_ = p.input[p.position : p.position+mtiLen]
	p.position += mtiLen
	// TODO parse Version, MessageClass, MessageSubClass
	return &MTI{}, nil
}

func (p *Parser) parseBitmap() ([]*Bitmap, error) {
	bitmaps := []*Bitmap{}
	// TODO bimap max 3?
	for {
		// read 64 bit
		// TODO check length
		bitmapByte := p.input[p.position : p.position+bitmapLen]
		p.position += bitmapLen

		bitmap := Bitmap{}
		bitmap.original = bitmapByte
		// TODO in advance, comfirm to alocate 64 size
		bitmap.fields = []Field{}

		// read one byte at a time
		for i := 0; i < bitmapLen; i++ {
			// read one bit at a time
			for j := 0; j < bitLenPerByte; j++ {
				// if Nth bit from left set 1
				// Nth field exists on DataElement
				if bitmap.original[i]<<j&0x80 == 0x80 {
					bitmap.fields = append(bitmap.fields, Field(8*i+j+1))
				}
			}
		}

		bitmaps = append(bitmaps, &bitmap)

		// if 1th bit from left doesn't set 1
		// next 64bit is not Bitmap
		if bitmap.original[0]&0x80 != 0x80 {
			break
		}
	}
	return bitmaps, nil
}

func (p *Parser) parseDataElement(bitmaps []*Bitmap) (DataElements, error) {
	des := DataElements{}
	for _, bm := range bitmaps {
		for _, f := range bm.fields {
			fn := p.parseDataElementFns[f]
			de, err := fn()
			if err != nil {
				return nil, err
			}
			des[f] = de
		}
	}
	return DataElements{}, nil
}

func (p *Parser) parse() (*Message, error) {
	mti, err := p.parseMti()
	if err != nil {
		return nil, err
	}
	bmps, err := p.parseBitmap()
	if err != nil {
		return nil, err
	}
	des, err := p.parseDataElement(bmps)
	if err != nil {
		return nil, err
	}
	return &Message{mti, bmps, des}, nil
}

func newParser(input []byte) *Parser {
	p := &Parser{input: input, position: 0}
	// ここでDataElementのパース関数が分かるようなら登録する
	// 登録用の関数があってもいいかも
	return p
}
