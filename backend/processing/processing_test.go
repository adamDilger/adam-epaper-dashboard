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
