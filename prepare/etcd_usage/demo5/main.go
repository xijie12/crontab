package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"context"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		delResp *clientv3.DeleteResponse
		kvpair *mvccpb.KeyValue
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

	//删除KV
	delResp,err = kv.Delete(context.TODO(),"/cron/jobs/job2",clientv3.WithPrevKV())
	//删除以/cron/jobs/为前缀的多个key
	//delResp,err = kv.Delete(context.TODO(),"/cron/jobs/",clientv3.WithPrefix(),clientv3.WithPrevKV())
	//删除以/cron/jobs/job1的开始往后的key
	//delResp,err = kv.Delete(context.TODO(),"/cron/jobs/",clientv3.WithFromKey(),clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(delResp)
		return
	}
	//被删除之前的value是什么
	if len(delResp.PrevKvs) != 0 {
		for _,kvpair = range delResp.PrevKvs{
			fmt.Println("删除了：",string(kvpair.Key),string(kvpair.Value))
		}
	}
}
