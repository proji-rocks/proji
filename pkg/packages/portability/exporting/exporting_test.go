package exporting

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/nikoksr/proji/pkg/api/v1/domain"
	"github.com/nikoksr/proji/pkg/packages/portability/importing"
)

func newStringPointer(s string) *string {
	return &s
}

func Test_ToConfig(t *testing.T) {
	t.Parallel()

	type args struct {
		pkg *domain.Package
		dir string
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid package",
			args: args{
				pkg: &domain.Package{
					Label:       "tst",
					Name:        "test",
					UpstreamURL: newStringPointer("https://github.com/nikoksr/proji/"),
					Description: newStringPointer("A test package."),
					DirTree: domain.DirTree{
						&domain.DirEntry{IsDir: false, Path: "docs"},
						&domain.DirEntry{IsDir: false, Path: "tests"},
						&domain.DirEntry{IsDir: true, Path: "file1.go", Template: &domain.Template{Path: "file1.go"}},
						&domain.DirEntry{IsDir: true, Path: "file2.go", Template: &domain.Template{Path: "file2.go"}},
						&domain.DirEntry{IsDir: true, Path: "file3.go", Template: &domain.Template{Path: "file3.go"}},
					},
					Plugins: &domain.PluginScheduler{
						Pre: []*domain.Plugin{
							{Path: "script1.lua"},
							{Path: "script2.lua"},
						},
						Post: []*domain.Plugin{
							{Path: "script3.lua"},
							{Path: "script4.lua"},
						},
					},
				},
				dir: "",
			},
			wantErr: false,
		},
		{
			name: "Invalid destination directory",
			args: args{
				pkg: &domain.Package{},
				dir: "/invalid/dir",
			},
			wantErr: true,
		},
		{
			name: "Invalid package (nil)",
			args: args{
				pkg: nil,
				dir: "",
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Empty dir will create a temporary file.
			configPath, err := ToConfig(context.Background(), tc.args.pkg, tc.args.dir)
			if err != nil {
				if !tc.wantErr {
					t.Fatalf("ToConfig() error = %v, wantErr %v", err, tc.wantErr)
				}

				return // If an error occurred, we don't need to check the config path.
			}

			if configPath == "" {
				t.Fatalf("ToConfig() config path is empty")
			}
			defer func() { _ = os.Remove(configPath) }() // Clean after yourself, dirty.

			// Check if config directory is as expected.
			configDir := filepath.Dir(configPath)

			if tc.args.dir != "" {
				// If a directory was specified, it should be the same as the config directory.
				if configDir != tc.args.dir {
					t.Fatalf(
						"ToConfig() config path is not in the expected directory (-want +got):\n%s",
						cmp.Diff(tc.args.dir, filepath.Dir(configPath)),
					)
				}
			} else {
				// This should create a config file in a temporary directory.
				if !strings.HasPrefix(configDir, os.TempDir()) {
					t.Fatalf("ToConfig() config path is not in a temporary directory")
				}
			}

			// Try to load package from config.
			//
			// TODO: Is this okay? Using a different internal function for verifying the functionality of another
			//       internal function?
			pkgGot, err := importing.LocalPackage(context.Background(), configPath)
			if err != nil {
				t.Fatalf("Failed to open config created by ToConfig(): %v", err)
			}

			// Compare fields manually; didn't find a simple solution to do this with gocmp. This is good enough for now,
			// just verbose.
			if pkgGot.Label != tc.args.pkg.Label {
				t.Fatalf("ToConfig() label = %v, want %v", pkgGot.Label, tc.args.pkg.Label)
			}
			if pkgGot.Name != tc.args.pkg.Name {
				t.Fatalf("ToConfig() name = %v, want %v", pkgGot.Name, tc.args.pkg.Name)
			}

			diff := cmp.Diff(pkgGot.UpstreamURL, tc.args.pkg.UpstreamURL)
			if diff != "" {
				t.Fatalf("ToConfig() upstreamURL mismatch (-want +got):\n%s", diff)
			}

			diff = cmp.Diff(pkgGot.Description, tc.args.pkg.Description)
			if diff != "" {
				t.Fatalf("ToConfig() description mismatch (-want +got):\n%s", diff)
			}

			diff = cmp.Diff(pkgGot.DirTree, tc.args.pkg.DirTree)
			if diff != "" {
				t.Fatalf("ToConfig() dirTree mismatch (-want +got):\n%s", diff)
			}

			diff = cmp.Diff(pkgGot.Plugins, tc.args.pkg.Plugins)
			if diff != "" {
				t.Fatalf("ToConfig() plugins mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
