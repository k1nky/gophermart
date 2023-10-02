package entity

import "testing"

func TestOrderNumber_IsValid(t *testing.T) {
	tests := []struct {
		name string
		n    OrderNumber
		want bool
	}{
		{
			name: "",
			n:    "12345678903",
			want: true,
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
