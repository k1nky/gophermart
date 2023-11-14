package order

import "testing"

func TestOrderNumber_IsValid(t *testing.T) {
	tests := []struct {
		name string
		n    OrderNumber
		want bool
	}{
		{
			name: "Valid",
			n:    "12345678903",
			want: true,
		},
		{
			name: "Invalid",
			n:    "12345678901",
			want: false,
		},
		{
			name: "Incorrect symbols",
			n:    "12345678903a",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.IsValid(); got != tt.want {
				t.Errorf("OrderNumber.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
