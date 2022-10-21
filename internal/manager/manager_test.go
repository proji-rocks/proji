package manager

import (
	"context"
	"testing"

	bolt "go.etcd.io/bbolt"

	"github.com/nikoksr/proji/internal/config"
	database "github.com/nikoksr/proji/pkg/database/bolt"
)

func TestNewPackageManager(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		address string
		db      *database.DB
		auth    *config.Auth
	}
	cases := []struct {
		name     string
		args     args
		wantType string
		wantErr  bool
	}{
		{
			name:     "local package manager with no database",
			args:     args{ctx: context.Background(), address: "", db: nil, auth: nil},
			wantType: "local",
			wantErr:  true,
		},
		{
			name:     "local package manager with no auth",
			args:     args{ctx: context.Background(), address: "", db: &database.DB{Core: &bolt.DB{}}, auth: nil},
			wantType: "local",
			wantErr:  false,
		},
		{
			name:     "remote package manager with fully qualified address",
			args:     args{ctx: context.Background(), address: "http://localhost:8080", db: nil, auth: nil},
			wantType: "remote",
			wantErr:  false,
		},
		{
			name:     "remote package manager with malformed address",
			args:     args{ctx: context.Background(), address: "htt://localhost:8080", db: nil, auth: nil},
			wantType: "remote",
			wantErr:  true,
		},
		{
			name:     "remote package manager with short address",
			args:     args{ctx: context.Background(), address: "localhost:8080", db: nil, auth: nil},
			wantType: "remote",
			wantErr:  false,
		},
		{
			name:     "remote package manager with port only as address",
			args:     args{ctx: context.Background(), address: ":8080", db: nil, auth: nil},
			wantType: "remote",
			wantErr:  false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewPackageManager(tc.args.ctx, tc.args.address, tc.args.db, tc.args.auth)
			if (err != nil) != tc.wantErr {
				t.Fatalf("NewPackageManager() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err == nil && got == nil {
				t.Fatal("NewPackageManager() returned nil manager")
			}
			if got != nil && tc.wantType != got.String() {
				t.Fatalf("NewPackageManager() type got = %v, want %v", got, tc.wantType)
			}
		})
	}
}

func TestNewProjectManager(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		db  *database.DB
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "project manager with no database",
			args:    args{ctx: context.Background(), db: nil},
			wantErr: true,
		},
		{
			name:    "project manager with database",
			args:    args{ctx: context.Background(), db: &database.DB{Core: &bolt.DB{}}},
			wantErr: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewProjectManager(tc.args.ctx, tc.args.db)
			if (err != nil) != tc.wantErr {
				t.Fatalf("NewProjectManager() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err == nil && got == nil {
				t.Fatal("NewProjectManager() returned nil manager")
			}
		})
	}
}
