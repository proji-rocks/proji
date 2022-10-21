package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestProject_Bucket(t *testing.T) {
	t.Parallel()

	project := Project{}
	bucket := project.Bucket()
	if bucket != bucketProjects {
		t.Fatalf("expected bucket to be %q, got %s", bucketProjects, bucket)
	}
}

func TestNewProject(t *testing.T) {
	t.Parallel()

	want := &ProjectAdd{
		Path:    "/some/where",
		Name:    "test",
		Package: "tst",
	}
	got := NewProject(want.Package, want.Path, want.Name)

	if got.Path != want.Path {
		t.Fatalf("expected Path to be %q, got %q", want.Path, got.Path)
	}
	if got.Name != want.Name {
		t.Fatalf("expected Name to be %q, got %q", want.Name, got.Name)
	}
	if got.Package != want.Package {
		t.Fatalf("expected Package to be %q, got %q", want.Package, got.Package)
	}
}

func TestProjectAdd_MarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		project *ProjectAdd
		want    *Project
	}{
		{
			name: "marshal project #1",
			project: &ProjectAdd{
				Path:        "/some/where",
				Name:        "test",
				Package:     "tst",
				Description: stringToPointer("Some description bla bla bla."),
			},
			want: &Project{
				Path:        "/some/where",
				Name:        "test",
				Package:     "tst",
				Description: stringToPointer("Some description bla bla bla."),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			name: "marshal project #2",
			project: &ProjectAdd{
				Path:    "/some/where",
				Name:    "test",
				Package: "tst",
			},
			want: &Project{
				Path:      "/some/where",
				Name:      "test",
				Package:   "tst",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name:    "marshal project #3",
			project: &ProjectAdd{},
			want: &Project{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Marshal the project into JSON.
			gotRaw, err := tc.project.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}

			// Unmarshal the JSON into a Project to make testing easier and validate that unmarshalling was also
			// successful.
			var got Project
			err = json.Unmarshal(gotRaw, &got)
			if err != nil {
				t.Fatalf("UnmarshalJSON() error = %v", err)
			}

			// ID has to be generated.
			if got.ID == "" {
				t.Fatalf("ID is empty")
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

			// Compare the fields of the Package and the PackageAdd. Make cmp ignore the ID field.
			diff := cmp.Diff(tc.want, &got, cmpopts.IgnoreFields(Project{}, "ID"))
			if diff != "" {
				t.Fatalf("MarshalJSON() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
