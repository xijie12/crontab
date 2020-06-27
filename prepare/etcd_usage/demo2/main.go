package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"context"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		putResp *clientv3.PutResponse
	)
	config = clientv3.Config{
		Endpoints:[]string{"120.78.140.135:2379"},
		DialTimeout:5 * time.Second,
	}
	if client,err = clientv3.New(config);err != nil {
		fmt.Println(err)
		return
	}
	//用于读写etcd的键值对
	kv = clientv3.NewKV(client)

	//putResp,err = kv.Put(context.TODO(),"/cron/jobs/job1","hello")
	//请求包含上一个版本的kv
	putResp,err = kv.Put(context.TODO(),"/cron/jobs/job1","bye",clientv3.WithPrevKV())
	if err != nil{
		fmt.Println(err)
		return
	}else{
		//kv值是按版本存储的，每个版本都会存储
		fmt.Println("Revision:",putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("上一个版本的kv信息：",putResp.PrevKv)
		}
	}
}
