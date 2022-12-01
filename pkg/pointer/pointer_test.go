package pointer_test

import (
	"testing"

	"github.com/nikoksr/proji/pkg/pointer"
)

func TestAsPointer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    any
		want any
	}{
		{
			name: "string",
			v:    "test",
			want: "test",
		},
		{
			name: "int",
			v:    1,
			want: 1,
		},
		{
			name: "bool",
			v:    true,
			want: true,
		},
		{
			name: "float",
			v:    1.0,
			want: 1.0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := pointer.To(tt.v)
			if got == nil {
				t.Fatalf("Expected pointer to be set")
			}
			if *got != tt.want {
				t.Errorf("AsPointer() = %v, want %v", *got, tt.want)
			}
		})
	}
}
