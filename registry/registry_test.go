package registry

import (
	"testing"
)

func Test_Register(t *testing.T) {
	regIns := NewRegistry(&Config{
		Token:       "18ee7064-3cdd-4ed5-a139-fd8d9add5847",
		Namespace:   "Test",
		ServiceName: "go_wallet_manage_svr",
		IP:          "127.0.0.1",
		Port:        1024,
		Weight:      100,
		Metadata:    "[]",
		HeartBeat:   5,
	})
	if err := regIns.Register(); err != nil {
		t.Fatal(err)
	}
}
