package main

import (
	"context"
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

func TestRunNSBSpeedWorkersPerDC(t *testing.T) {
	results := []iptestResult{
		{ipAddr: "1.1.1.1", dataCenter: "HKG"},
		{ipAddr: "1.1.1.2", dataCenter: "HKG"},
		{ipAddr: "1.1.1.3", dataCenter: "HKG"},
		{ipAddr: "2.2.2.1", dataCenter: "NRT"},
		{ipAddr: "2.2.2.2", dataCenter: "NRT"},
		{ipAddr: "2.2.2.3", dataCenter: "NRT"},
	}
	if got := calcNSBSpeedTargetTotal(results, 2, true); got != 4 {
		t.Fatalf("calcNSBSpeedTargetTotal() = %d, want 4", got)
	}
	canceled := runNSBSpeedWorkersWithMode(context.Background(), results, 2, 2, true, 1.0, nil, nil, func(idx int) (float64, string) {
		return 2048, "" // 2.0 MB/s
	})
	if canceled {
		t.Fatalf("expected not canceled")
	}
	counts := map[string]int{}
	for _, r := range results {
		if r.speedQualified {
			counts[r.dataCenter]++
		}
	}
	if counts["HKG"] != 2 || counts["NRT"] != 2 {
		t.Fatalf("expected 2 qualified per DC, got %+v", counts)
	}
}
