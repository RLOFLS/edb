package edb

import (
	"testing"
)

func TestDns(t *testing.T) {
	tests := []struct {
		config Config
		dns    string
	}{
		{
			Config{},
			"",
		},
		{
			Config{Host: "127.0.0.1", Port: "3306", Username: "root", Password: "123", Database: "test", Charset: "utf8"},
			"",
		},
		{
			Config{Driver: "mysql", Host: "127.0.0.1", Port: "3306", Username: "root", Password: "123", Database: "test", Charset: "utf8"},
			"root:123@tcp(127.0.0.1:3306)/test?charset=utf8",
		},
	}

	for _, test := range tests {
		actDNS := test.config.DNS()
		if test.dns != actDNS {
			t.Errorf("expect: %s, actually: %s", test.dns, actDNS)
		}
	}
}
