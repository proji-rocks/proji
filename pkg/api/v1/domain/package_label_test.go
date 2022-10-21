package domain

import "testing"

func Test_generateLabelFromName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		label string
	}{
		{
			name:  "",
			label: "",
		},
		{
			name:  " ",
			label: "",
		},
		{
			name:  "a",
			label: "a",
		},
		{
			name:  "a_",
			label: "a",
		},
		{
			name:  "_a_",
			label: "a",
		},
		{
			name:  "a-b",
			label: "ab",
		},
		{
			name:  "a-b-c",
			label: "abc",
		},
		{
			name:  "ABC",
			label: "abc",
		},
		{
			name:  "AbcDefGhi",
			label: "adg",
		},
		{
			name:  "AbcDefGhi-JklMno",
			label: "aj",
		},
		{
			name:  "AbcDefGhi-JklMno_PqrStu",
			label: "aj",
		},
		{
			name:  "abc def ghi jkl mno pqr stu",
			label: "adgj",
		},
		{
			name:  "AbcDefGhi_JklMno_PqrStu",
			label: "ajp",
		},
		{
			name:  "AbcDefGhi_JklMno_PqrStu_VwXyZ",
			label: "ajpv",
		},
		{
			name:  "a.b.c.d",
			label: "abcd",
		},
		{
			name:  "a b c d",
			label: "abcd",
		},
		{
			name:  "a-b-c-d",
			label: "abcd",
		},
		{
			name:  "a_b_c_d",
			label: "abcd",
		},
		{
			name:  "a%20b%20c%20d",
			label: "abcd",
		},
		{
			name:  "aaaaaaaaaaabbbbbbbbbbbbbbccccccccccccccdddddddddddeeeeeeeeeeeeee",
			label: "ace",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			label := generateLabelFromName(tc.name)
			if label != tc.label {
				t.Errorf("expected label %s, got %s", tc.label, label)
			}
		})
	}
}
