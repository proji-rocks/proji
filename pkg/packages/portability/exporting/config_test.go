package exporting

import (
	"bytes"
	"context"
	"os"
	"testing"
)

func Test_write(t *testing.T) {
	t.Parallel()

	type args struct {
		validFile bool
		data      *bytes.Buffer
	}

	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Nil file",
			args: args{
				validFile: false,
				data:      nil,
			},
			wantErr: true,
		},
		{
			name: "Nil data",
			args: args{
				validFile: true,
				data:      nil,
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var file *os.File
			var err error
			if tc.args.validFile {
				file, err = os.CreateTemp("", "test_proji.config.*.json")
				if err != nil {
					t.Fatalf("failed to create temporary file: %v", err)
				}
				defer func() { _ = file.Close() }()
			}

			if err := write(context.Background(), file, tc.args.data); (err != nil) != tc.wantErr {
				t.Errorf("write() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
