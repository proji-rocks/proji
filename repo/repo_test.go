package repo

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name: "Test ParseURL 1",
			args: args{URL: "https://github.com/nikoksr/proji"},
			want: &url.URL{
				Scheme: "https",
				Host:   "github.com",
				Path:   "/nikoksr/proji",
			},
			wantErr: false,
		},
		{
			name: "Test ParseURL 2",
			args: args{URL: "https://github.com/nikoksr/proji.git"},
			want: &url.URL{
				Scheme: "https",
				Host:   "github.com",
				Path:   "/nikoksr/proji",
			},
			wantErr: false,
		},
		{
			name: "Test ParseURL 3",
			args: args{URL: "gh:/nikoksr/proji/configs/test.conf"},
			want: &url.URL{
				Scheme: "https",
				Host:   "github.com",
				Path:   "/nikoksr/proji/configs/test.conf",
			},
			wantErr: false,
		},
		{
			name: "Test ParseURL 4",
			args: args{URL: "gl:/nikoksr/proji-test.git"},
			want: &url.URL{
				Scheme: "https",
				Host:   "gitlab.com",
				Path:   "/nikoksr/proji-test",
			},
			wantErr: false,
		},
		{
			name:    "Test ParseURL 5",
			args:    args{URL: ""},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURL(tt.args.URL)
			assert.Equal(t, err != nil, tt.wantErr, "ParseURL() error = %v, wantErr %v", err, tt.wantErr)
			if got == nil {
				return
			}
			assert.Equal(t, tt.want, got, "ParseURL() got = %v, want %v", got, tt.want)
		})
	}
}
