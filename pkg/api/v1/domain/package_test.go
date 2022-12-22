package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/nikoksr/proji/pkg/pointer"
	"github.com/pelletier/go-toml/v2"
)

func TestNewPackage(t *testing.T) {
	t.Parallel()

	type args struct {
		name string
	}
	cases := []struct {
		name string
		args args
		want *PackageAdd
	}{
		{
			name: "new package - simple",
			args: args{name: "test"},
			want: &PackageAdd{Name: "test", Label: "tst"},
		},
		{
			name: "new package - with no name",
			args: args{name: ""},
			want: &PackageAdd{Name: "Unknown", Label: "xxx"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := NewPackageWithAutoLabel(tc.args.name)
			if got == nil {
				t.Fatal("NewPackage() returned nil")
			}
			if got.Name != tc.want.Name {
				t.Fatalf("NewPackage() name got = %q, want %q", got.Name, tc.want.Name)
			}
			if got.Label != tc.want.Label {
				t.Fatalf("NewPackage() label got = %q, want %q", got.Label, tc.want.Label)
			}
		})
	}
}

func TestPackage_Bucket(t *testing.T) {
	t.Parallel()

	pkg := Package{}
	bucket := pkg.Bucket()
	if bucket != bucketPackages {
		t.Fatalf("expected bucket to be %q, got %s", bucketPackages, bucket)
	}
}

func TestPackageAdd_Marshalling(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		pkg  *PackageAdd
		want *Package
	}{
		{
			name: "marshal package #1",
			pkg:  &PackageAdd{Name: "test", Label: "tst"},
			want: &Package{Name: "test", Label: "tst", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		{
			name: "marshal package #2",
			pkg:  &PackageAdd{Name: "xxxxxx", Label: "x"},
			want: &Package{Name: "xxxxxx", Label: "x", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
		{
			name: "marshal package #3",
			pkg:  &PackageAdd{Name: "", Label: ""},
			want: &Package{Name: "", Label: "", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Marshal the package into JSON.
			gotRaw, err := tc.pkg.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}

			// Unmarshal the JSON into a Package to make testing easier and validate that unmarshalling was also
			// successful.
			var got Package
			err = json.Unmarshal(gotRaw, &got)
			if err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}

			// Normalize the timestamps to make testing easier. We cut off everything after the second. This may lead to
			// false positives, but it's good enough for our purposes. A typical false positive would be if the
			// timestamps were only off by a few milliseconds but in different seconds.
			// For example, '2022-01-01 00:00:10.999' and '2022-01-01 00:00:11.000' would result in a false positive
			// the first timestamp gets rounded to '2022-01-01 00:00:10' ant the second one gets rounded to
			// '2022-01-01 00:00:11'. This should happen rarely since the tests get executed in less than a second but
			// still something to keep in mind. A better solution would be to allow a relative tolerance for the
			// timestamps.
			tc.want.CreatedAt = tc.want.CreatedAt.Truncate(time.Second)
			tc.want.UpdatedAt = tc.want.UpdatedAt.Truncate(time.Second)
			got.CreatedAt = got.CreatedAt.Truncate(time.Second)
			got.UpdatedAt = got.UpdatedAt.Truncate(time.Second)

			// Compare the fields of the Package and the PackageAdd.
			if diff := cmp.Diff(tc.want, &got); diff != "" {
				t.Fatalf("MarshalJSON() mismatch (-want +got):\n%s", diff)
			}

			// Repeat the same process for TOML.
			gotRaw, err = tc.pkg.MarshalTOML()
			if err != nil {
				t.Fatalf("MarshalTOML() error = %v", err)
			}

			err = toml.Unmarshal(gotRaw, &got)
			if err != nil {
				t.Fatalf("UnmarshalTOML() error = %v", err)
			}

			// ... ^
			tc.want.CreatedAt = tc.want.CreatedAt.Truncate(time.Second)
			tc.want.UpdatedAt = tc.want.UpdatedAt.Truncate(time.Second)
			got.CreatedAt = got.CreatedAt.Truncate(time.Second)
			got.UpdatedAt = got.UpdatedAt.Truncate(time.Second)

			// Compare the fields of the Package and the PackageAdd.

			if diff := cmp.Diff(tc.want, &got); diff != "" {
				t.Fatalf("MarshalJSON() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPackage_Conversions(t *testing.T) {
	t.Parallel()

	type wantPackages struct {
		updated *PackageUpdate
	}
	cases := []struct {
		name string
		pkg  *Package
		want wantPackages
	}{
		{
			name: "convert package - simple",
			pkg:  &Package{Name: "test", Label: "tst"},
			want: wantPackages{
				updated: &PackageUpdate{Name: "test", Label: "tst"},
			},
		},
		{
			name: "convert package - complex",
			pkg: &Package{
				Name:        "test",
				Label:       "tst",
				UpstreamURL: pointer.To("https://github.com/user/repo/tree/branch"),
				SHA:         pointer.To("1234567890abcdef"),
				Description: pointer.To("This is a test package."),
				DirTree: &DirTree{
					Entries: []*DirEntry{
						{IsDir: true, Path: "test"},
						{IsDir: false, Path: "test/test.txt"},
						{IsDir: false, Path: "README.md", Template: &Template{ID: "123", Path: "basic.tmpl"}},
					},
				},
				Plugins: &PluginScheduler{
					Pre: []*Plugin{
						{ID: "123", Path: "basic.lua"},
						{ID: "456", Path: "advanced.lua"},
					},
					Post: []*Plugin{
						{ID: "789", Path: "basic.lua"},
						{ID: "012", Path: "advanced.lua"},
					},
				},
				CreatedAt: time.Unix(0, 0),
				UpdatedAt: time.Unix(0, 0),
			},
			want: wantPackages{
				updated: &PackageUpdate{
					Name:        "test",
					Label:       "tst",
					UpstreamURL: pointer.To("https://github.com/user/repo/tree/branch"),
					SHA:         pointer.To("1234567890abcdef"),
					Description: pointer.To("This is a test package."),
					DirTree: &DirTree{
						Entries: []*DirEntry{
							{IsDir: true, Path: "test"},
							{IsDir: false, Path: "test/test.txt"},
							{IsDir: false, Path: "README.md", Template: &Template{ID: "123", Path: "basic.tmpl"}},
						},
					},
					Plugins: &PluginScheduler{
						Pre: []*Plugin{
							{ID: "123", Path: "basic.lua"},
							{ID: "456", Path: "advanced.lua"},
						},
						Post: []*Plugin{
							{ID: "789", Path: "basic.lua"},
							{ID: "012", Path: "advanced.lua"},
						},
					},
				},
			},
		},
		{
			name: "convert package - empty",
			pkg:  &Package{},
			want: wantPackages{
				updated: &PackageUpdate{},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if diff := cmp.Diff(tc.want.updated, tc.pkg.AsUpdatable()); diff != "" {
				t.Fatalf("Package.ToUpdate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
