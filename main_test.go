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

func TestParseSql(t *testing.T) {
	var validtests = []string{"select * from hosts", "select * from networks", "select * from hosts where fqdn like 'server.example.com'", "select * from hosts where network like '192.168.1'", "select * from networks where network like '192.168.1'", "select * from hosts where ipaddress like '192.168.1.1'", "select * from hosts where mac like 'de:ad:be:ef:ca:fe'"}
	var invalidtests = []string{"random junk", "more junk", "even more junk"}
	for i, v := range validtests {
		if !ParseSql(v) {
			t.Error("Test ", i, ": Expected: True  Actual: ", v)
		}
	}
	for i, v := range invalidtests {
		if ParseSql(v) {
			t.Error("Test ", i, ": Expected: True  Actual: ", v)
		}
	}
}

func TestValidIP(t *testing.T) {
	var validtests = []string{"1.1.1.1", "192.168.1.1", "192.168.100.101"}
	var invalidtests = []string{"256.256.256.256", "99999999", "a.b.c.d"}
	for i, v := range validtests {
		if !ValidIP(v) {
			t.Error("Test ", i, ": Expected: True  Actual: ", v)
		}
	}
	for i, v := range invalidtests {
		if ValidIP(v) {
			t.Error("Test ", i, ": Expected: True  Actual: ", v)
		}
	}
}
