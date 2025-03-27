package processing_test

import (
	"epaper-dashboard/processing"
	"testing"
)

func TestConvertBoolArrayToBytes(t *testing.T) {
	tests := []struct {
		name string
		data [][]bool
		want []byte
	}{
		{
			name: "test case 1",
			data: [][]bool{
				{true, false, true, false, true, false, true, false},
				{false, true, true, true, true, true, true, false},
			},
			want: []byte{0b10101010, 0b01111110},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processing.ConvertBoolArrayToBytes(tt.data)

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ConvertBoolArrayToBytes() row [%d] = %b, want %b", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestConvertBoolArrayToBytesRLE(t *testing.T) {
	tests := []struct {
		name string
		data [][]bool
		want []uint8
	}{
		{
			name: "test case 1",
			data: [][]bool{
				{true, true, true, true, true, false, false, false},
				{false, true, true, true, true, true, true, false},
			},
			want: []uint8{
				0b10000101,
				0b00000011,
				0b00000001,
				0b10000110,
				0b00000001,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processing.ConvertBoolArrayToBytesRLE(tt.data)

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ConvertBoolArrayToBytesRLE() row [%d] = %b, want %b", i, got[i], tt.want[i])

				}
			}
		})
	}
}
