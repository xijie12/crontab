package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"fmt"
	"context"
)

func main() {
	//事务txn实现分布式锁,乐观锁
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		lease clientv3.Lease
		leaseId clientv3.LeaseID
		leaseGrandResp *clientv3.LeaseGrantResponse
		keepResp *clientv3.LeaseKeepAliveResponse
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		ctx context.Context
		cancelFunc context.CancelFunc
		kv clientv3.KV
		txn clientv3.Txn
		txnResp *clientv3.TxnResponse
	)

	//客户端配置
	config = clientv3.Config{
		Endpoints:[]string{"120.78.140.135:2379"},
		DialTimeout:5 * time.Second,
	}

	//建立连接
	client, err = clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	//lease实现锁自动过期：
	//op操作
	//txn事务：if else then

	//1.上锁（创建租约，自动续租，拿着租约去抢占一个key）
	//申请一个lease（租约）
	lease = clientv3.NewLease(client)

	//申请一个5秒的租约
	leaseGrandResp,err = lease.Grant(context.TODO(),5)
	if err != nil {
		fmt.Println(err)
		return
	}

	//拿到租约的ID
	leaseId = leaseGrandResp.ID

	//准备一个用于取消自动续租的context
	ctx, cancelFunc = context.WithCancel(context.TODO())

	//确保函数退出后，自动续租会停止
	defer cancelFunc() //keepRespChan == nil
	defer lease.Revoke(context.TODO(),leaseId)

	//5秒后会取消自动续租
	keepRespChan,err = lease.KeepAlive(ctx,leaseId)

	//处理续约应答的协程
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

	//if 不存在key，then设置它，else枪锁失败
	kv = clientv3.NewKV(client)

	//创建事务
	txn = kv.Txn(context.TODO())

	//定义事务
	//如果key不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"),"=",0)).
		Then(clientv3.OpPut("/cron/lock/job9","xxx",clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job9"))//否则枪锁失败

	//提交事务
	txnResp,err = txn.Commit()
	if err != nil {
		fmt.Println(err)
		return
	}

	//判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用：",string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	//2.处理业务
	//在锁内，很安全
	fmt.Println("处理业务")
	time.Sleep(5 * time.Second)

	//3.释放锁（取消自动续租，释放租约）
	//defer会把租约释放掉，关联的KV就会被删除了,在1步骤上
}
