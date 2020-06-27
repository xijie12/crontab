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
		lease clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		keepResp *clientv3.LeaseKeepAliveResponse
		keepRespChan <- chan *clientv3.LeaseKeepAliveResponse
		kv clientv3.KV
	)
	config = clientv3.Config{
		Endpoints:[]string{"120.78.140.135:2379"},
		DialTimeout:5 * time.Second,
	}
	if client,err = clientv3.New(config);err != nil {
		fmt.Println(err)
		return
	}

	//申请一个lease（租约）
	lease = clientv3.NewLease(client)

	//申请一个10秒的租约
	leaseGrantResp,err = lease.Grant(context.TODO(),10)
	if err != nil {
		fmt.Println(err)
		return
	}

	//拿到租约的ID
	leaseId = leaseGrantResp.ID

	//自动续租
	//ctx,_ := context.WithTimeout(context.TODO(),5 * time.Second)
	//5秒后会取消自动续租
	//keepRespChan, err = lease.KeepAlive(ctx,leaseId)
	keepRespChan, err = lease.KeepAlive(context.TODO(),leaseId)
	if err != nil {
		fmt.Println(err)
		return
	}

	//处理续租应答的协程
	go func() {
		for{
			select {
			case keepResp = <- keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约已经失效了")
					goto END
				}else{//每秒会续租一次，所以就会收到一次应答
					fmt.Println("收到自动续租应答：",keepResp.ID)
				}
			}
		}
		END:
	}()

	//获取KV API子集
	kv = clientv3.NewKV(client)

	//Put一个KV，让它与租约关联起来，从而实现10秒后自动过期
	putResp, err = kv.Put(context.TODO(),"/cron/lock/job1","",clientv3.WithLease(leaseId))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功：",putResp.Header.Revision)

	//定时的看一下key过期了没有
	for{
		getResp, err = kv.Get(context.TODO(),"/cron/lock/job1")
		if err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没过期：",getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}
