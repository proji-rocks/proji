package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestTemplate_Bucket(t *testing.T) {
	t.Parallel()

	template := Template{}
	bucket := template.Bucket()
	if bucket != bucketTemplates {
		t.Fatalf("expected bucket to be %q, got %q", bucketTemplates, bucket)
	}
}

func TestTemplateAdd_MarshalJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		template *TemplateAdd
		want     *Template
	}{
		{
			name: "marshal template #1",
			template: &TemplateAdd{
				Path:        "/path/to/template",
				UpstreamURL: stringToPointer("https://example.com"),
				Description: stringToPointer("template description"),
			},
			want: &Template{
				Path:        "/path/to/template",
				UpstreamURL: stringToPointer("https://example.com"),
				Description: stringToPointer("template description"),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			name: "marshal template #2",
			template: &TemplateAdd{
				Path:        "/path/to/template",
				UpstreamURL: stringToPointer("https://example.com"),
			},
			want: &Template{
				Path:        "/path/to/template",
				UpstreamURL: stringToPointer("https://example.com"),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			name: "marshal template #3",
			template: &TemplateAdd{
				Path:        "/path/to/template",
				Description: stringToPointer("template description"),
			},
			want: &Template{
				Path:        "/path/to/template",
				Description: stringToPointer("template description"),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		{
			name:     "marshal template #4",
			template: &TemplateAdd{},
			want: &Template{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Marshal the template into JSON.
			gotRaw, err := tc.template.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}

			// Unmarshal the JSON into a Template to make testing easier and validate that unmarshalling was also
			// successful.
			var got Template
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
			diff := cmp.Diff(tc.want, &got, cmpopts.IgnoreFields(Template{}, "ID"))
			if diff != "" {
				t.Fatalf("MarshalJSON() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
