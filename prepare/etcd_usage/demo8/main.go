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
		putOp clientv3.Op
		getOp clientv3.Op
		opResp clientv3.OpResponse
	)

	//配置客户端
	config = clientv3.Config{
		Endpoints:[]string{"120.78.140.135:2379"},
		DialTimeout:5 * time.Second,
	}

	//建立连接
	client,err = clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	//创建Op:operation操作
	putOp = clientv3.OpPut("/cron/jobs/job8","")

	//执行Op
	opResp,err = kv.Do(context.TODO(),putOp)
	if err != nil{
		fmt.Println(err)
		return
	}

	//kv.Put
	//kv.Get
	//kv.Delete
	fmt.Println("写入Revision:",opResp.Put().Header.Revision)

	//创建Op
	getOp = clientv3.OpGet("/cron/jobs/job8")
	//执行op
	opResp,err = kv.Do(context.TODO(),getOp)
	if err != nil {
		fmt.Println(err)
		return
	}
	//打印
	fmt.Println("数据的Revision:",opResp.Get().Kvs[0].CreateRevision)
	fmt.Println("数据的Revision:",opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据的value:",string(opResp.Get().Kvs[0].Value))
}
