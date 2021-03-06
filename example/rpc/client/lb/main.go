package main

import (
	"context"
	"fmt"
	"time"

	"zero/core/discov"
	"zero/example/rpc/remote/unary"
	"zero/rpcx"
)

func main() {
	cli := rpcx.MustNewClient(rpcx.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost:2379"},
			Key:   "rpcx",
		},
	})
	greet := unary.NewGreeterClient(cli.Conn())
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			resp, err := greet.Greet(context.Background(), &unary.Request{
				Name: "kevin",
			})
			if err != nil {
				fmt.Println("X", err.Error())
			} else {
				fmt.Println("=>", resp.Greet)
			}
		}
	}
}
