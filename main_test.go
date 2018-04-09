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
