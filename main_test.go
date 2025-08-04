package main

import (
	"github.com/tomeai/dataflow/sdk"
	"log"
	"testing"
)

func TestSave(t *testing.T) {
	rs, err := sdk.NewResultService("77963b7a931377ad4ab5ad6a9cd718aa")
	if err != nil {
		log.Fatal("初始化失败:", err)
	}

	record := sdk.Record{
		StoreKey: "https://linux.do/t/topic/76121",
		Data: map[string]any{
			"name": "opentome",
		},
		Metadata: map[string]string{
			"name": "xiaoming",
		},
	}

	if err := rs.SaveItem(record); err != nil {
		log.Fatal("发送失败:", err)
	}
}
