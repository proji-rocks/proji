package remote

import (
	"context"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseRepoURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		repoURL string
		want    string
		wantErr bool
	}{
		{
			name:    "Perfect GitHub URL",
			repoURL: "https://github.com/nikoksr/proji",
			want:    "https://github.com/nikoksr/proji",
			wantErr: false,
		},
		{
			name:    "URL with .git suffix",
			repoURL: "https://github.com/nikoksr/proji.git",
			want:    "https://github.com/nikoksr/proji",
			wantErr: false,
		},
		{
			name:    "URL with .git suffix spaces",
			repoURL: "  https://github.com/nikoksr/proji.git   ",
			want:    "https://github.com/nikoksr/proji",
			wantErr: false,
		},
		{
			name:    "URL with trailing slash",
			repoURL: "  https://github.com/nikoksr/proji/",
			want:    "https://github.com/nikoksr/proji",
			wantErr: false,
		},
		{
			name:    "GitHub abbreviated URL",
			repoURL: "gh:nikoksr/proji",
			want:    "https://github.com/nikoksr/proji",
			wantErr: false,
		},
		{
			name:    "GitLab abbreviated URL",
			repoURL: "gl:nikoksr/proji",
			want:    "https://gitlab.com/nikoksr/proji",
			wantErr: false,
		},
		{
			name:    "GitHub URL with regular branch",
			repoURL: "https://github.com/nikoksr/proji/tree/dev",
			want:    "https://github.com/nikoksr/proji/tree/dev",
			wantErr: false,
		},
		{
			name:    "GitLab URL with regular branch",
			repoURL: "https://gitlab.com/nikoksr/proji/-/tree/dev",
			want:    "https://gitlab.com/nikoksr/proji/-/tree/dev",
			wantErr: false,
		},
		{
			name:    "GitHub URL with @branch",
			repoURL: "https://github.com/nikoksr/proji@dev",
			want:    "https://github.com/nikoksr/proji/tree/dev",
			wantErr: false,
		},
		{
			name:    "GitHub URL with invalid @branch",
			repoURL: "https://github.com/nikoksr/proji@dev@main",
			want:    "",
			wantErr: true,
		},
		{
			name:    "GitLab URL with @branch",
			repoURL: "https://gitlab.com/nikoksr/proji@dev",
			want:    "https://gitlab.com/nikoksr/proji/-/tree/dev",
			wantErr: false,
		},
		{
			name:    "GitLab URL with invalid @branch",
			repoURL: "https://gitlab.com/nikoksr/proji@dev@main",
			want:    "",
			wantErr: true,
		},
		{
			name:    "GitHub URL with abbreviation and @branch",
			repoURL: "gh:nikoksr/proji@dev",
			want:    "https://github.com/nikoksr/proji/tree/dev",
			wantErr: false,
		},
		{
			name:    "GitLab URL with abbreviation and @branch",
			repoURL: "gl:nikoksr/proji@dev",
			want:    "https://gitlab.com/nikoksr/proji/-/tree/dev",
			wantErr: false,
		},
		{
			name:    "No url",
			repoURL: "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseRepoURL(tc.repoURL)
			if err != nil && !tc.wantErr {
				t.Errorf("ParseRepoURL returned unexpected error: %v", err)
				return
			}

			if err == nil && tc.wantErr {
				t.Error("ParseRepoURL did not return an error")
				return
			}

			if !tc.wantErr && got.String() != tc.want {
				t.Errorf("ParseRepoURL(%q) = %v, want %v", tc.repoURL, got, tc.want)
			}
		})
	}
}

func TestExtractInfoFromURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		repoURL string
		want    PackageInfo
		wantErr bool
	}{
		{
			name:    "GitHub url with branch",
			repoURL: "https://github.com/nikoksr/proji/tree/dev",
			want: PackageInfo{
				Repo: RepoInfo{
					Owner: "nikoksr",
					Name:  "proji",
					Ref:   "dev",
				},
				Path: "",
			},
			wantErr: false,
		},
		{
			name:    "GitHub url with branch and simple path",
			repoURL: "https://github.com/nikoksr/proji/blob/main/.editorconfig",
			want: PackageInfo{
				Repo: RepoInfo{
					Owner: "nikoksr",
					Name:  "proji",
					Ref:   "main",
				},
				Path: ".editorconfig",
			},
			wantErr: false,
		},
		{
			name:    "GitHub url with branch and complex path",
			repoURL: "https://github.com/nikoksr/proji/blob/main/pkg/remote/github/github.go",
			want: PackageInfo{
				Repo: RepoInfo{
					Owner: "nikoksr",
					Name:  "proji",
					Ref:   "main",
				},
				Path: "pkg/remote/github/github.go",
			},
			wantErr: false,
		},
		{
			name:    "GitHub url minimal",
			repoURL: "https://github.com/nikoksr/proji",
			want: PackageInfo{
				Repo: RepoInfo{
					Owner: "nikoksr",
					Name:  "proji",
					Ref:   "main",
				},
				Path: "",
			},
			wantErr: false,
		},
		{
			name:    "GitHub url without branch and repo name",
			repoURL: "https://github.com/nikoksr",
			want:    PackageInfo{},
			wantErr: true,
		},
		{
			name:    "GitHub url without branch, repo name and owner",
			repoURL: "https://github.com",
			want:    PackageInfo{},
			wantErr: true,
		},
		{
			name:    "GitLab url with branch #1",
			repoURL: "https://gitlab.com/nikoksr/proji/-/tree/dev",
			want: PackageInfo{
				Repo: RepoInfo{
					Owner: "nikoksr",
					Name:  "proji",
					Ref:   "dev",
				},
				Path: "",
			},
			wantErr: false,
		},
		{
			name:    "GitLab url with branch #2",
			repoURL: "https://gitlab.com/nikoksr/proji/tree/dev",
			want: PackageInfo{
				Repo: RepoInfo{
					Owner: "nikoksr",
					Name:  "proji",
					Ref:   "dev",
				},
				Path: "",
			},
			wantErr: false,
		},
		{
			name:    "GitLab url with branch and simple path",
			repoURL: "https://gitlab.com/inkscape/inkscape/-/blob/master/.gitignore",
			want: PackageInfo{
				Repo: RepoInfo{
					Owner: "inkscape",
					Name:  "inkscape",
					Ref:   "master",
				},
				Path: ".gitignore",
			},
			wantErr: false,
		},
		{
			name:    "GitLab url with branch and complex path",
			repoURL: "https://gitlab.com/inkscape/inkscape/-/blob/master/src/extension/plugins/CMakeLists.txt",
			want: PackageInfo{
				Repo: RepoInfo{
					Owner: "inkscape",
					Name:  "inkscape",
					Ref:   "master",
				},
				Path: "src/extension/plugins/CMakeLists.txt",
			},
			wantErr: false,
		},
		{
			name:    "GitLab url without branch",
			repoURL: "https://gitlab.com/nikoksr/proji",
			want: PackageInfo{
				Repo: RepoInfo{
					Owner: "nikoksr",
					Name:  "proji",
					Ref:   "main",
				},
				Path: "",
			},
			wantErr: false,
		},
		{
			name:    "GitLab url without branch and repo name",
			repoURL: "https://gitlab.com/nikoksr",
			want:    PackageInfo{},
			wantErr: true,
		},
		{
			name:    "GitLab url without branch, repo name and owner",
			repoURL: "https://gitlab.com",
			want:    PackageInfo{},
			wantErr: true,
		},
		{
			name:    "No url",
			repoURL: "",
			want:    PackageInfo{},
			wantErr: true,
		},
		{
			name:    "Nil url",
			repoURL: "nil",
			want:    PackageInfo{},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repoURL, err := url.Parse(tc.repoURL)
			if err != nil && !tc.wantErr {
				t.Errorf("unexpected error when parsing the repoURL: %v", err)
				return
			}

			if tc.repoURL == "nil" {
				repoURL = nil
			}

			got, err := ExtractPackageInfoFromURL(context.Background(), repoURL)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ExtractPackageInfoFromURL() error = %v, wantErr %v", err, tc.wantErr)
			}
			diff := cmp.Diff(got, tc.want)
			if diff != "" {
				t.Fatalf("ExtractPackageInfoFromURL mismatch (-want +got):\n%s", diff)
			}

			// For coverage... we could make ExtractRepoInfoFromURL() use ExtractPackageInfoFromURL() internally to
			// avoid this but feels more unfitting than this kinda redundant test.
			gotRepoInfo, err := ExtractRepoInfoFromURL(context.Background(), repoURL)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ExtractRepoInfoFromURL() error = %v, wantErr %v", err, tc.wantErr)
			}
			if gotRepoInfo != tc.want.Repo {
				t.Fatalf("ExtractRepoInfoFromURL mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
