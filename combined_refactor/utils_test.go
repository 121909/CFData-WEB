package main

import (
	"reflect"
	"testing"
)

func TestParseNSBFilesInputs(t *testing.T) {
	files := []nsbUploadedFile{
		{
			Name:    "a.txt",
			Content: "1.1.1.1 443\n2.2.2.2\n",
		},
		{
			Name:    "b.csv",
			Content: "IP地址,端口号,数据中心\n2.2.2.2,443,HKG\n3.3.3.3,8443,NRT\n",
		},
	}
	got := parseNSBFilesInputs(files, 443)
	want := []string{"1.1.1.1 443", "2.2.2.2 443", "3.3.3.3 8443"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseNSBFilesInputs() = %v, want %v", got, want)
	}
}

func TestMatchesTargetDC(t *testing.T) {
	cases := []struct {
		dc     string
		target string
		want   bool
	}{
		{"HKG", "", true},
		{"HKG", "hkg", true},
		{"NRT", "HKG, NRT, SIN", true},
		{"SIN", "HKG，NRT;SIN", true},
		{"LAX", "HKG,NRT,SIN", false},
	}
	for _, tc := range cases {
		if got := matchesTargetDC(tc.dc, tc.target); got != tc.want {
			t.Errorf("matchesTargetDC(%q, %q) = %v, want %v", tc.dc, tc.target, got, tc.want)
		}
	}
}
