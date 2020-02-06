package iso8583

import (
	"reflect"
	"testing"
)

func TestParser_parseBitmap(t *testing.T) {
	tests := []struct {
		name    string
		p       *Parser
		want    []*Bitmap
		wantErr bool
	}{
		{"all 0", newParser([]byte{0, 0, 0, 0, 0, 0, 0, 0}),
			[]*Bitmap{
				&Bitmap{
					[]byte{0, 0, 0, 0, 0, 0, 0, 0},
					[]Field{},
				},
			}, false},
		{"exist second bitmap", newParser([]byte{0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
			[]*Bitmap{
				&Bitmap{
					[]byte{0x80, 0, 0, 0, 0, 0, 0, 0},
					[]Field{Field(1)},
				},
				&Bitmap{
					[]byte{0, 0, 0, 0, 0, 0, 0, 0},
					[]Field{},
				},
			}, false},
		{"0100 0000 ...", newParser([]byte{0x40, 0, 0, 0, 0, 0, 0, 0}),
			[]*Bitmap{
				&Bitmap{
					[]byte{0x40, 0, 0, 0, 0, 0, 0, 0},
					[]Field{Field(2)},
				},
			}, false},
		{"00000000  0100 0000 ...", newParser([]byte{0, 0x40, 0, 0, 0, 0, 0, 0}),
			[]*Bitmap{
				&Bitmap{
					[]byte{0, 0x40, 0, 0, 0, 0, 0, 0},
					[]Field{Field(10)},
				},
			}, false},
		{"...  0000 0001", newParser([]byte{0, 0, 0, 0, 0, 0, 0, 1}),
			[]*Bitmap{
				&Bitmap{
					[]byte{0, 0, 0, 0, 0, 0, 0, 1},
					[]Field{Field(64)},
				},
			}, false},
		{"...  0000 0001  ...", newParser([]byte{0, 0, 0, 0, 0, 0, 1, 0}),
			[]*Bitmap{
				&Bitmap{
					[]byte{0, 0, 0, 0, 0, 0, 1, 0},
					[]Field{Field(56)},
				},
			}, false},

		{"multi field", newParser([]byte{0x76, 0, 0, 0, 0, 0, 1, 1}),
			[]*Bitmap{
				&Bitmap{
					[]byte{0x76, 0, 0, 0, 0, 0, 1, 1},
					[]Field{Field(2), Field(3), Field(4), Field(6), Field(7), Field(56), Field(64)},
				},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.parseBitmap()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseBitmap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseBitmap() = %v, want %v", got, tt.want)
			}
		})
	}
}
