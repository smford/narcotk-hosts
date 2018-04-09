package main

import (
	"testing"
)

func TestPrepareMac(t *testing.T) {
	var tests = []string{"DeAdbEefcaFE", "de:ad:be:ef:ca:fe", "de-ad-be-ef-ca-fe"}
	for i, v := range tests {
		if PrepareMac(v) != "de:ad:be:ef:ca:fe" {
			t.Error("Test ", i, ": Expected de:ad:be:ef:ca:fe got ", v)
		}
	}
}

func TestMakePaddedIp(t *testing.T) {
	var tests = []string{"1.1.1.1", "192.168.1.1", "192.168.100.101"}
	var expectedresults = []string{"001001001001", "192168001001", "192168100101"}
	for i, v := range tests {
		if MakePaddedIp(v) != expectedresults[i] {
			t.Error("Test ", i, ": Expected: ", expectedresults[i], "  Actual: ", MakePaddedIp(v))
		}
	}
}
