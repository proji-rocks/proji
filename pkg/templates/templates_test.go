package templates

import (
	"bytes"
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/google/go-cmp/cmp"
)

func TestNewEngine(t *testing.T) {
	t.Parallel()

	type args struct {
		StartTag string
		EndTag   string
	}

	cases := []struct {
		name string
		args args
		want *TemplateEngine
	}{
		{
			name: "default",
			args: args{
				StartTag: "",
				EndTag:   "",
			},
			want: &TemplateEngine{
				StartTag:     "%{{",
				EndTag:       "}}%",
				MissingKeyFn: defaultMissingKeyFn,
			},
		},
		{
			name: "custom #1",
			args: args{
				StartTag: "%{{",
				EndTag:   "}}%",
			},
			want: &TemplateEngine{
				StartTag:     "%{{",
				EndTag:       "}}%",
				MissingKeyFn: defaultMissingKeyFn,
			},
		},
		{
			name: "custom #2",
			args: args{
				StartTag: "!!",
				EndTag:   "??",
			},
			want: &TemplateEngine{
				StartTag:     "!!",
				EndTag:       "??",
				MissingKeyFn: defaultMissingKeyFn,
			},
		},
	}

	funcFilter := cmp.FilterValues(func(x, y MissingKeyFn) bool {
		return y != nil
	}, cmp.Ignore())

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := NewEngine(tc.args.StartTag, tc.args.EndTag)
			if diff := cmp.Diff(tc.want, got, funcFilter); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_normalizeKey(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
	}

	cases := []struct {
		name string
		args args
		want string
	}{
		{
			name: "dash separated",
			args: args{
				key: "project-name",
			},
			want: "projectname",
		},
		{
			name: "underscore separated",
			args: args{
				key: "project_name",
			},
			want: "projectname",
		},
		{
			name: "space separated",
			args: args{
				key: "project name",
			},
			want: "projectname",
		},
		{
			name: "uppercase",
			args: args{
				key: "PROJECTNAME",
			},
			want: "projectname",
		},
		{
			name: "empty",
			args: args{
				key: "",
			},
			want: "",
		},
		{
			name: "various",
			args: args{
				key: " PrOjEcT-nAmE_",
			},
			want: "projectname",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := normalizeKey(tc.args.key)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTemplateEngine_Parse(t *testing.T) {
	t.Parallel()

	type args struct {
		template     string
		missingKeyFn MissingKeyFn
	}

	cases := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				template: "",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "simple",
			args: args{
				template: "Hello, %{{name}}%!",
				missingKeyFn: func(key string) (string, error) {
					if key == "Name" {
						return "Proji", nil
					}
					return "", errors.Newf("unexpected key: %s", key)
				},
			},
			want:    "Hello, Proji!",
			wantErr: false,
		},
		{
			name: "advanced #1",
			args: args{
				template: "Hello, %{{name}}%! %{{message}}%",
				missingKeyFn: func(key string) (string, error) {
					if key == "Name" {
						return "Proji", nil
					} else if key == "Message" {
						return "I hope you have a great day.", nil
					}
					return "", errors.Newf("unexpected key: %s", key)
				},
			},
			want:    "Hello, Proji! I hope you have a great day.",
			wantErr: false,
		},
		{
			name: "advanced #2",
			args: args{
				template: `package %{{package}}%

import "fmt"

func Greet() {
    fmt.Println("Hello, %{{name}}%!")
}`,
				missingKeyFn: func(key string) (string, error) {
					if key == "Package" {
						return "main", nil
					} else if key == "Name" {
						return "Proji", nil
					}
					return "", errors.Newf("unexpected key: %s", key)
				},
			},
			want: `package main

import "fmt"

func Greet() {
    fmt.Println("Hello, Proji!")
}`,
			wantErr: false,
		},
		{
			name: "invalid template",
			args: args{
				template:     "Hello, %{{name}}!",
				missingKeyFn: defaultMissingKeyFn,
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			engine := NewEngine("", "")
			if engine == nil {
				t.Fatal("engine is nil")
			}

			// This is crucial to ensure that the correct values will be passed to the template engine.
			engine.MissingKeyFn = tc.args.missingKeyFn

			buf := bytes.NewBufferString("")

			err := engine.Parse(context.Background(), buf, []byte(tc.args.template))
			if (err != nil) != tc.wantErr {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.want, buf.String()); diff != "" {
				t.Fatalf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
