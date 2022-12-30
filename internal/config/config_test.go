package config

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/pflag"
)

func Test_load(t *testing.T) {
	t.Parallel()

	type args struct {
		path  string
		flags *pflag.FlagSet
	}
	cases := []struct {
		name      string
		args      args
		want      *Config
		isValid   bool
		wantErr   bool
		doCleanUp bool
	}{
		{
			name: "Load valid config #1",
			args: args{
				path:  "./testdata/proji_valid_conf_1.toml",
				flags: pflag.NewFlagSet("test", pflag.ContinueOnError),
			},
			want: &Config{
				Database: Database{
					DSN: "/home/user_a/.config/proji/data/proji.db",
				},
				Import: Import{
					Exclude: "^(.git|.env|.idea|.vscode)$",
				},
				Monitoring: Monitoring{
					Sentry: Sentry{
						Enabled: false,
					},
				},
				System: System{
					TextEditor: "vim",
				},
			},
			isValid:   true,
			wantErr:   false,
			doCleanUp: false,
		},
		{
			name: "Load valid config #2",
			args: args{
				path:  "./testdata/proji_valid_conf_2.toml",
				flags: pflag.NewFlagSet("test", pflag.ContinueOnError),
			},
			want: &Config{
				Database: Database{
					DSN: "/home/user_b/.local/share/proji.db",
				},
				Import: Import{
					Exclude: "",
				},
				Monitoring: Monitoring{
					Sentry: Sentry{
						Enabled: false,
					},
				},
				System: System{
					TextEditor: "",
				},
			},
			isValid:   true,
			wantErr:   false,
			doCleanUp: false,
		},
		{
			name: "Load non-existent config; creates default config",
			args: args{
				path: filepath.Join(os.TempDir(), "proji-test", "Config_load", defaultConfigFile),
			},
			want: &Config{
				Database: Database{
					DSN: filepath.Join(os.TempDir(), "proji-test", "Config_load", defaultDataDir, defaultDatabaseFile),
				},
				Import: Import{
					Exclude: defaultExcludePattern,
				},
				Monitoring: Monitoring{
					Sentry: Sentry{
						Enabled: defaultSentryState,
					},
				},
				System: System{
					TextEditor: "",
				},
			},
			isValid:   true,
			wantErr:   false,
			doCleanUp: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Load config from examples folder
			conf, err := load(context.Background(), tc.args.path, tc.args.flags)
			if (err != nil) != tc.wantErr {
				t.Fatalf("load() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if tc.wantErr {
				return // No need to continue
			}

			// At this point, config should never be nil
			if conf == nil {
				t.Fatal("load() returned nil config")
				return
			}

			// Check if config is as expected
			if diff := cmp.Diff(tc.want, conf, cmpopts.IgnoreFields(Config{}, "provider")); diff != "" {
				t.Fatalf("load() returned unexpected config (-want +got):\n%s", diff)
			}

			// Check if config is valid
			if err = conf.Validate(); tc.isValid == (err != nil) {
				t.Fatalf("Validate() failed unexpectedly: %v", err)
			}

			if tc.doCleanUp {
				dir := filepath.Dir(tc.args.path)
				err = os.RemoveAll(dir)
				if err != nil {
					t.Errorf("Failed to cleanup test config at %q: %v", dir, err)
				}
			}
		})
	}
}

func TestLoad(t *testing.T) {
	t.Parallel()

	// Compared to Test_load, this test focuses on the Singleton pattern that Config implements. We're loading the same
	// config twice and check if the same instance is returned.
	confPath := filepath.Join(os.TempDir(), "proji-test", "ConfigLoad", defaultConfigFile)
	defer func() {
		dir := filepath.Dir(confPath)
		err := os.RemoveAll(dir)
		if err != nil {
			t.Errorf("Failed to cleanup test config at %q: %v", dir, err)
		}
	}()

	// Load config; this should initialize the singleton
	conf1, err := Load(context.Background(), confPath, nil)
	if err != nil {
		t.Fatalf("Load() failed unexpectedly: %v", err)
	}
	if conf1 == nil {
		t.Fatal("Load() returned nil config")
	}

	// Load config again; this should return the same instance
	conf2, err := Load(context.Background(), confPath, nil)
	if err != nil {
		t.Fatalf("Load() failed unexpectedly: %v", err)
	}
	if conf2 == nil {
		t.Fatal("Load() returned nil config")
	}

	// Check if both instances are the same
	if conf1 != conf2 {
		t.Fatal("Load() returned different instances")
	}
}

func Test_defaultConfigPath(t *testing.T) {
	t.Parallel()

	// This test is a bit tricky because we need to make sure that the default config path is correct. We can't just
	// compare it to a constant because the default config path depends on the OS.
	// These tests will later be run on different OSes, so we can't just hardcode the expected path. Only execute
	// the path check when we're on the OS that we expect.
	path, err := defaultConfigPath()

	switch runtime.GOOS {
	case "darwin", "freebsd", "linux", "netbsd", "openbsd":
		// Verify config path for *nix systems
		if err != nil {
			t.Fatalf("defaultConfigPath() failed unexpectedly: %v", err)
		}
		if path != filepath.Join(os.Getenv("HOME"), ".config", "proji", defaultConfigFile) {
			t.Fatalf("defaultConfigPath() returned unexpected path: %q", path)
		}

		// We're on a *nix system, so let's check if the getWindowsUserConfigPath() function returns an error
		_, err = getWindowsUserConfigPath()
		if err == nil {
			t.Fatalf("getWindowsUserConfigPath() returned no error even though we're on a *nix system (%q)", runtime.GOOS)
		}

	case "windows":
		// Verify config path for Windows
		if err != nil {
			t.Fatalf("defaultConfigPath() failed unexpectedly: %v", err)
		}
		if path != filepath.Join(os.Getenv("LOCALAPPDATA"), "proji", defaultConfigFile) {
			t.Fatalf("defaultConfigPath() returned unexpected path: %q", path)
		}

		// We're on Windows, so let's check if the getUnixUserConfigPath() function returns an error
		_, err = getUnixUserConfigPath()
		if err == nil {
			t.Fatal("getUnixUserConfigPath() returned no error even though we're on Windows")
		}
	default:
		if err == nil {
			t.Fatal("defaultConfigPath() was expected to return an error for an unsupported OS")
		}
	}
}

func TestConfig_BaseDir(t *testing.T) {
	t.Parallel()

	type args struct {
		path string
	}
	cases := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default",
			args: args{
				path: "/home/user/.config/proji.toml",
			},
			want: filepath.Join("/home", "user", ".config"),
		},
		{
			name: "complex",
			args: args{
				path: "/home/user/.config/proji/this/is/a/very/long/path/to/a/config/file.toml",
			},
			want: filepath.Join("/home", "user", ".config", "proji", "this", "is", "a", "very", "long", "path", "to", "a", "config"),
		},
		{
			name: "current dir",
			args: args{
				path: "./proji.toml",
			},
			want: ".",
		},
		// TODO: Add Windows tests
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			conf := &Config{
				provider: newProvider(tc.args.path),
			}
			if got := conf.BaseDir(); got != filepath.FromSlash(tc.want) {
				t.Fatalf("BaseDir() returned unexpected path: %q;\nbase directory of %q should be %q", got, tc.args.path, tc.want)
			}
		})
	}
}
