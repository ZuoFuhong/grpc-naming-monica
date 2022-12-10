package api

import (
	"fmt"
	"testing"
)

func Test_FetchAPI(t *testing.T) {
	nodes, err := Fetch("Test", "go_wallet_manage_svr")
	if err != nil {
		t.Fatal(err)
	}
	for _, node := range nodes {
		fmt.Println(node)
	}
}
