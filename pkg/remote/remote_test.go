package remote

import "testing"

func TestIsStatusCodeOK(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		code int
		want bool
	}{
		{
			name: "100",
			code: 100,
			want: false,
		},
		{
			name: "199",
			code: 199,
			want: false,
		},
		{
			name: "200",
			code: 200,
			want: true,
		},
		{
			name: "299",
			code: 299,
			want: true,
		},
		{
			name: "300",
			code: 300,
			want: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := IsStatusCodeOK(tc.code)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDefaultPathSkipper(t *testing.T) {
	t.Parallel()

	if DefaultPathSkipper("") {
		t.Fatalf("DefaultPathSkipper did not return false")
	}
}
